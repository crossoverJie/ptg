// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.5.1
// source: reflect/gen/common.proto

package v1

import (
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

type Common struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header string `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
}

func (x *Common) Reset() {
	*x = Common{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reflect_gen_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Common) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Common) ProtoMessage() {}

func (x *Common) ProtoReflect() protoreflect.Message {
	mi := &file_reflect_gen_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Common.ProtoReflect.Descriptor instead.
func (*Common) Descriptor() ([]byte, []int) {
	return file_reflect_gen_common_proto_rawDescGZIP(), []int{0}
}

func (x *Common) GetHeader() string {
	if x != nil {
		return x.Header
	}
	return ""
}

var File_reflect_gen_common_proto protoreflect.FileDescriptor

var file_reflect_gen_common_proto_rawDesc = []byte{
	0x0a, 0x18, 0x72, 0x65, 0x66, 0x6c, 0x65, 0x63, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x22, 0x20, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x6f, 0x76, 0x65, 0x72, 0x4a, 0x69,
	0x65, 0x2f, 0x70, 0x74, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_reflect_gen_common_proto_rawDescOnce sync.Once
	file_reflect_gen_common_proto_rawDescData = file_reflect_gen_common_proto_rawDesc
)

func file_reflect_gen_common_proto_rawDescGZIP() []byte {
	file_reflect_gen_common_proto_rawDescOnce.Do(func() {
		file_reflect_gen_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_reflect_gen_common_proto_rawDescData)
	})
	return file_reflect_gen_common_proto_rawDescData
}

var file_reflect_gen_common_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_reflect_gen_common_proto_goTypes = []interface{}{
	(*Common)(nil), // 0: order.v1.Common
}
var file_reflect_gen_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_reflect_gen_common_proto_init() }
func file_reflect_gen_common_proto_init() {
	if File_reflect_gen_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_reflect_gen_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Common); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_reflect_gen_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_reflect_gen_common_proto_goTypes,
		DependencyIndexes: file_reflect_gen_common_proto_depIdxs,
		MessageInfos:      file_reflect_gen_common_proto_msgTypes,
	}.Build()
	File_reflect_gen_common_proto = out.File
	file_reflect_gen_common_proto_rawDesc = nil
	file_reflect_gen_common_proto_goTypes = nil
	file_reflect_gen_common_proto_depIdxs = nil
}
