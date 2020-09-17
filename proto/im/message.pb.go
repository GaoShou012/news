// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: proto/im/message.proto

package im

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

type Conn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Uuid string `protobuf:"bytes,2,opt,name=Uuid,proto3" json:"Uuid,omitempty"`
}

func (x *Conn) Reset() {
	*x = Conn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_im_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Conn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Conn) ProtoMessage() {}

func (x *Conn) ProtoReflect() protoreflect.Message {
	mi := &file_proto_im_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Conn.ProtoReflect.Descriptor instead.
func (*Conn) Descriptor() ([]byte, []int) {
	return file_proto_im_message_proto_rawDescGZIP(), []int{0}
}

func (x *Conn) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Conn) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code uint64 `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Desc string `protobuf:"bytes,2,opt,name=Desc,proto3" json:"Desc,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_im_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_proto_im_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_proto_im_message_proto_rawDescGZIP(), []int{1}
}

func (x *Error) GetCode() uint64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Error) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code uint64 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Desc string `protobuf:"bytes,2,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_im_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_proto_im_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_proto_im_message_proto_rawDescGZIP(), []int{2}
}

func (x *State) GetCode() uint64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *State) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

type Head struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=Path,proto3" json:"Path,omitempty"`
}

func (x *Head) Reset() {
	*x = Head{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_im_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Head) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Head) ProtoMessage() {}

func (x *Head) ProtoReflect() protoreflect.Message {
	mi := &file_proto_im_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Head.ProtoReflect.Descriptor instead.
func (*Head) Descriptor() ([]byte, []int) {
	return file_proto_im_message_proto_rawDescGZIP(), []int{3}
}

func (x *Head) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Head *Head  `protobuf:"bytes,1,opt,name=Head,proto3" json:"Head,omitempty"`
	Body []byte `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_im_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_proto_im_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_proto_im_message_proto_rawDescGZIP(), []int{4}
}

func (x *Message) GetHead() *Head {
	if x != nil {
		return x.Head
	}
	return nil
}

func (x *Message) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_proto_im_message_proto protoreflect.FileDescriptor

var file_proto_im_message_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6d, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x69, 0x6d, 0x22, 0x2a, 0x0a, 0x04, 0x43, 0x6f, 0x6e, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x75,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x75, 0x69, 0x64, 0x22, 0x2f,
	0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x44,
	0x65, 0x73, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x44, 0x65, 0x73, 0x63, 0x22,
	0x2f, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x65, 0x73, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63,
	0x22, 0x1a, 0x0a, 0x04, 0x48, 0x65, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x50, 0x61, 0x74, 0x68, 0x22, 0x41, 0x0a, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x48, 0x65, 0x61, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x69, 0x6d,
	0x2e, 0x48, 0x65, 0x61, 0x64, 0x52, 0x04, 0x48, 0x65, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x42,
	0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x42,
	0x15, 0x5a, 0x13, 0x77, 0x63, 0x68, 0x61, 0x74, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x69, 0x6d, 0x3b, 0x69, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_im_message_proto_rawDescOnce sync.Once
	file_proto_im_message_proto_rawDescData = file_proto_im_message_proto_rawDesc
)

func file_proto_im_message_proto_rawDescGZIP() []byte {
	file_proto_im_message_proto_rawDescOnce.Do(func() {
		file_proto_im_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_im_message_proto_rawDescData)
	})
	return file_proto_im_message_proto_rawDescData
}

var file_proto_im_message_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_im_message_proto_goTypes = []interface{}{
	(*Conn)(nil),    // 0: proto.im.Conn
	(*Error)(nil),   // 1: proto.im.Error
	(*State)(nil),   // 2: proto.im.State
	(*Head)(nil),    // 3: proto.im.Head
	(*Message)(nil), // 4: proto.im.Message
}
var file_proto_im_message_proto_depIdxs = []int32{
	3, // 0: proto.im.Message.Head:type_name -> proto.im.Head
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_im_message_proto_init() }
func file_proto_im_message_proto_init() {
	if File_proto_im_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_im_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Conn); i {
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
		file_proto_im_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_proto_im_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*State); i {
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
		file_proto_im_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Head); i {
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
		file_proto_im_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			RawDescriptor: file_proto_im_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_im_message_proto_goTypes,
		DependencyIndexes: file_proto_im_message_proto_depIdxs,
		MessageInfos:      file_proto_im_message_proto_msgTypes,
	}.Build()
	File_proto_im_message_proto = out.File
	file_proto_im_message_proto_rawDesc = nil
	file_proto_im_message_proto_goTypes = nil
	file_proto_im_message_proto_depIdxs = nil
}