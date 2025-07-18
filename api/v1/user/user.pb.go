// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: v1/user/user.proto

package userv1

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

type ProfileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProfileRequest) Reset() {
	*x = ProfileRequest{}
	mi := &file_v1_user_user_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileRequest) ProtoMessage() {}

func (x *ProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_user_user_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileRequest.ProtoReflect.Descriptor instead.
func (*ProfileRequest) Descriptor() ([]byte, []int) {
	return file_v1_user_user_proto_rawDescGZIP(), []int{0}
}

type ProfileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *shared.User           `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProfileResponse) Reset() {
	*x = ProfileResponse{}
	mi := &file_v1_user_user_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProfileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileResponse) ProtoMessage() {}

func (x *ProfileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_user_user_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileResponse.ProtoReflect.Descriptor instead.
func (*ProfileResponse) Descriptor() ([]byte, []int) {
	return file_v1_user_user_proto_rawDescGZIP(), []int{1}
}

func (x *ProfileResponse) GetUser() *shared.User {
	if x != nil {
		return x.User
	}
	return nil
}

type UserByIDRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserByIDRequest) Reset() {
	*x = UserByIDRequest{}
	mi := &file_v1_user_user_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserByIDRequest) ProtoMessage() {}

func (x *UserByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_user_user_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserByIDRequest.ProtoReflect.Descriptor instead.
func (*UserByIDRequest) Descriptor() ([]byte, []int) {
	return file_v1_user_user_proto_rawDescGZIP(), []int{2}
}

func (x *UserByIDRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type UserByEmailRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserByEmailRequest) Reset() {
	*x = UserByEmailRequest{}
	mi := &file_v1_user_user_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserByEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserByEmailRequest) ProtoMessage() {}

func (x *UserByEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_user_user_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserByEmailRequest.ProtoReflect.Descriptor instead.
func (*UserByEmailRequest) Descriptor() ([]byte, []int) {
	return file_v1_user_user_proto_rawDescGZIP(), []int{3}
}

func (x *UserByEmailRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type UserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *shared.User           `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserResponse) Reset() {
	*x = UserResponse{}
	mi := &file_v1_user_user_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserResponse) ProtoMessage() {}

func (x *UserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_user_user_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserResponse.ProtoReflect.Descriptor instead.
func (*UserResponse) Descriptor() ([]byte, []int) {
	return file_v1_user_user_proto_rawDescGZIP(), []int{4}
}

func (x *UserResponse) GetUser() *shared.User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_v1_user_user_proto protoreflect.FileDescriptor

const file_v1_user_user_proto_rawDesc = "" +
	"\n" +
	"\x12v1/user/user.proto\x12\auser.v1\x1a\x14v1/shared/user.proto\"\x10\n" +
	"\x0eProfileRequest\"6\n" +
	"\x0fProfileResponse\x12#\n" +
	"\x04user\x18\x01 \x01(\v2\x0f.shared.v1.UserR\x04user\"*\n" +
	"\x0fUserByIDRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\"*\n" +
	"\x12UserByEmailRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\"3\n" +
	"\fUserResponse\x12#\n" +
	"\x04user\x18\x01 \x01(\v2\x0f.shared.v1.UserR\x04user2\xcb\x01\n" +
	"\vUserService\x12<\n" +
	"\aProfile\x12\x17.user.v1.ProfileRequest\x1a\x18.user.v1.ProfileResponse\x12;\n" +
	"\bUserByID\x12\x18.user.v1.UserByIDRequest\x1a\x15.user.v1.UserResponse\x12A\n" +
	"\vUserByEmail\x12\x1b.user.v1.UserByEmailRequest\x1a\x15.user.v1.UserResponseB@Z>github.com/studyguides-com/study-guides-api/api/v1/user;userv1b\x06proto3"

var (
	file_v1_user_user_proto_rawDescOnce sync.Once
	file_v1_user_user_proto_rawDescData []byte
)

func file_v1_user_user_proto_rawDescGZIP() []byte {
	file_v1_user_user_proto_rawDescOnce.Do(func() {
		file_v1_user_user_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_user_user_proto_rawDesc), len(file_v1_user_user_proto_rawDesc)))
	})
	return file_v1_user_user_proto_rawDescData
}

var file_v1_user_user_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_v1_user_user_proto_goTypes = []any{
	(*ProfileRequest)(nil),     // 0: user.v1.ProfileRequest
	(*ProfileResponse)(nil),    // 1: user.v1.ProfileResponse
	(*UserByIDRequest)(nil),    // 2: user.v1.UserByIDRequest
	(*UserByEmailRequest)(nil), // 3: user.v1.UserByEmailRequest
	(*UserResponse)(nil),       // 4: user.v1.UserResponse
	(*shared.User)(nil),        // 5: shared.v1.User
}
var file_v1_user_user_proto_depIdxs = []int32{
	5, // 0: user.v1.ProfileResponse.user:type_name -> shared.v1.User
	5, // 1: user.v1.UserResponse.user:type_name -> shared.v1.User
	0, // 2: user.v1.UserService.Profile:input_type -> user.v1.ProfileRequest
	2, // 3: user.v1.UserService.UserByID:input_type -> user.v1.UserByIDRequest
	3, // 4: user.v1.UserService.UserByEmail:input_type -> user.v1.UserByEmailRequest
	1, // 5: user.v1.UserService.Profile:output_type -> user.v1.ProfileResponse
	4, // 6: user.v1.UserService.UserByID:output_type -> user.v1.UserResponse
	4, // 7: user.v1.UserService.UserByEmail:output_type -> user.v1.UserResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_v1_user_user_proto_init() }
func file_v1_user_user_proto_init() {
	if File_v1_user_user_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_user_user_proto_rawDesc), len(file_v1_user_user_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_user_user_proto_goTypes,
		DependencyIndexes: file_v1_user_user_proto_depIdxs,
		MessageInfos:      file_v1_user_user_proto_msgTypes,
	}.Build()
	File_v1_user_user_proto = out.File
	file_v1_user_user_proto_goTypes = nil
	file_v1_user_user_proto_depIdxs = nil
}
