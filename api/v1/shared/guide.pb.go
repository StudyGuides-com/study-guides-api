// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: v1/shared/guide.proto

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

type GuideData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	ParserType    ParserType             `protobuf:"varint,3,opt,name=parser_type,json=parserType,proto3,enum=shared.v1.ParserType" json:"parser_type,omitempty"`
	Sections      []*SectionData         `protobuf:"bytes,4,rep,name=sections,proto3" json:"sections,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GuideData) Reset() {
	*x = GuideData{}
	mi := &file_v1_shared_guide_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GuideData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuideData) ProtoMessage() {}

func (x *GuideData) ProtoReflect() protoreflect.Message {
	mi := &file_v1_shared_guide_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuideData.ProtoReflect.Descriptor instead.
func (*GuideData) Descriptor() ([]byte, []int) {
	return file_v1_shared_guide_proto_rawDescGZIP(), []int{0}
}

func (x *GuideData) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *GuideData) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *GuideData) GetParserType() ParserType {
	if x != nil {
		return x.ParserType
	}
	return ParserType_PARSER_TYPE_UNSPECIFIED
}

func (x *GuideData) GetSections() []*SectionData {
	if x != nil {
		return x.Sections
	}
	return nil
}

var File_v1_shared_guide_proto protoreflect.FileDescriptor

const file_v1_shared_guide_proto_rawDesc = "" +
	"\n" +
	"\x15v1/shared/guide.proto\x12\tshared.v1\x1a\x17v1/shared/section.proto\x1a\x1av1/shared/parsertype.proto\"\xaf\x01\n" +
	"\tGuideData\x12\x14\n" +
	"\x05title\x18\x01 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x126\n" +
	"\vparser_type\x18\x03 \x01(\x0e2\x15.shared.v1.ParserTypeR\n" +
	"parserType\x122\n" +
	"\bsections\x18\x04 \x03(\v2\x16.shared.v1.SectionDataR\bsectionsBDZBgithub.com/studyguides-com/study-guides-api/api/v1/shared;sharedv1b\x06proto3"

var (
	file_v1_shared_guide_proto_rawDescOnce sync.Once
	file_v1_shared_guide_proto_rawDescData []byte
)

func file_v1_shared_guide_proto_rawDescGZIP() []byte {
	file_v1_shared_guide_proto_rawDescOnce.Do(func() {
		file_v1_shared_guide_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_shared_guide_proto_rawDesc), len(file_v1_shared_guide_proto_rawDesc)))
	})
	return file_v1_shared_guide_proto_rawDescData
}

var file_v1_shared_guide_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_v1_shared_guide_proto_goTypes = []any{
	(*GuideData)(nil),   // 0: shared.v1.GuideData
	(ParserType)(0),     // 1: shared.v1.ParserType
	(*SectionData)(nil), // 2: shared.v1.SectionData
}
var file_v1_shared_guide_proto_depIdxs = []int32{
	1, // 0: shared.v1.GuideData.parser_type:type_name -> shared.v1.ParserType
	2, // 1: shared.v1.GuideData.sections:type_name -> shared.v1.SectionData
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_v1_shared_guide_proto_init() }
func file_v1_shared_guide_proto_init() {
	if File_v1_shared_guide_proto != nil {
		return
	}
	file_v1_shared_section_proto_init()
	file_v1_shared_parsertype_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_shared_guide_proto_rawDesc), len(file_v1_shared_guide_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_shared_guide_proto_goTypes,
		DependencyIndexes: file_v1_shared_guide_proto_depIdxs,
		MessageInfos:      file_v1_shared_guide_proto_msgTypes,
	}.Build()
	File_v1_shared_guide_proto = out.File
	file_v1_shared_guide_proto_goTypes = nil
	file_v1_shared_guide_proto_depIdxs = nil
}
