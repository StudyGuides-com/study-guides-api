// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: v1/shared/exporttype.proto

package sharedv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ExportType int32

const (
	ExportType_EXPORT_TYPE_UNSPECIFIED ExportType = 0
	ExportType_EXPORT_TYPE_PROD        ExportType = 1
	ExportType_EXPORT_TYPE_TEST        ExportType = 2
	ExportType_EXPORT_TYPE_DEV         ExportType = 3
)

// Enum value maps for ExportType.
var (
	ExportType_name = map[int32]string{
		0: "EXPORT_TYPE_UNSPECIFIED",
		1: "EXPORT_TYPE_PROD",
		2: "EXPORT_TYPE_TEST",
		3: "EXPORT_TYPE_DEV",
	}
	ExportType_value = map[string]int32{
		"EXPORT_TYPE_UNSPECIFIED": 0,
		"EXPORT_TYPE_PROD":        1,
		"EXPORT_TYPE_TEST":        2,
		"EXPORT_TYPE_DEV":         3,
	}
)

func (x ExportType) Enum() *ExportType {
	p := new(ExportType)
	*p = x
	return p
}

func (x ExportType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ExportType) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_shared_exporttype_proto_enumTypes[0].Descriptor()
}

func (ExportType) Type() protoreflect.EnumType {
	return &file_v1_shared_exporttype_proto_enumTypes[0]
}

func (x ExportType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ExportType.Descriptor instead.
func (ExportType) EnumDescriptor() ([]byte, []int) {
	return file_v1_shared_exporttype_proto_rawDescGZIP(), []int{0}
}

var File_v1_shared_exporttype_proto protoreflect.FileDescriptor

const file_v1_shared_exporttype_proto_rawDesc = "" +
	"\n" +
	"\x1av1/shared/exporttype.proto\x12\tshared.v1*j\n" +
	"\n" +
	"ExportType\x12\x1b\n" +
	"\x17EXPORT_TYPE_UNSPECIFIED\x10\x00\x12\x14\n" +
	"\x10EXPORT_TYPE_PROD\x10\x01\x12\x14\n" +
	"\x10EXPORT_TYPE_TEST\x10\x02\x12\x13\n" +
	"\x0fEXPORT_TYPE_DEV\x10\x03BDZBgithub.com/studyguides-com/study-guides-api/api/v1/shared;sharedv1b\x06proto3"

var (
	file_v1_shared_exporttype_proto_rawDescOnce sync.Once
	file_v1_shared_exporttype_proto_rawDescData []byte
)

func file_v1_shared_exporttype_proto_rawDescGZIP() []byte {
	file_v1_shared_exporttype_proto_rawDescOnce.Do(func() {
		file_v1_shared_exporttype_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_shared_exporttype_proto_rawDesc), len(file_v1_shared_exporttype_proto_rawDesc)))
	})
	return file_v1_shared_exporttype_proto_rawDescData
}

var file_v1_shared_exporttype_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_v1_shared_exporttype_proto_goTypes = []any{
	(ExportType)(0), // 0: shared.v1.ExportType
}
var file_v1_shared_exporttype_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_shared_exporttype_proto_init() }
func file_v1_shared_exporttype_proto_init() {
	if File_v1_shared_exporttype_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_shared_exporttype_proto_rawDesc), len(file_v1_shared_exporttype_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_shared_exporttype_proto_goTypes,
		DependencyIndexes: file_v1_shared_exporttype_proto_depIdxs,
		EnumInfos:         file_v1_shared_exporttype_proto_enumTypes,
	}.Build()
	File_v1_shared_exporttype_proto = out.File
	file_v1_shared_exporttype_proto_goTypes = nil
	file_v1_shared_exporttype_proto_depIdxs = nil
}
