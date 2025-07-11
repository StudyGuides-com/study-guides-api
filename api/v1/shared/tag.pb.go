// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: v1/shared/tag.proto

package sharedv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type Tag struct {
	state              protoimpl.MessageState  `protogen:"open.v1"`
	Id                 string                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	BatchId            *string                 `protobuf:"bytes,2,opt,name=batch_id,json=batchId,proto3,oneof" json:"batch_id,omitempty"`
	Hash               string                  `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Name               string                  `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Description        *string                 `protobuf:"bytes,5,opt,name=description,proto3,oneof" json:"description,omitempty"`
	Type               TagType                 `protobuf:"varint,6,opt,name=type,proto3,enum=shared.v1.TagType" json:"type,omitempty"`
	Context            ContextType             `protobuf:"varint,7,opt,name=context,proto3,enum=shared.v1.ContextType" json:"context,omitempty"`
	ParentTagId        *string                 `protobuf:"bytes,8,opt,name=parent_tag_id,json=parentTagId,proto3,oneof" json:"parent_tag_id,omitempty"`
	ContentRating      ContentRating           `protobuf:"varint,9,opt,name=content_rating,json=contentRating,proto3,enum=shared.v1.ContentRating" json:"content_rating,omitempty"`
	ContentDescriptors []ContentDescriptorType `protobuf:"varint,10,rep,packed,name=content_descriptors,json=contentDescriptors,proto3,enum=shared.v1.ContentDescriptorType" json:"content_descriptors,omitempty"`
	MetaTags           []string                `protobuf:"bytes,11,rep,name=meta_tags,json=metaTags,proto3" json:"meta_tags,omitempty"`
	Public             bool                    `protobuf:"varint,12,opt,name=public,proto3" json:"public,omitempty"`
	AccessCount        int32                   `protobuf:"varint,13,opt,name=access_count,json=accessCount,proto3" json:"access_count,omitempty"`
	Metadata           *Metadata               `protobuf:"bytes,14,opt,name=metadata,proto3" json:"metadata,omitempty"`
	CreatedAt          *timestamppb.Timestamp  `protobuf:"bytes,15,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt          *timestamppb.Timestamp  `protobuf:"bytes,16,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	OwnerId            *string                 `protobuf:"bytes,17,opt,name=owner_id,json=ownerId,proto3,oneof" json:"owner_id,omitempty"`
	HasQuestions       bool                    `protobuf:"varint,18,opt,name=has_questions,json=hasQuestions,proto3" json:"has_questions,omitempty"`
	HasChildren        bool                    `protobuf:"varint,19,opt,name=has_children,json=hasChildren,proto3" json:"has_children,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Tag) Reset() {
	*x = Tag{}
	mi := &file_v1_shared_tag_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_v1_shared_tag_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_v1_shared_tag_proto_rawDescGZIP(), []int{0}
}

func (x *Tag) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Tag) GetBatchId() string {
	if x != nil && x.BatchId != nil {
		return *x.BatchId
	}
	return ""
}

func (x *Tag) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *Tag) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tag) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Tag) GetType() TagType {
	if x != nil {
		return x.Type
	}
	return TagType_Category
}

func (x *Tag) GetContext() ContextType {
	if x != nil {
		return x.Context
	}
	return ContextType_Colleges
}

func (x *Tag) GetParentTagId() string {
	if x != nil && x.ParentTagId != nil {
		return *x.ParentTagId
	}
	return ""
}

func (x *Tag) GetContentRating() ContentRating {
	if x != nil {
		return x.ContentRating
	}
	return ContentRating_Unspecified
}

func (x *Tag) GetContentDescriptors() []ContentDescriptorType {
	if x != nil {
		return x.ContentDescriptors
	}
	return nil
}

func (x *Tag) GetMetaTags() []string {
	if x != nil {
		return x.MetaTags
	}
	return nil
}

func (x *Tag) GetPublic() bool {
	if x != nil {
		return x.Public
	}
	return false
}

func (x *Tag) GetAccessCount() int32 {
	if x != nil {
		return x.AccessCount
	}
	return 0
}

func (x *Tag) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Tag) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Tag) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Tag) GetOwnerId() string {
	if x != nil && x.OwnerId != nil {
		return *x.OwnerId
	}
	return ""
}

func (x *Tag) GetHasQuestions() bool {
	if x != nil {
		return x.HasQuestions
	}
	return false
}

func (x *Tag) GetHasChildren() bool {
	if x != nil {
		return x.HasChildren
	}
	return false
}

var File_v1_shared_tag_proto protoreflect.FileDescriptor

const file_v1_shared_tag_proto_rawDesc = "" +
	"\n" +
	"\x13v1/shared/tag.proto\x12\tshared.v1\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x17v1/shared/tagtype.proto\x1a\x1dv1/shared/contentrating.proto\x1a%v1/shared/contentdescriptortype.proto\x1a\x18v1/shared/metadata.proto\x1a\x1bv1/shared/contexttype.proto\"\xbe\x06\n" +
	"\x03Tag\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1e\n" +
	"\bbatch_id\x18\x02 \x01(\tH\x00R\abatchId\x88\x01\x01\x12\x12\n" +
	"\x04hash\x18\x03 \x01(\tR\x04hash\x12\x12\n" +
	"\x04name\x18\x04 \x01(\tR\x04name\x12%\n" +
	"\vdescription\x18\x05 \x01(\tH\x01R\vdescription\x88\x01\x01\x12&\n" +
	"\x04type\x18\x06 \x01(\x0e2\x12.shared.v1.TagTypeR\x04type\x120\n" +
	"\acontext\x18\a \x01(\x0e2\x16.shared.v1.ContextTypeR\acontext\x12'\n" +
	"\rparent_tag_id\x18\b \x01(\tH\x02R\vparentTagId\x88\x01\x01\x12?\n" +
	"\x0econtent_rating\x18\t \x01(\x0e2\x18.shared.v1.ContentRatingR\rcontentRating\x12Q\n" +
	"\x13content_descriptors\x18\n" +
	" \x03(\x0e2 .shared.v1.ContentDescriptorTypeR\x12contentDescriptors\x12\x1b\n" +
	"\tmeta_tags\x18\v \x03(\tR\bmetaTags\x12\x16\n" +
	"\x06public\x18\f \x01(\bR\x06public\x12!\n" +
	"\faccess_count\x18\r \x01(\x05R\vaccessCount\x12/\n" +
	"\bmetadata\x18\x0e \x01(\v2\x13.shared.v1.MetadataR\bmetadata\x129\n" +
	"\n" +
	"created_at\x18\x0f \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\x10 \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\x12\x1e\n" +
	"\bowner_id\x18\x11 \x01(\tH\x03R\aownerId\x88\x01\x01\x12#\n" +
	"\rhas_questions\x18\x12 \x01(\bR\fhasQuestions\x12!\n" +
	"\fhas_children\x18\x13 \x01(\bR\vhasChildrenB\v\n" +
	"\t_batch_idB\x0e\n" +
	"\f_descriptionB\x10\n" +
	"\x0e_parent_tag_idB\v\n" +
	"\t_owner_idBDZBgithub.com/studyguides-com/study-guides-api/api/v1/shared;sharedv1b\x06proto3"

var (
	file_v1_shared_tag_proto_rawDescOnce sync.Once
	file_v1_shared_tag_proto_rawDescData []byte
)

func file_v1_shared_tag_proto_rawDescGZIP() []byte {
	file_v1_shared_tag_proto_rawDescOnce.Do(func() {
		file_v1_shared_tag_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_shared_tag_proto_rawDesc), len(file_v1_shared_tag_proto_rawDesc)))
	})
	return file_v1_shared_tag_proto_rawDescData
}

var file_v1_shared_tag_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_v1_shared_tag_proto_goTypes = []any{
	(*Tag)(nil),                   // 0: shared.v1.Tag
	(TagType)(0),                  // 1: shared.v1.TagType
	(ContextType)(0),              // 2: shared.v1.ContextType
	(ContentRating)(0),            // 3: shared.v1.ContentRating
	(ContentDescriptorType)(0),    // 4: shared.v1.ContentDescriptorType
	(*Metadata)(nil),              // 5: shared.v1.Metadata
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_v1_shared_tag_proto_depIdxs = []int32{
	1, // 0: shared.v1.Tag.type:type_name -> shared.v1.TagType
	2, // 1: shared.v1.Tag.context:type_name -> shared.v1.ContextType
	3, // 2: shared.v1.Tag.content_rating:type_name -> shared.v1.ContentRating
	4, // 3: shared.v1.Tag.content_descriptors:type_name -> shared.v1.ContentDescriptorType
	5, // 4: shared.v1.Tag.metadata:type_name -> shared.v1.Metadata
	6, // 5: shared.v1.Tag.created_at:type_name -> google.protobuf.Timestamp
	6, // 6: shared.v1.Tag.updated_at:type_name -> google.protobuf.Timestamp
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_v1_shared_tag_proto_init() }
func file_v1_shared_tag_proto_init() {
	if File_v1_shared_tag_proto != nil {
		return
	}
	file_v1_shared_tagtype_proto_init()
	file_v1_shared_contentrating_proto_init()
	file_v1_shared_contentdescriptortype_proto_init()
	file_v1_shared_metadata_proto_init()
	file_v1_shared_contexttype_proto_init()
	file_v1_shared_tag_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_shared_tag_proto_rawDesc), len(file_v1_shared_tag_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_shared_tag_proto_goTypes,
		DependencyIndexes: file_v1_shared_tag_proto_depIdxs,
		MessageInfos:      file_v1_shared_tag_proto_msgTypes,
	}.Build()
	File_v1_shared_tag_proto = out.File
	file_v1_shared_tag_proto_goTypes = nil
	file_v1_shared_tag_proto_depIdxs = nil
}
