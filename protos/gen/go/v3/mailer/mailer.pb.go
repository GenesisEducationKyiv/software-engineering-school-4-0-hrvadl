// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: v3/mailer/mailer.proto

package mailer

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ExchangeFetchedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventID   string  `protobuf:"bytes,1,opt,name=EventID,proto3" json:"EventID,omitempty"`
	EventType string  `protobuf:"bytes,2,opt,name=EventType,proto3" json:"EventType,omitempty"`
	From      string  `protobuf:"bytes,3,opt,name=From,proto3" json:"From,omitempty"`
	To        string  `protobuf:"bytes,4,opt,name=To,proto3" json:"To,omitempty"`
	Rate      float32 `protobuf:"fixed32,5,opt,name=Rate,proto3" json:"Rate,omitempty"`
}

func (x *ExchangeFetchedEvent) Reset() {
	*x = ExchangeFetchedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v3_mailer_mailer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExchangeFetchedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExchangeFetchedEvent) ProtoMessage() {}

func (x *ExchangeFetchedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_v3_mailer_mailer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExchangeFetchedEvent.ProtoReflect.Descriptor instead.
func (*ExchangeFetchedEvent) Descriptor() ([]byte, []int) {
	return file_v3_mailer_mailer_proto_rawDescGZIP(), []int{0}
}

func (x *ExchangeFetchedEvent) GetEventID() string {
	if x != nil {
		return x.EventID
	}
	return ""
}

func (x *ExchangeFetchedEvent) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

func (x *ExchangeFetchedEvent) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *ExchangeFetchedEvent) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *ExchangeFetchedEvent) GetRate() float32 {
	if x != nil {
		return x.Rate
	}
	return 0
}

var File_v3_mailer_mailer_proto protoreflect.FileDescriptor

var file_v3_mailer_mailer_proto_rawDesc = []byte{
	0x0a, 0x16, 0x76, 0x33, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2f, 0x6d, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72,
	0x2e, 0x76, 0x33, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x86, 0x01, 0x0a, 0x14, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x46, 0x65, 0x74,
	0x63, 0x68, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x54, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x54, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x52, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x04, 0x52, 0x61, 0x74, 0x65, 0x42, 0x59, 0x5a, 0x57, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x45,
	0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x79, 0x69, 0x76, 0x2f, 0x73, 0x6f, 0x66,
	0x74, 0x77, 0x61, 0x72, 0x65, 0x2d, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x65, 0x72, 0x69, 0x6e,
	0x67, 0x2d, 0x73, 0x63, 0x68, 0x6f, 0x6f, 0x6c, 0x2d, 0x34, 0x2d, 0x30, 0x2d, 0x68, 0x72, 0x76,
	0x61, 0x64, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x76, 0x33, 0x2f, 0x6d, 0x61,
	0x69, 0x6c, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v3_mailer_mailer_proto_rawDescOnce sync.Once
	file_v3_mailer_mailer_proto_rawDescData = file_v3_mailer_mailer_proto_rawDesc
)

func file_v3_mailer_mailer_proto_rawDescGZIP() []byte {
	file_v3_mailer_mailer_proto_rawDescOnce.Do(func() {
		file_v3_mailer_mailer_proto_rawDescData = protoimpl.X.CompressGZIP(file_v3_mailer_mailer_proto_rawDescData)
	})
	return file_v3_mailer_mailer_proto_rawDescData
}

var file_v3_mailer_mailer_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_v3_mailer_mailer_proto_goTypes = []any{
	(*ExchangeFetchedEvent)(nil), // 0: mailer.v3.ExchangeFetchedEvent
}
var file_v3_mailer_mailer_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v3_mailer_mailer_proto_init() }
func file_v3_mailer_mailer_proto_init() {
	if File_v3_mailer_mailer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v3_mailer_mailer_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ExchangeFetchedEvent); i {
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
			RawDescriptor: file_v3_mailer_mailer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v3_mailer_mailer_proto_goTypes,
		DependencyIndexes: file_v3_mailer_mailer_proto_depIdxs,
		MessageInfos:      file_v3_mailer_mailer_proto_msgTypes,
	}.Build()
	File_v3_mailer_mailer_proto = out.File
	file_v3_mailer_mailer_proto_rawDesc = nil
	file_v3_mailer_mailer_proto_goTypes = nil
	file_v3_mailer_mailer_proto_depIdxs = nil
}
