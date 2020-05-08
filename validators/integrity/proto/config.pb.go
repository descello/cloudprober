// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.2
// source: github.com/google/cloudprober/validators/integrity/proto/config.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Validator struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Validate the data integrity of the response using a pattern that is
	// repeated throughout the length of the response, with last len(response) %
	// len(pattern) bytes being zero bytes.
	//
	// For example if response length is 100 bytes and pattern length is 8 bytes,
	// first 96 bytes of the response should be pattern repeated 12 times, and
	// last 4 bytes should be set to zero byte ('\0')
	//
	// Types that are assignable to Pattern:
	//	*Validator_PatternString
	//	*Validator_PatternNumBytes
	Pattern isValidator_Pattern `protobuf_oneof:"pattern"`
}

func (x *Validator) Reset() {
	*x = Validator{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_google_cloudprober_validators_integrity_proto_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Validator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Validator) ProtoMessage() {}

func (x *Validator) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_google_cloudprober_validators_integrity_proto_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Validator.ProtoReflect.Descriptor instead.
func (*Validator) Descriptor() ([]byte, []int) {
	return file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDescGZIP(), []int{0}
}

func (m *Validator) GetPattern() isValidator_Pattern {
	if m != nil {
		return m.Pattern
	}
	return nil
}

func (x *Validator) GetPatternString() string {
	if x, ok := x.GetPattern().(*Validator_PatternString); ok {
		return x.PatternString
	}
	return ""
}

func (x *Validator) GetPatternNumBytes() int32 {
	if x, ok := x.GetPattern().(*Validator_PatternNumBytes); ok {
		return x.PatternNumBytes
	}
	return 0
}

type isValidator_Pattern interface {
	isValidator_Pattern()
}

type Validator_PatternString struct {
	// Pattern string for pattern repetition based integrity checks.
	// For example, cloudprobercloudprobercloudprober...
	PatternString string `protobuf:"bytes,1,opt,name=pattern_string,json=patternString,oneof"`
}

type Validator_PatternNumBytes struct {
	// Pattern is derived from the first few bytes of the payload. This is
	// useful when pattern is not known in advance, for example cloudprober's
	// ping probe repeates the timestamp (8 bytes) in the packet payload.
	// An error is returned if response is smaller than pattern_num_bytes.
	PatternNumBytes int32 `protobuf:"varint,2,opt,name=pattern_num_bytes,json=patternNumBytes,oneof"`
}

func (*Validator_PatternString) isValidator_Pattern() {}

func (*Validator_PatternNumBytes) isValidator_Pattern() {}

var File_github_com_google_cloudprober_validators_integrity_proto_config_proto protoreflect.FileDescriptor

var file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDesc = []byte{
	0x0a, 0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x72, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x67,
	0x72, 0x69, 0x74, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x20, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x70, 0x72,
	0x6f, 0x62, 0x65, 0x72, 0x2e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x2e,
	0x69, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x69, 0x74, 0x79, 0x22, 0x6d, 0x0a, 0x09, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x27, 0x0a, 0x0e, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72,
	0x6e, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x0d, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12,
	0x2c, 0x0a, 0x11, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x5f, 0x6e, 0x75, 0x6d, 0x5f, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x0f, 0x70, 0x61,
	0x74, 0x74, 0x65, 0x72, 0x6e, 0x4e, 0x75, 0x6d, 0x42, 0x79, 0x74, 0x65, 0x73, 0x42, 0x09, 0x0a,
	0x07, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x72, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x6f, 0x72, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x69, 0x74, 0x79, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f,
}

var (
	file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDescOnce sync.Once
	file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDescData = file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDesc
)

func file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDescGZIP() []byte {
	file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDescOnce.Do(func() {
		file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDescData)
	})
	return file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDescData
}

var file_github_com_google_cloudprober_validators_integrity_proto_config_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_google_cloudprober_validators_integrity_proto_config_proto_goTypes = []interface{}{
	(*Validator)(nil), // 0: cloudprober.validators.integrity.Validator
}
var file_github_com_google_cloudprober_validators_integrity_proto_config_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_google_cloudprober_validators_integrity_proto_config_proto_init() }
func file_github_com_google_cloudprober_validators_integrity_proto_config_proto_init() {
	if File_github_com_google_cloudprober_validators_integrity_proto_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_google_cloudprober_validators_integrity_proto_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Validator); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_github_com_google_cloudprober_validators_integrity_proto_config_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Validator_PatternString)(nil),
		(*Validator_PatternNumBytes)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_google_cloudprober_validators_integrity_proto_config_proto_goTypes,
		DependencyIndexes: file_github_com_google_cloudprober_validators_integrity_proto_config_proto_depIdxs,
		MessageInfos:      file_github_com_google_cloudprober_validators_integrity_proto_config_proto_msgTypes,
	}.Build()
	File_github_com_google_cloudprober_validators_integrity_proto_config_proto = out.File
	file_github_com_google_cloudprober_validators_integrity_proto_config_proto_rawDesc = nil
	file_github_com_google_cloudprober_validators_integrity_proto_config_proto_goTypes = nil
	file_github_com_google_cloudprober_validators_integrity_proto_config_proto_depIdxs = nil
}
