// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: v1/question/question.proto

package questionv1

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

type ForTagRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TagId         string                 `protobuf:"bytes,1,opt,name=tag_id,json=tagId,proto3" json:"tag_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ForTagRequest) Reset() {
	*x = ForTagRequest{}
	mi := &file_v1_question_question_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ForTagRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForTagRequest) ProtoMessage() {}

func (x *ForTagRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_question_question_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForTagRequest.ProtoReflect.Descriptor instead.
func (*ForTagRequest) Descriptor() ([]byte, []int) {
	return file_v1_question_question_proto_rawDescGZIP(), []int{0}
}

func (x *ForTagRequest) GetTagId() string {
	if x != nil {
		return x.TagId
	}
	return ""
}

type QuestionsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Questions     []*shared.Question     `protobuf:"bytes,1,rep,name=questions,proto3" json:"questions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QuestionsResponse) Reset() {
	*x = QuestionsResponse{}
	mi := &file_v1_question_question_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QuestionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuestionsResponse) ProtoMessage() {}

func (x *QuestionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_question_question_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuestionsResponse.ProtoReflect.Descriptor instead.
func (*QuestionsResponse) Descriptor() ([]byte, []int) {
	return file_v1_question_question_proto_rawDescGZIP(), []int{1}
}

func (x *QuestionsResponse) GetQuestions() []*shared.Question {
	if x != nil {
		return x.Questions
	}
	return nil
}

type QuestionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Question      *shared.Question       `protobuf:"bytes,1,opt,name=question,proto3" json:"question,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QuestionResponse) Reset() {
	*x = QuestionResponse{}
	mi := &file_v1_question_question_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QuestionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuestionResponse) ProtoMessage() {}

func (x *QuestionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_question_question_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuestionResponse.ProtoReflect.Descriptor instead.
func (*QuestionResponse) Descriptor() ([]byte, []int) {
	return file_v1_question_question_proto_rawDescGZIP(), []int{2}
}

func (x *QuestionResponse) GetQuestion() *shared.Question {
	if x != nil {
		return x.Question
	}
	return nil
}

type ReportQuestionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	QuestionId    string                 `protobuf:"bytes,1,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
	ReportType    shared.ReportType      `protobuf:"varint,2,opt,name=report_type,json=reportType,proto3,enum=shared.v1.ReportType" json:"report_type,omitempty"`
	Reason        string                 `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReportQuestionRequest) Reset() {
	*x = ReportQuestionRequest{}
	mi := &file_v1_question_question_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReportQuestionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportQuestionRequest) ProtoMessage() {}

func (x *ReportQuestionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_question_question_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportQuestionRequest.ProtoReflect.Descriptor instead.
func (*ReportQuestionRequest) Descriptor() ([]byte, []int) {
	return file_v1_question_question_proto_rawDescGZIP(), []int{3}
}

func (x *ReportQuestionRequest) GetQuestionId() string {
	if x != nil {
		return x.QuestionId
	}
	return ""
}

func (x *ReportQuestionRequest) GetReportType() shared.ReportType {
	if x != nil {
		return x.ReportType
	}
	return shared.ReportType(0)
}

func (x *ReportQuestionRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type ReportQuestionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReportQuestionResponse) Reset() {
	*x = ReportQuestionResponse{}
	mi := &file_v1_question_question_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReportQuestionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportQuestionResponse) ProtoMessage() {}

func (x *ReportQuestionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_question_question_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportQuestionResponse.ProtoReflect.Descriptor instead.
func (*ReportQuestionResponse) Descriptor() ([]byte, []int) {
	return file_v1_question_question_proto_rawDescGZIP(), []int{4}
}

func (x *ReportQuestionResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_v1_question_question_proto protoreflect.FileDescriptor

const file_v1_question_question_proto_rawDesc = "" +
	"\n" +
	"\x1av1/question/question.proto\x12\vquestion.v1\x1a\x18v1/shared/question.proto\x1a\x1av1/shared/reporttype.proto\"&\n" +
	"\rForTagRequest\x12\x15\n" +
	"\x06tag_id\x18\x01 \x01(\tR\x05tagId\"F\n" +
	"\x11QuestionsResponse\x121\n" +
	"\tquestions\x18\x01 \x03(\v2\x13.shared.v1.QuestionR\tquestions\"C\n" +
	"\x10QuestionResponse\x12/\n" +
	"\bquestion\x18\x01 \x01(\v2\x13.shared.v1.QuestionR\bquestion\"\x88\x01\n" +
	"\x15ReportQuestionRequest\x12\x1f\n" +
	"\vquestion_id\x18\x01 \x01(\tR\n" +
	"questionId\x126\n" +
	"\vreport_type\x18\x02 \x01(\x0e2\x15.shared.v1.ReportTypeR\n" +
	"reportType\x12\x16\n" +
	"\x06reason\x18\x03 \x01(\tR\x06reason\"2\n" +
	"\x16ReportQuestionResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess2\xaa\x01\n" +
	"\x0fQuestionService\x12D\n" +
	"\x06ForTag\x12\x1a.question.v1.ForTagRequest\x1a\x1e.question.v1.QuestionsResponse\x12Q\n" +
	"\x06Report\x12\".question.v1.ReportQuestionRequest\x1a#.question.v1.ReportQuestionResponseBHZFgithub.com/studyguides-com/study-guides-api/api/v1/question;questionv1b\x06proto3"

var (
	file_v1_question_question_proto_rawDescOnce sync.Once
	file_v1_question_question_proto_rawDescData []byte
)

func file_v1_question_question_proto_rawDescGZIP() []byte {
	file_v1_question_question_proto_rawDescOnce.Do(func() {
		file_v1_question_question_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_question_question_proto_rawDesc), len(file_v1_question_question_proto_rawDesc)))
	})
	return file_v1_question_question_proto_rawDescData
}

var file_v1_question_question_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_v1_question_question_proto_goTypes = []any{
	(*ForTagRequest)(nil),          // 0: question.v1.ForTagRequest
	(*QuestionsResponse)(nil),      // 1: question.v1.QuestionsResponse
	(*QuestionResponse)(nil),       // 2: question.v1.QuestionResponse
	(*ReportQuestionRequest)(nil),  // 3: question.v1.ReportQuestionRequest
	(*ReportQuestionResponse)(nil), // 4: question.v1.ReportQuestionResponse
	(*shared.Question)(nil),        // 5: shared.v1.Question
	(shared.ReportType)(0),         // 6: shared.v1.ReportType
}
var file_v1_question_question_proto_depIdxs = []int32{
	5, // 0: question.v1.QuestionsResponse.questions:type_name -> shared.v1.Question
	5, // 1: question.v1.QuestionResponse.question:type_name -> shared.v1.Question
	6, // 2: question.v1.ReportQuestionRequest.report_type:type_name -> shared.v1.ReportType
	0, // 3: question.v1.QuestionService.ForTag:input_type -> question.v1.ForTagRequest
	3, // 4: question.v1.QuestionService.Report:input_type -> question.v1.ReportQuestionRequest
	1, // 5: question.v1.QuestionService.ForTag:output_type -> question.v1.QuestionsResponse
	4, // 6: question.v1.QuestionService.Report:output_type -> question.v1.ReportQuestionResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_v1_question_question_proto_init() }
func file_v1_question_question_proto_init() {
	if File_v1_question_question_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_question_question_proto_rawDesc), len(file_v1_question_question_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_question_question_proto_goTypes,
		DependencyIndexes: file_v1_question_question_proto_depIdxs,
		MessageInfos:      file_v1_question_question_proto_msgTypes,
	}.Build()
	File_v1_question_question_proto = out.File
	file_v1_question_question_proto_goTypes = nil
	file_v1_question_question_proto_depIdxs = nil
}
