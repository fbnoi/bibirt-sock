// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: websocket/pb/ctrl.proto

package pb

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

type Heartbeat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UpTimestamp   int64 `protobuf:"varint,1,opt,name=upTimestamp,proto3" json:"upTimestamp,omitempty"`
	DownTimestamp int64 `protobuf:"varint,2,opt,name=downTimestamp,proto3" json:"downTimestamp,omitempty"`
}

func (x *Heartbeat) Reset() {
	*x = Heartbeat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_websocket_pb_ctrl_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Heartbeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Heartbeat) ProtoMessage() {}

func (x *Heartbeat) ProtoReflect() protoreflect.Message {
	mi := &file_websocket_pb_ctrl_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Heartbeat.ProtoReflect.Descriptor instead.
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return file_websocket_pb_ctrl_proto_rawDescGZIP(), []int{0}
}

func (x *Heartbeat) GetUpTimestamp() int64 {
	if x != nil {
		return x.UpTimestamp
	}
	return 0
}

func (x *Heartbeat) GetDownTimestamp() int64 {
	if x != nil {
		return x.DownTimestamp
	}
	return 0
}

type Close struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Close) Reset() {
	*x = Close{}
	if protoimpl.UnsafeEnabled {
		mi := &file_websocket_pb_ctrl_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Close) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Close) ProtoMessage() {}

func (x *Close) ProtoReflect() protoreflect.Message {
	mi := &file_websocket_pb_ctrl_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Close.ProtoReflect.Descriptor instead.
func (*Close) Descriptor() ([]byte, []int) {
	return file_websocket_pb_ctrl_proto_rawDescGZIP(), []int{1}
}

type Maintain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Maintain) Reset() {
	*x = Maintain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_websocket_pb_ctrl_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Maintain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Maintain) ProtoMessage() {}

func (x *Maintain) ProtoReflect() protoreflect.Message {
	mi := &file_websocket_pb_ctrl_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Maintain.ProtoReflect.Descriptor instead.
func (*Maintain) Descriptor() ([]byte, []int) {
	return file_websocket_pb_ctrl_proto_rawDescGZIP(), []int{2}
}

var File_websocket_pb_ctrl_proto protoreflect.FileDescriptor

var file_websocket_pb_ctrl_proto_rawDesc = []byte{
	0x0a, 0x17, 0x77, 0x65, 0x62, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2f, 0x70, 0x62, 0x2f, 0x63,
	0x74, 0x72, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x53, 0x0a,
	0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x75, 0x70,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0b, 0x75, 0x70, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x24, 0x0a, 0x0d,
	0x64, 0x6f, 0x77, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0d, 0x64, 0x6f, 0x77, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x22, 0x07, 0x0a, 0x05, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x22, 0x0a, 0x0a, 0x08, 0x4d,
	0x61, 0x69, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x3b, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_websocket_pb_ctrl_proto_rawDescOnce sync.Once
	file_websocket_pb_ctrl_proto_rawDescData = file_websocket_pb_ctrl_proto_rawDesc
)

func file_websocket_pb_ctrl_proto_rawDescGZIP() []byte {
	file_websocket_pb_ctrl_proto_rawDescOnce.Do(func() {
		file_websocket_pb_ctrl_proto_rawDescData = protoimpl.X.CompressGZIP(file_websocket_pb_ctrl_proto_rawDescData)
	})
	return file_websocket_pb_ctrl_proto_rawDescData
}

var file_websocket_pb_ctrl_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_websocket_pb_ctrl_proto_goTypes = []interface{}{
	(*Heartbeat)(nil), // 0: pb.Heartbeat
	(*Close)(nil),     // 1: pb.Close
	(*Maintain)(nil),  // 2: pb.Maintain
}
var file_websocket_pb_ctrl_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_websocket_pb_ctrl_proto_init() }
func file_websocket_pb_ctrl_proto_init() {
	if File_websocket_pb_ctrl_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_websocket_pb_ctrl_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Heartbeat); i {
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
		file_websocket_pb_ctrl_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Close); i {
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
		file_websocket_pb_ctrl_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Maintain); i {
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
			RawDescriptor: file_websocket_pb_ctrl_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_websocket_pb_ctrl_proto_goTypes,
		DependencyIndexes: file_websocket_pb_ctrl_proto_depIdxs,
		MessageInfos:      file_websocket_pb_ctrl_proto_msgTypes,
	}.Build()
	File_websocket_pb_ctrl_proto = out.File
	file_websocket_pb_ctrl_proto_rawDesc = nil
	file_websocket_pb_ctrl_proto_goTypes = nil
	file_websocket_pb_ctrl_proto_depIdxs = nil
}
