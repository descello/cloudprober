// Copyright 2022 The Cloudprober Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package probestatus implements a surfacer that exposes probes' status over web
interface. This surfacer builds an in-memory timeseries database from the
incoming EventMetrics.
*/
package probestatus

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cloudprober/cloudprober/config/runconfig"
	"github.com/cloudprober/cloudprober/logger"
	"github.com/cloudprober/cloudprober/metrics"
	"github.com/cloudprober/cloudprober/surfacers/common/options"
	configpb "github.com/cloudprober/cloudprober/surfacers/probestatus/proto"
	"github.com/cloudprober/cloudprober/sysvars"
)

const (
	metricsBufferSize = 10000
)

// queriesQueueSize defines how many queries can we queue before we start
// blocking on previous queries to finish.
const queriesQueueSize = 10

// httpWriter is a wrapper for http.ResponseWriter that includes a channel
// to signal the completion of the writing of the response.
type httpWriter struct {
	w        http.ResponseWriter
	doneChan chan struct{}
}

type pageCache struct {
	mu         sync.RWMutex
	content    []byte
	cachedTime time.Time
	maxAge     time.Duration
}

func (pc *pageCache) contentIfValid() ([]byte, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	if time.Now().Sub(pc.cachedTime) > pc.maxAge {
		return nil, false
	}
	return pc.content, true
}

func (pc *pageCache) setContent(content []byte) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.content, pc.cachedTime = content, time.Now()
}

// Surfacer implements a status surfacer for Cloudprober.
type Surfacer struct {
	c         *configpb.SurfacerConf // Configuration
	opts      *options.Options
	emChan    chan *metrics.EventMetrics // Buffered channel to store incoming EventMetrics
	queryChan chan *httpWriter           // Query channel
	l         *logger.Logger

	resolution   time.Duration
	metrics      map[string]map[string]*timeseries
	probeNames   []string
	probeTargets map[string][]string

	// Dashboard page cache.
	pageCache *pageCache

	// Dashboard Metadata
	dashboardDurations []time.Duration
}

// New returns a probestatus surfacer based on the config provided. It sets up
// a goroutine to process both the incoming EventMetrics and the web requests
// for the URL handler /metrics.
func New(ctx context.Context, config *configpb.SurfacerConf, opts *options.Options, l *logger.Logger) *Surfacer {
	if config == nil {
		config = &configpb.SurfacerConf{}
	}

	res := time.Duration(config.GetResolutionSec()) * time.Second
	if res == 0 {
		res = time.Minute
	}

	ps := &Surfacer{
		c:            config,
		opts:         opts,
		emChan:       make(chan *metrics.EventMetrics, metricsBufferSize),
		queryChan:    make(chan *httpWriter, queriesQueueSize),
		metrics:      make(map[string]map[string]*timeseries),
		probeTargets: make(map[string][]string),
		resolution:   res,
		l:            l,
	}

	ps.dashboardDurations = dashboardDurations(ps.resolution * time.Duration(ps.c.GetTimeseriesSize()))
	ps.pageCache = &pageCache{
		maxAge: time.Duration(ps.c.GetCacheTimeSec()) * time.Second,
	}

	// Start a goroutine to process the incoming EventMetrics as well as
	// the incoming web queries. To avoid data access race conditions, we do
	// one thing at a time.
	go func() {
		for {
			select {
			case <-ctx.Done():
				ps.l.Infof("Context canceled, stopping the input/output processing loop.")
				return
			case em := <-ps.emChan:
				ps.record(em)
			case hw := <-ps.queryChan:
				ps.writeData(hw.w)
				close(hw.doneChan)
			}
		}
	}()

	http.HandleFunc(config.GetUrl(), func(w http.ResponseWriter, r *http.Request) {
		// doneChan is used to track the completion of the response writing. This is
		// required as response is written in a different goroutine.
		doneChan := make(chan struct{}, 1)
		ps.queryChan <- &httpWriter{w, doneChan}
		<-doneChan
	})

	l.Infof("Initialized status surfacer at the URL: %s", "probesstatus")
	return ps
}

// Write queues the incoming data into a channel. This channel is watched by a
// goroutine that actually processes the data and updates the in-memory
// database.
func (ps *Surfacer) Write(_ context.Context, em *metrics.EventMetrics) {
	select {
	case ps.emChan <- em:
	default:
		ps.l.Errorf("Surfacer's write channel is full, dropping new data.")
	}
}

// record processes the incoming EventMetrics and updates the in-memory
// database.
func (ps *Surfacer) record(em *metrics.EventMetrics) {
	probeName, targetName := em.Label("probe"), em.Label("dst")
	if probeName == "sysvars" || em.Metric("total") == nil {
		return
	}

	total, totalOk := em.Metric("total").(metrics.NumValue)
	success, successOk := em.Metric("success").(metrics.NumValue)
	if !totalOk || !successOk {
		return
	}

	probeTS := ps.metrics[probeName]
	if probeTS == nil {
		probeTS = make(map[string]*timeseries)
		ps.metrics[probeName] = probeTS
		ps.probeNames = append(ps.probeNames, probeName)
	}

	targetTS := probeTS[targetName]
	if targetTS == nil {
		if len(probeTS) == int(ps.c.GetMaxTargetsPerProbe())-1 {
			ps.l.Warningf("Reached the per-probe timeseries capacity (%d) with target \"%s\". All new targets will be silently dropped.", ps.c.GetMaxTargetsPerProbe(), targetName)
		}
		if len(probeTS) >= int(ps.c.GetMaxTargetsPerProbe()) {
			return
		}
		targetTS = newTimeseries(ps.resolution, int(ps.c.GetTimeseriesSize()))
		probeTS[targetName] = targetTS
		ps.probeTargets[probeName] = append(ps.probeTargets[probeName], targetName)
	}

	targetTS.addDatum(em.Timestamp, &datum{
		total:   total.Int64(),
		success: success.Int64(),
		latency: em.Metric("latency").Clone(),
	})
}

func (ps *Surfacer) probeStatus(probeName string, durations []time.Duration) ([]string, []string) {
	var lines, debugLines []string

	for _, targetName := range ps.probeTargets[probeName] {
		lines = append(lines, "<tr><td><b>"+targetName+"</b></td>")
		ts := ps.metrics[probeName][targetName]
		data := ts.getRecentData(24 * time.Hour)

		for _, td := range durations {
			t, s := ts.computeDelta(data, td)
			lines = append(lines, fmt.Sprintf("<td>%.4f</td>", float64(s)/float64(t)))
		}

		debugLines = append(debugLines, fmt.Sprintf("Target: %s, Oldest timestamp: %s<br>",
			targetName, ts.currentTS.Add(time.Duration(-len(data))*ts.res)))

		for _, i := range []int{0, len(data) - 1} {
			d := data[i]
			debugLines = append(debugLines, fmt.Sprintf("#%d total=%d, success=%d, latency=%s <br>",
				i, d.total, d.success, d.latency.String()))
		}
	}
	return lines, debugLines
}

func (ps *Surfacer) writeData(w io.Writer) {
	content, valid := ps.pageCache.contentIfValid()
	if valid {
		w.Write(content)
		return
	}

	startTime := sysvars.StartTime().Truncate(time.Millisecond)
	uptime := time.Since(startTime).Truncate(time.Millisecond)

	probesStatus := make(map[string]template.HTML)
	probesStatusDebug := make(map[string]template.HTML)

	for _, probeName := range ps.probeNames {
		probeLines, probeDebugLines := ps.probeStatus(probeName, ps.dashboardDurations)
		probesStatus[probeName] = template.HTML(strings.Join(probeLines, "\n"))
		probesStatusDebug[probeName] = template.HTML(strings.Join(probeDebugLines, "\n"))
	}

	var statusBuf bytes.Buffer

	tmpl, err := template.New("statusTmpl").Parse(probeStatusTmpl)
	if err != nil {
		ps.l.Errorf("Error parsing probe status template: %v", err)
		return
	}
	tmpl.Execute(&statusBuf, struct {
		Durations         []string
		ProbeNames        []string
		ProbesStatus      map[string]template.HTML
		ProbesStatusDebug map[string]template.HTML
		Version           string
		StartTime, Uptime fmt.Stringer
	}{
		Durations:         shortDur(ps.dashboardDurations),
		ProbeNames:        ps.probeNames,
		ProbesStatus:      probesStatus,
		ProbesStatusDebug: probesStatusDebug,
		Version:           runconfig.Version(),
		StartTime:         startTime,
		Uptime:            uptime,
	})

	ps.pageCache.setContent(statusBuf.Bytes())
	w.Write(statusBuf.Bytes())
}
