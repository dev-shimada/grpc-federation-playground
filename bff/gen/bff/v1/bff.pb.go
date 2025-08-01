// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: bff/v1/bff.proto

package bffv1

import (
	_ "github.com/dev-shimada/grpc-federation-playground/bff/gen/message/v1"
	_ "github.com/dev-shimada/grpc-federation-playground/bff/gen/user/v1"
	_ "github.com/mercari/grpc-federation/grpc/federation"
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

type GetMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MessageId     string                 `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMessageRequest) Reset() {
	*x = GetMessageRequest{}
	mi := &file_bff_v1_bff_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessageRequest) ProtoMessage() {}

func (x *GetMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bff_v1_bff_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessageRequest.ProtoReflect.Descriptor instead.
func (*GetMessageRequest) Descriptor() ([]byte, []int) {
	return file_bff_v1_bff_proto_rawDescGZIP(), []int{0}
}

func (x *GetMessageRequest) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

type GetMessageResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       *Message               `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	User          *User                  `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMessageResponse) Reset() {
	*x = GetMessageResponse{}
	mi := &file_bff_v1_bff_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMessageResponse) ProtoMessage() {}

func (x *GetMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_bff_v1_bff_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMessageResponse.ProtoReflect.Descriptor instead.
func (*GetMessageResponse) Descriptor() ([]byte, []int) {
	return file_bff_v1_bff_proto_rawDescGZIP(), []int{1}
}

func (x *GetMessageResponse) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *GetMessageResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Text          string                 `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_bff_v1_bff_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_bff_v1_bff_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_bff_v1_bff_proto_rawDescGZIP(), []int{2}
}

func (x *Message) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Message) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Email         string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     string                 `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_bff_v1_bff_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_bff_v1_bff_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_bff_v1_bff_proto_rawDescGZIP(), []int{3}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *User) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

var File_bff_v1_bff_proto protoreflect.FileDescriptor

const file_bff_v1_bff_proto_rawDesc = "" +
	"\n" +
	"\x10bff/v1/bff.proto\x12\x06bff.v1\x1a grpc/federation/federation.proto\x1a\x12user/v1/user.proto\x1a\x18message/v1/message.proto\"2\n" +
	"\x11GetMessageRequest\x12\x1d\n" +
	"\n" +
	"message_id\x18\x01 \x01(\tR\tmessageId\"\xd0\x01\n" +
	"\x12GetMessageResponse\x127\n" +
	"\amessage\x18\x01 \x01(\v2\x0f.bff.v1.MessageB\f\x9aJ\t\x12\amessageR\amessage\x12+\n" +
	"\x04user\x18\x02 \x01(\v2\f.bff.v1.UserB\t\x9aJ\x06\x12\x04userR\x04user:T\x9aJQ\n" +
	"(\n" +
	"\amessagej\x1d\n" +
	"\aMessage\x12\x12\n" +
	"\x02id\x12\f$.message_id\n" +
	"%\n" +
	"\x04userj\x1d\n" +
	"\x04User\x12\x15\n" +
	"\x02id\x12\x0fmessage.user_id\"q\n" +
	"\aMessage\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x12\n" +
	"\x04text\x18\x02 \x01(\tR\x04text:9\x9aJ6\n" +
	"4\n" +
	"\x03res\x18\x01r+\n" +
	"\x1dmessage.v1.MessageService/Get\x12\n" +
	"\n" +
	"\x02id\x12\x04$.id\"\xb3\x01\n" +
	"\x04User\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" +
	"\x05email\x18\x02 \x01(\tR\x05email\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name\x12\x1d\n" +
	"\n" +
	"created_at\x18\x04 \x01(\tR\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\x05 \x01(\tR\tupdatedAt:3\x9aJ0\n" +
	".\n" +
	"\x03res\x18\x01r%\n" +
	"\x17user.v1.UserService/Get\x12\n" +
	"\n" +
	"\x02id\x12\x04$.id2V\n" +
	"\n" +
	"BffService\x12C\n" +
	"\n" +
	"GetMessage\x12\x19.bff.v1.GetMessageRequest\x1a\x1a.bff.v1.GetMessageResponse\x1a\x03\x9aJ\x00B\xc8\x01\x9aJ.\x12\x12user/v1/user.proto\x12\x18message/v1/message.proto\n" +
	"\n" +
	"com.bff.v1B\bBffProtoP\x01ZFgithub.com/dev-shimada/grpc-federation-playground/bff/gen/bff/v1;bffv1\xa2\x02\x03BXX\xaa\x02\x06Bff.V1\xca\x02\x06Bff\\V1\xe2\x02\x12Bff\\V1\\GPBMetadata\xea\x02\aBff::V1b\x06proto3"

var (
	file_bff_v1_bff_proto_rawDescOnce sync.Once
	file_bff_v1_bff_proto_rawDescData []byte
)

func file_bff_v1_bff_proto_rawDescGZIP() []byte {
	file_bff_v1_bff_proto_rawDescOnce.Do(func() {
		file_bff_v1_bff_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_bff_v1_bff_proto_rawDesc), len(file_bff_v1_bff_proto_rawDesc)))
	})
	return file_bff_v1_bff_proto_rawDescData
}

var file_bff_v1_bff_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_bff_v1_bff_proto_goTypes = []any{
	(*GetMessageRequest)(nil),  // 0: bff.v1.GetMessageRequest
	(*GetMessageResponse)(nil), // 1: bff.v1.GetMessageResponse
	(*Message)(nil),            // 2: bff.v1.Message
	(*User)(nil),               // 3: bff.v1.User
}
var file_bff_v1_bff_proto_depIdxs = []int32{
	2, // 0: bff.v1.GetMessageResponse.message:type_name -> bff.v1.Message
	3, // 1: bff.v1.GetMessageResponse.user:type_name -> bff.v1.User
	0, // 2: bff.v1.BffService.GetMessage:input_type -> bff.v1.GetMessageRequest
	1, // 3: bff.v1.BffService.GetMessage:output_type -> bff.v1.GetMessageResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_bff_v1_bff_proto_init() }
func file_bff_v1_bff_proto_init() {
	if File_bff_v1_bff_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_bff_v1_bff_proto_rawDesc), len(file_bff_v1_bff_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_bff_v1_bff_proto_goTypes,
		DependencyIndexes: file_bff_v1_bff_proto_depIdxs,
		MessageInfos:      file_bff_v1_bff_proto_msgTypes,
	}.Build()
	File_bff_v1_bff_proto = out.File
	file_bff_v1_bff_proto_goTypes = nil
	file_bff_v1_bff_proto_depIdxs = nil
}
