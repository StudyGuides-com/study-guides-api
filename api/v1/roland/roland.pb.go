// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: v1/roland/roland.proto

package rolandv1

import (
	shared "github.com/studyguides-com/study-guides-api/api/v1/shared"
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

// SaveBundleRequest represents a request to save a bundle
type SaveBundleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Bundle        *shared.Bundle         `protobuf:"bytes,1,opt,name=bundle,proto3" json:"bundle,omitempty"`
	Force         bool                   `protobuf:"varint,2,opt,name=force,proto3" json:"force,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SaveBundleRequest) Reset() {
	*x = SaveBundleRequest{}
	mi := &file_v1_roland_roland_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SaveBundleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBundleRequest) ProtoMessage() {}

func (x *SaveBundleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBundleRequest.ProtoReflect.Descriptor instead.
func (*SaveBundleRequest) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{0}
}

func (x *SaveBundleRequest) GetBundle() *shared.Bundle {
	if x != nil {
		return x.Bundle
	}
	return nil
}

func (x *SaveBundleRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

// SaveBundleResponse represents the response from saving a bundle
type SaveBundleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Created       bool                   `protobuf:"varint,1,opt,name=created,proto3" json:"created,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SaveBundleResponse) Reset() {
	*x = SaveBundleResponse{}
	mi := &file_v1_roland_roland_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SaveBundleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBundleResponse) ProtoMessage() {}

func (x *SaveBundleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBundleResponse.ProtoReflect.Descriptor instead.
func (*SaveBundleResponse) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{1}
}

func (x *SaveBundleResponse) GetCreated() bool {
	if x != nil {
		return x.Created
	}
	return false
}

// BundlesRequest represents a request to get all bundles
type BundlesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BundlesRequest) Reset() {
	*x = BundlesRequest{}
	mi := &file_v1_roland_roland_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BundlesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BundlesRequest) ProtoMessage() {}

func (x *BundlesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BundlesRequest.ProtoReflect.Descriptor instead.
func (*BundlesRequest) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{2}
}

// BundlesResponse represents the response containing all bundles
type BundlesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Bundles       []*shared.Bundle       `protobuf:"bytes,1,rep,name=bundles,proto3" json:"bundles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BundlesResponse) Reset() {
	*x = BundlesResponse{}
	mi := &file_v1_roland_roland_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BundlesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BundlesResponse) ProtoMessage() {}

func (x *BundlesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BundlesResponse.ProtoReflect.Descriptor instead.
func (*BundlesResponse) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{3}
}

func (x *BundlesResponse) GetBundles() []*shared.Bundle {
	if x != nil {
		return x.Bundles
	}
	return nil
}

// BundlesByParserTypeRequest represents a request to get bundles by parser type
type BundlesByParserTypeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ParserType    shared.ParserType      `protobuf:"varint,1,opt,name=parser_type,json=parserType,proto3,enum=shared.v1.ParserType" json:"parser_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BundlesByParserTypeRequest) Reset() {
	*x = BundlesByParserTypeRequest{}
	mi := &file_v1_roland_roland_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BundlesByParserTypeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BundlesByParserTypeRequest) ProtoMessage() {}

func (x *BundlesByParserTypeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BundlesByParserTypeRequest.ProtoReflect.Descriptor instead.
func (*BundlesByParserTypeRequest) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{4}
}

func (x *BundlesByParserTypeRequest) GetParserType() shared.ParserType {
	if x != nil {
		return x.ParserType
	}
	return shared.ParserType(0)
}

// BundlesByParserTypeResponse represents the response containing bundles by parser type
type BundlesByParserTypeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Bundles       []*shared.Bundle       `protobuf:"bytes,1,rep,name=bundles,proto3" json:"bundles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BundlesByParserTypeResponse) Reset() {
	*x = BundlesByParserTypeResponse{}
	mi := &file_v1_roland_roland_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BundlesByParserTypeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BundlesByParserTypeResponse) ProtoMessage() {}

func (x *BundlesByParserTypeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BundlesByParserTypeResponse.ProtoReflect.Descriptor instead.
func (*BundlesByParserTypeResponse) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{5}
}

func (x *BundlesByParserTypeResponse) GetBundles() []*shared.Bundle {
	if x != nil {
		return x.Bundles
	}
	return nil
}

// UpdateGobRequest represents a request to update the gob payload for a bundle
type UpdateGobRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	GobPayload    []byte                 `protobuf:"bytes,2,opt,name=gob_payload,json=gobPayload,proto3" json:"gob_payload,omitempty"`
	Force         bool                   `protobuf:"varint,3,opt,name=force,proto3" json:"force,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateGobRequest) Reset() {
	*x = UpdateGobRequest{}
	mi := &file_v1_roland_roland_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateGobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateGobRequest) ProtoMessage() {}

func (x *UpdateGobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateGobRequest.ProtoReflect.Descriptor instead.
func (*UpdateGobRequest) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateGobRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateGobRequest) GetGobPayload() []byte {
	if x != nil {
		return x.GobPayload
	}
	return nil
}

func (x *UpdateGobRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

// UpdateGobResponse represents the response from updating gob payload
type UpdateGobResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Updated       bool                   `protobuf:"varint,1,opt,name=updated,proto3" json:"updated,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateGobResponse) Reset() {
	*x = UpdateGobResponse{}
	mi := &file_v1_roland_roland_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateGobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateGobResponse) ProtoMessage() {}

func (x *UpdateGobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateGobResponse.ProtoReflect.Descriptor instead.
func (*UpdateGobResponse) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateGobResponse) GetUpdated() bool {
	if x != nil {
		return x.Updated
	}
	return false
}

// DeleteAllBundlesRequest represents a request to delete all bundles
type DeleteAllBundlesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteAllBundlesRequest) Reset() {
	*x = DeleteAllBundlesRequest{}
	mi := &file_v1_roland_roland_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteAllBundlesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAllBundlesRequest) ProtoMessage() {}

func (x *DeleteAllBundlesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAllBundlesRequest.ProtoReflect.Descriptor instead.
func (*DeleteAllBundlesRequest) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{8}
}

// DeleteAllBundlesResponse represents the response from deleting all bundles
type DeleteAllBundlesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteAllBundlesResponse) Reset() {
	*x = DeleteAllBundlesResponse{}
	mi := &file_v1_roland_roland_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteAllBundlesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAllBundlesResponse) ProtoMessage() {}

func (x *DeleteAllBundlesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAllBundlesResponse.ProtoReflect.Descriptor instead.
func (*DeleteAllBundlesResponse) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteAllBundlesResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// DeleteBundleByIDRequest represents a request to delete a bundle by ID
type DeleteBundleByIDRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteBundleByIDRequest) Reset() {
	*x = DeleteBundleByIDRequest{}
	mi := &file_v1_roland_roland_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteBundleByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteBundleByIDRequest) ProtoMessage() {}

func (x *DeleteBundleByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteBundleByIDRequest.ProtoReflect.Descriptor instead.
func (*DeleteBundleByIDRequest) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{10}
}

func (x *DeleteBundleByIDRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// DeleteBundleByIDResponse represents the response from deleting a bundle by ID
type DeleteBundleByIDResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteBundleByIDResponse) Reset() {
	*x = DeleteBundleByIDResponse{}
	mi := &file_v1_roland_roland_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteBundleByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteBundleByIDResponse) ProtoMessage() {}

func (x *DeleteBundleByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteBundleByIDResponse.ProtoReflect.Descriptor instead.
func (*DeleteBundleByIDResponse) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{11}
}

func (x *DeleteBundleByIDResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// DeleteBundlesByShortIDRequest represents a request to delete bundles by short ID
type DeleteBundlesByShortIDRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortId       string                 `protobuf:"bytes,1,opt,name=short_id,json=shortId,proto3" json:"short_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteBundlesByShortIDRequest) Reset() {
	*x = DeleteBundlesByShortIDRequest{}
	mi := &file_v1_roland_roland_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteBundlesByShortIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteBundlesByShortIDRequest) ProtoMessage() {}

func (x *DeleteBundlesByShortIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteBundlesByShortIDRequest.ProtoReflect.Descriptor instead.
func (*DeleteBundlesByShortIDRequest) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{12}
}

func (x *DeleteBundlesByShortIDRequest) GetShortId() string {
	if x != nil {
		return x.ShortId
	}
	return ""
}

// DeleteBundlesByShortIDResponse represents the response from deleting bundles by short ID
type DeleteBundlesByShortIDResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DeletedCount  int32                  `protobuf:"varint,1,opt,name=deleted_count,json=deletedCount,proto3" json:"deleted_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteBundlesByShortIDResponse) Reset() {
	*x = DeleteBundlesByShortIDResponse{}
	mi := &file_v1_roland_roland_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteBundlesByShortIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteBundlesByShortIDResponse) ProtoMessage() {}

func (x *DeleteBundlesByShortIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteBundlesByShortIDResponse.ProtoReflect.Descriptor instead.
func (*DeleteBundlesByShortIDResponse) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{13}
}

func (x *DeleteBundlesByShortIDResponse) GetDeletedCount() int32 {
	if x != nil {
		return x.DeletedCount
	}
	return 0
}

// MarkBundleExportedRequest represents a request to mark a bundle as exported
type MarkBundleExportedRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ExportType    shared.ExportType      `protobuf:"varint,2,opt,name=export_type,json=exportType,proto3,enum=shared.v1.ExportType" json:"export_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MarkBundleExportedRequest) Reset() {
	*x = MarkBundleExportedRequest{}
	mi := &file_v1_roland_roland_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MarkBundleExportedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkBundleExportedRequest) ProtoMessage() {}

func (x *MarkBundleExportedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkBundleExportedRequest.ProtoReflect.Descriptor instead.
func (*MarkBundleExportedRequest) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{14}
}

func (x *MarkBundleExportedRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MarkBundleExportedRequest) GetExportType() shared.ExportType {
	if x != nil {
		return x.ExportType
	}
	return shared.ExportType(0)
}

// MarkBundleExportedResponse represents the response from marking a bundle as exported
type MarkBundleExportedResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Updated       bool                   `protobuf:"varint,1,opt,name=updated,proto3" json:"updated,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MarkBundleExportedResponse) Reset() {
	*x = MarkBundleExportedResponse{}
	mi := &file_v1_roland_roland_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MarkBundleExportedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkBundleExportedResponse) ProtoMessage() {}

func (x *MarkBundleExportedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_roland_roland_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkBundleExportedResponse.ProtoReflect.Descriptor instead.
func (*MarkBundleExportedResponse) Descriptor() ([]byte, []int) {
	return file_v1_roland_roland_proto_rawDescGZIP(), []int{15}
}

func (x *MarkBundleExportedResponse) GetUpdated() bool {
	if x != nil {
		return x.Updated
	}
	return false
}

var File_v1_roland_roland_proto protoreflect.FileDescriptor

const file_v1_roland_roland_proto_rawDesc = "" +
	"\n" +
	"\x16v1/roland/roland.proto\x12\troland.v1\x1a\x16v1/shared/bundle.proto\x1a\x1av1/shared/exporttype.proto\x1a\x1av1/shared/parsertype.proto\"T\n" +
	"\x11SaveBundleRequest\x12)\n" +
	"\x06bundle\x18\x01 \x01(\v2\x11.shared.v1.BundleR\x06bundle\x12\x14\n" +
	"\x05force\x18\x02 \x01(\bR\x05force\".\n" +
	"\x12SaveBundleResponse\x12\x18\n" +
	"\acreated\x18\x01 \x01(\bR\acreated\"\x10\n" +
	"\x0eBundlesRequest\">\n" +
	"\x0fBundlesResponse\x12+\n" +
	"\abundles\x18\x01 \x03(\v2\x11.shared.v1.BundleR\abundles\"T\n" +
	"\x1aBundlesByParserTypeRequest\x126\n" +
	"\vparser_type\x18\x01 \x01(\x0e2\x15.shared.v1.ParserTypeR\n" +
	"parserType\"J\n" +
	"\x1bBundlesByParserTypeResponse\x12+\n" +
	"\abundles\x18\x01 \x03(\v2\x11.shared.v1.BundleR\abundles\"Y\n" +
	"\x10UpdateGobRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1f\n" +
	"\vgob_payload\x18\x02 \x01(\fR\n" +
	"gobPayload\x12\x14\n" +
	"\x05force\x18\x03 \x01(\bR\x05force\"-\n" +
	"\x11UpdateGobResponse\x12\x18\n" +
	"\aupdated\x18\x01 \x01(\bR\aupdated\"\x19\n" +
	"\x17DeleteAllBundlesRequest\"4\n" +
	"\x18DeleteAllBundlesResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\")\n" +
	"\x17DeleteBundleByIDRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"4\n" +
	"\x18DeleteBundleByIDResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\":\n" +
	"\x1dDeleteBundlesByShortIDRequest\x12\x19\n" +
	"\bshort_id\x18\x01 \x01(\tR\ashortId\"E\n" +
	"\x1eDeleteBundlesByShortIDResponse\x12#\n" +
	"\rdeleted_count\x18\x01 \x01(\x05R\fdeletedCount\"c\n" +
	"\x19MarkBundleExportedRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x126\n" +
	"\vexport_type\x18\x02 \x01(\x0e2\x15.shared.v1.ExportTypeR\n" +
	"exportType\"6\n" +
	"\x1aMarkBundleExportedResponse\x12\x18\n" +
	"\aupdated\x18\x01 \x01(\bR\aupdated2\xd6\x05\n" +
	"\rRolandService\x12I\n" +
	"\n" +
	"SaveBundle\x12\x1c.roland.v1.SaveBundleRequest\x1a\x1d.roland.v1.SaveBundleResponse\x12@\n" +
	"\aBundles\x12\x19.roland.v1.BundlesRequest\x1a\x1a.roland.v1.BundlesResponse\x12d\n" +
	"\x13BundlesByParserType\x12%.roland.v1.BundlesByParserTypeRequest\x1a&.roland.v1.BundlesByParserTypeResponse\x12F\n" +
	"\tUpdateGob\x12\x1b.roland.v1.UpdateGobRequest\x1a\x1c.roland.v1.UpdateGobResponse\x12[\n" +
	"\x10DeleteAllBundles\x12\".roland.v1.DeleteAllBundlesRequest\x1a#.roland.v1.DeleteAllBundlesResponse\x12[\n" +
	"\x10DeleteBundleByID\x12\".roland.v1.DeleteBundleByIDRequest\x1a#.roland.v1.DeleteBundleByIDResponse\x12m\n" +
	"\x16DeleteBundlesByShortID\x12(.roland.v1.DeleteBundlesByShortIDRequest\x1a).roland.v1.DeleteBundlesByShortIDResponse\x12a\n" +
	"\x12MarkBundleExported\x12$.roland.v1.MarkBundleExportedRequest\x1a%.roland.v1.MarkBundleExportedResponseBDZBgithub.com.studyguides-com/study-guides-api/api/v1/roland;rolandv1b\x06proto3"

var (
	file_v1_roland_roland_proto_rawDescOnce sync.Once
	file_v1_roland_roland_proto_rawDescData []byte
)

func file_v1_roland_roland_proto_rawDescGZIP() []byte {
	file_v1_roland_roland_proto_rawDescOnce.Do(func() {
		file_v1_roland_roland_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_roland_roland_proto_rawDesc), len(file_v1_roland_roland_proto_rawDesc)))
	})
	return file_v1_roland_roland_proto_rawDescData
}

var file_v1_roland_roland_proto_msgTypes = make([]protoimpl.MessageInfo, 16)
var file_v1_roland_roland_proto_goTypes = []any{
	(*SaveBundleRequest)(nil),              // 0: roland.v1.SaveBundleRequest
	(*SaveBundleResponse)(nil),             // 1: roland.v1.SaveBundleResponse
	(*BundlesRequest)(nil),                 // 2: roland.v1.BundlesRequest
	(*BundlesResponse)(nil),                // 3: roland.v1.BundlesResponse
	(*BundlesByParserTypeRequest)(nil),     // 4: roland.v1.BundlesByParserTypeRequest
	(*BundlesByParserTypeResponse)(nil),    // 5: roland.v1.BundlesByParserTypeResponse
	(*UpdateGobRequest)(nil),               // 6: roland.v1.UpdateGobRequest
	(*UpdateGobResponse)(nil),              // 7: roland.v1.UpdateGobResponse
	(*DeleteAllBundlesRequest)(nil),        // 8: roland.v1.DeleteAllBundlesRequest
	(*DeleteAllBundlesResponse)(nil),       // 9: roland.v1.DeleteAllBundlesResponse
	(*DeleteBundleByIDRequest)(nil),        // 10: roland.v1.DeleteBundleByIDRequest
	(*DeleteBundleByIDResponse)(nil),       // 11: roland.v1.DeleteBundleByIDResponse
	(*DeleteBundlesByShortIDRequest)(nil),  // 12: roland.v1.DeleteBundlesByShortIDRequest
	(*DeleteBundlesByShortIDResponse)(nil), // 13: roland.v1.DeleteBundlesByShortIDResponse
	(*MarkBundleExportedRequest)(nil),      // 14: roland.v1.MarkBundleExportedRequest
	(*MarkBundleExportedResponse)(nil),     // 15: roland.v1.MarkBundleExportedResponse
	(*shared.Bundle)(nil),                  // 16: shared.v1.Bundle
	(shared.ParserType)(0),                 // 17: shared.v1.ParserType
	(shared.ExportType)(0),                 // 18: shared.v1.ExportType
}
var file_v1_roland_roland_proto_depIdxs = []int32{
	16, // 0: roland.v1.SaveBundleRequest.bundle:type_name -> shared.v1.Bundle
	16, // 1: roland.v1.BundlesResponse.bundles:type_name -> shared.v1.Bundle
	17, // 2: roland.v1.BundlesByParserTypeRequest.parser_type:type_name -> shared.v1.ParserType
	16, // 3: roland.v1.BundlesByParserTypeResponse.bundles:type_name -> shared.v1.Bundle
	18, // 4: roland.v1.MarkBundleExportedRequest.export_type:type_name -> shared.v1.ExportType
	0,  // 5: roland.v1.RolandService.SaveBundle:input_type -> roland.v1.SaveBundleRequest
	2,  // 6: roland.v1.RolandService.Bundles:input_type -> roland.v1.BundlesRequest
	4,  // 7: roland.v1.RolandService.BundlesByParserType:input_type -> roland.v1.BundlesByParserTypeRequest
	6,  // 8: roland.v1.RolandService.UpdateGob:input_type -> roland.v1.UpdateGobRequest
	8,  // 9: roland.v1.RolandService.DeleteAllBundles:input_type -> roland.v1.DeleteAllBundlesRequest
	10, // 10: roland.v1.RolandService.DeleteBundleByID:input_type -> roland.v1.DeleteBundleByIDRequest
	12, // 11: roland.v1.RolandService.DeleteBundlesByShortID:input_type -> roland.v1.DeleteBundlesByShortIDRequest
	14, // 12: roland.v1.RolandService.MarkBundleExported:input_type -> roland.v1.MarkBundleExportedRequest
	1,  // 13: roland.v1.RolandService.SaveBundle:output_type -> roland.v1.SaveBundleResponse
	3,  // 14: roland.v1.RolandService.Bundles:output_type -> roland.v1.BundlesResponse
	5,  // 15: roland.v1.RolandService.BundlesByParserType:output_type -> roland.v1.BundlesByParserTypeResponse
	7,  // 16: roland.v1.RolandService.UpdateGob:output_type -> roland.v1.UpdateGobResponse
	9,  // 17: roland.v1.RolandService.DeleteAllBundles:output_type -> roland.v1.DeleteAllBundlesResponse
	11, // 18: roland.v1.RolandService.DeleteBundleByID:output_type -> roland.v1.DeleteBundleByIDResponse
	13, // 19: roland.v1.RolandService.DeleteBundlesByShortID:output_type -> roland.v1.DeleteBundlesByShortIDResponse
	15, // 20: roland.v1.RolandService.MarkBundleExported:output_type -> roland.v1.MarkBundleExportedResponse
	13, // [13:21] is the sub-list for method output_type
	5,  // [5:13] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_v1_roland_roland_proto_init() }
func file_v1_roland_roland_proto_init() {
	if File_v1_roland_roland_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_roland_roland_proto_rawDesc), len(file_v1_roland_roland_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   16,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_roland_roland_proto_goTypes,
		DependencyIndexes: file_v1_roland_roland_proto_depIdxs,
		MessageInfos:      file_v1_roland_roland_proto_msgTypes,
	}.Build()
	File_v1_roland_roland_proto = out.File
	file_v1_roland_roland_proto_goTypes = nil
	file_v1_roland_roland_proto_depIdxs = nil
}
