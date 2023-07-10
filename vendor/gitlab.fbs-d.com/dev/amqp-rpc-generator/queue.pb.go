// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: queue.proto

package amqp_rpc_generator

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Type int32

const (
	Type_RPC             Type = 0
	Type_Exchange        Type = 1
	Type_QueueType       Type = 2
	Type_TrialDelay      Type = 3
	Type_BatchTrialDelay Type = 4
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "RPC",
		1: "Exchange",
		2: "QueueType",
		3: "TrialDelay",
		4: "BatchTrialDelay",
	}
	Type_value = map[string]int32{
		"RPC":             0,
		"Exchange":        1,
		"QueueType":       2,
		"TrialDelay":      3,
		"BatchTrialDelay": 4,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_queue_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_queue_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_queue_proto_rawDescGZIP(), []int{0}
}

type Queue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type            Type   `protobuf:"varint,1,opt,name=Type,proto3,enum=queue.Type" json:"Type"`
	Name            string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name"`
	ExchangeName    string `protobuf:"bytes,3,opt,name=ExchangeName,proto3" json:"ExchangeName"`
	ExchangeType    string `protobuf:"bytes,4,opt,name=ExchangeType,proto3" json:"ExchangeType"`
	RoutingKey      string `protobuf:"bytes,5,opt,name=RoutingKey,proto3" json:"RoutingKey"`
	OnlyClient      bool   `protobuf:"varint,6,opt,name=OnlyClient,proto3" json:"OnlyClient"`
	Trials          int64  `protobuf:"varint,7,opt,name=Trials,proto3" json:"Trials"`
	Delay           int64  `protobuf:"varint,8,opt,name=Delay,proto3" json:"Delay"`
	EventSkipErrors bool   `protobuf:"varint,9,opt,name=EventSkipErrors,proto3" json:"EventSkipErrors"`
}

func (x *Queue) Reset() {
	*x = Queue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_queue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Queue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Queue) ProtoMessage() {}

func (x *Queue) ProtoReflect() protoreflect.Message {
	mi := &file_queue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Queue.ProtoReflect.Descriptor instead.
func (*Queue) Descriptor() ([]byte, []int) {
	return file_queue_proto_rawDescGZIP(), []int{0}
}

func (x *Queue) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_RPC
}

func (x *Queue) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Queue) GetExchangeName() string {
	if x != nil {
		return x.ExchangeName
	}
	return ""
}

func (x *Queue) GetExchangeType() string {
	if x != nil {
		return x.ExchangeType
	}
	return ""
}

func (x *Queue) GetRoutingKey() string {
	if x != nil {
		return x.RoutingKey
	}
	return ""
}

func (x *Queue) GetOnlyClient() bool {
	if x != nil {
		return x.OnlyClient
	}
	return false
}

func (x *Queue) GetTrials() int64 {
	if x != nil {
		return x.Trials
	}
	return 0
}

func (x *Queue) GetDelay() int64 {
	if x != nil {
		return x.Delay
	}
	return 0
}

func (x *Queue) GetEventSkipErrors() bool {
	if x != nil {
		return x.EventSkipErrors
	}
	return false
}

var file_queue_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*Queue)(nil),
		Field:         1001,
		Name:          "queue.queue",
		Tag:           "bytes,1001,opt,name=queue",
		Filename:      "queue.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional queue.Queue queue = 1001;
	E_Queue = &file_queue_proto_extTypes[0]
)

var File_queue_proto protoreflect.FileDescriptor

var file_queue_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x02, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x12, 0x1f, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b,
	0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x45, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x45, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x12, 0x1e, 0x0a,
	0x0a, 0x4f, 0x6e, 0x6c, 0x79, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0a, 0x4f, 0x6e, 0x6c, 0x79, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x54, 0x72, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x54,
	0x72, 0x69, 0x61, 0x6c, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x12, 0x28, 0x0a, 0x0f, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x53, 0x6b, 0x69, 0x70, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x6b, 0x69, 0x70, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x73, 0x2a, 0x51, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x07, 0x0a,
	0x03, 0x52, 0x50, 0x43, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x51, 0x75, 0x65, 0x75, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x72, 0x69, 0x61, 0x6c, 0x44, 0x65, 0x6c, 0x61,
	0x79, 0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x42, 0x61, 0x74, 0x63, 0x68, 0x54, 0x72, 0x69, 0x61,
	0x6c, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x10, 0x04, 0x3a, 0x43, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xe9, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x42, 0x3c, 0x5a,
	0x3a, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x66, 0x62, 0x73, 0x2d, 0x64, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x64, 0x65, 0x76, 0x2f, 0x61, 0x6d, 0x71, 0x70, 0x2d, 0x72, 0x70, 0x63, 0x2d, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x3b, 0x61, 0x6d, 0x71, 0x70, 0x5f, 0x72, 0x70,
	0x63, 0x5f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_queue_proto_rawDescOnce sync.Once
	file_queue_proto_rawDescData = file_queue_proto_rawDesc
)

func file_queue_proto_rawDescGZIP() []byte {
	file_queue_proto_rawDescOnce.Do(func() {
		file_queue_proto_rawDescData = protoimpl.X.CompressGZIP(file_queue_proto_rawDescData)
	})
	return file_queue_proto_rawDescData
}

var file_queue_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_queue_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_queue_proto_goTypes = []interface{}{
	(Type)(0),                          // 0: queue.Type
	(*Queue)(nil),                      // 1: queue.Queue
	(*descriptorpb.MethodOptions)(nil), // 2: google.protobuf.MethodOptions
}
var file_queue_proto_depIdxs = []int32{
	0, // 0: queue.Queue.Type:type_name -> queue.Type
	2, // 1: queue.queue:extendee -> google.protobuf.MethodOptions
	1, // 2: queue.queue:type_name -> queue.Queue
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	2, // [2:3] is the sub-list for extension type_name
	1, // [1:2] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_queue_proto_init() }
func file_queue_proto_init() {
	if File_queue_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_queue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Queue); i {
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
			RawDescriptor: file_queue_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_queue_proto_goTypes,
		DependencyIndexes: file_queue_proto_depIdxs,
		EnumInfos:         file_queue_proto_enumTypes,
		MessageInfos:      file_queue_proto_msgTypes,
		ExtensionInfos:    file_queue_proto_extTypes,
	}.Build()
	File_queue_proto = out.File
	file_queue_proto_rawDesc = nil
	file_queue_proto_goTypes = nil
	file_queue_proto_depIdxs = nil
}
