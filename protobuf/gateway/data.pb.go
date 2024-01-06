// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: protobuf/gateway/data.proto

package gateway

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

type DataMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Flag      int32  `protobuf:"varint,1,opt,name=flag,proto3" json:"flag,omitempty"` // 功能Flag 0: 心跳包  1: 主动读取数据 2: 推送数据  3: 错误信息  4: 设置
	MessageId string `protobuf:"bytes,2,opt,name=messageId,proto3" json:"messageId,omitempty"`
	AgentId   string `protobuf:"bytes,3,opt,name=agentId,proto3" json:"agentId,omitempty"`
	Data      []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"` // 二进制数据字段
}

func (x *DataMessage) Reset() {
	*x = DataMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_gateway_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataMessage) ProtoMessage() {}

func (x *DataMessage) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_gateway_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataMessage.ProtoReflect.Descriptor instead.
func (*DataMessage) Descriptor() ([]byte, []int) {
	return file_protobuf_gateway_data_proto_rawDescGZIP(), []int{0}
}

func (x *DataMessage) GetFlag() int32 {
	if x != nil {
		return x.Flag
	}
	return 0
}

func (x *DataMessage) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *DataMessage) GetAgentId() string {
	if x != nil {
		return x.AgentId
	}
	return ""
}

func (x *DataMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type MessageIdInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId string `protobuf:"bytes,1,opt,name=messageId,proto3" json:"messageId,omitempty"`
	Time      string `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	AgentId   string `protobuf:"bytes,3,opt,name=agentId,proto3" json:"agentId,omitempty"`
}

func (x *MessageIdInfo) Reset() {
	*x = MessageIdInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_gateway_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageIdInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageIdInfo) ProtoMessage() {}

func (x *MessageIdInfo) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_gateway_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageIdInfo.ProtoReflect.Descriptor instead.
func (*MessageIdInfo) Descriptor() ([]byte, []int) {
	return file_protobuf_gateway_data_proto_rawDescGZIP(), []int{1}
}

func (x *MessageIdInfo) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *MessageIdInfo) GetTime() string {
	if x != nil {
		return x.Time
	}
	return ""
}

func (x *MessageIdInfo) GetAgentId() string {
	if x != nil {
		return x.AgentId
	}
	return ""
}

type PingPong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Flag int32 `protobuf:"varint,1,opt,name=flag,proto3" json:"flag,omitempty"` // 0:Ping 1:Pong
}

func (x *PingPong) Reset() {
	*x = PingPong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_gateway_data_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingPong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingPong) ProtoMessage() {}

func (x *PingPong) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_gateway_data_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingPong.ProtoReflect.Descriptor instead.
func (*PingPong) Descriptor() ([]byte, []int) {
	return file_protobuf_gateway_data_proto_rawDescGZIP(), []int{2}
}

func (x *PingPong) GetFlag() int32 {
	if x != nil {
		return x.Flag
	}
	return 0
}

var File_protobuf_gateway_data_proto protoreflect.FileDescriptor

var file_protobuf_gateway_data_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6c,
	0x75, 0x62, 0x72, 0x69, 0x63, 0x61, 0x6e, 0x74, 0x22, 0x6d, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x5b, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x49, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x22, 0x1e, 0x0a, 0x08, 0x50, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x6e, 0x67,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x66, 0x6c, 0x61, 0x67, 0x32, 0xc9, 0x01, 0x0a, 0x0e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3c, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x16, 0x2e, 0x6c, 0x75, 0x62, 0x72, 0x69, 0x63, 0x61, 0x6e, 0x74, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x6c, 0x75, 0x62, 0x72, 0x69, 0x63,
	0x61, 0x6e, 0x74, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x32, 0x0a, 0x04, 0x70, 0x69, 0x6e, 0x67, 0x12, 0x13, 0x2e,
	0x6c, 0x75, 0x62, 0x72, 0x69, 0x63, 0x61, 0x6e, 0x74, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x50, 0x6f,
	0x6e, 0x67, 0x1a, 0x13, 0x2e, 0x6c, 0x75, 0x62, 0x72, 0x69, 0x63, 0x61, 0x6e, 0x74, 0x2e, 0x50,
	0x69, 0x6e, 0x67, 0x50, 0x6f, 0x6e, 0x67, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x0d, 0x70, 0x75, 0x73,
	0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x18, 0x2e, 0x6c, 0x75, 0x62,
	0x72, 0x69, 0x63, 0x61, 0x6e, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64,
	0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x18, 0x2e, 0x6c, 0x75, 0x62, 0x72, 0x69, 0x63, 0x61, 0x6e, 0x74,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x00,
	0x42, 0x12, 0x5a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x61, 0x74,
	0x65, 0x77, 0x61, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_gateway_data_proto_rawDescOnce sync.Once
	file_protobuf_gateway_data_proto_rawDescData = file_protobuf_gateway_data_proto_rawDesc
)

func file_protobuf_gateway_data_proto_rawDescGZIP() []byte {
	file_protobuf_gateway_data_proto_rawDescOnce.Do(func() {
		file_protobuf_gateway_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_gateway_data_proto_rawDescData)
	})
	return file_protobuf_gateway_data_proto_rawDescData
}

var file_protobuf_gateway_data_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_protobuf_gateway_data_proto_goTypes = []any{
	(*DataMessage)(nil),   // 0: lubricant.DataMessage
	(*MessageIdInfo)(nil), // 1: lubricant.MessageIdInfo
	(*PingPong)(nil),      // 2: lubricant.PingPong
}
var file_protobuf_gateway_data_proto_depIdxs = []int32{
	0, // 0: lubricant.gatewayService.data:input_type -> lubricant.DataMessage
	2, // 1: lubricant.gatewayService.ping:input_type -> lubricant.PingPong
	1, // 2: lubricant.gatewayService.pushMessageId:input_type -> lubricant.MessageIdInfo
	0, // 3: lubricant.gatewayService.data:output_type -> lubricant.DataMessage
	2, // 4: lubricant.gatewayService.ping:output_type -> lubricant.PingPong
	1, // 5: lubricant.gatewayService.pushMessageId:output_type -> lubricant.MessageIdInfo
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protobuf_gateway_data_proto_init() }
func file_protobuf_gateway_data_proto_init() {
	if File_protobuf_gateway_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_gateway_data_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*DataMessage); i {
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
		file_protobuf_gateway_data_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*MessageIdInfo); i {
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
		file_protobuf_gateway_data_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PingPong); i {
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
			RawDescriptor: file_protobuf_gateway_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_gateway_data_proto_goTypes,
		DependencyIndexes: file_protobuf_gateway_data_proto_depIdxs,
		MessageInfos:      file_protobuf_gateway_data_proto_msgTypes,
	}.Build()
	File_protobuf_gateway_data_proto = out.File
	file_protobuf_gateway_data_proto_rawDesc = nil
	file_protobuf_gateway_data_proto_goTypes = nil
	file_protobuf_gateway_data_proto_depIdxs = nil
}
