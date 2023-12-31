// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: embeddedproto/pt100.proto

package embeddedproto

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PTConfigs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Configs []*PTConfig `protobuf:"bytes,1,rep,name=configs,proto3" json:"configs,omitempty"`
}

func (x *PTConfigs) Reset() {
	*x = PTConfigs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_pt100_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PTConfigs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PTConfigs) ProtoMessage() {}

func (x *PTConfigs) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_pt100_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PTConfigs.ProtoReflect.Descriptor instead.
func (*PTConfigs) Descriptor() ([]byte, []int) {
	return file_embeddedproto_pt100_proto_rawDescGZIP(), []int{0}
}

func (x *PTConfigs) GetConfigs() []*PTConfig {
	if x != nil {
		return x.Configs
	}
	return nil
}

type PTConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID           string  `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name         string  `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Correction   float32 `protobuf:"fixed32,3,opt,name=Correction,proto3" json:"Correction,omitempty"`
	PollInterval int32   `protobuf:"varint,5,opt,name=PollInterval,proto3" json:"PollInterval,omitempty"`
	Samples      uint32  `protobuf:"varint,6,opt,name=Samples,proto3" json:"Samples,omitempty"`
	Enabled      bool    `protobuf:"varint,7,opt,name=Enabled,proto3" json:"Enabled,omitempty"`
	Async        bool    `protobuf:"varint,8,opt,name=Async,proto3" json:"Async,omitempty"`
}

func (x *PTConfig) Reset() {
	*x = PTConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_pt100_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PTConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PTConfig) ProtoMessage() {}

func (x *PTConfig) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_pt100_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PTConfig.ProtoReflect.Descriptor instead.
func (*PTConfig) Descriptor() ([]byte, []int) {
	return file_embeddedproto_pt100_proto_rawDescGZIP(), []int{1}
}

func (x *PTConfig) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *PTConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PTConfig) GetCorrection() float32 {
	if x != nil {
		return x.Correction
	}
	return 0
}

func (x *PTConfig) GetPollInterval() int32 {
	if x != nil {
		return x.PollInterval
	}
	return 0
}

func (x *PTConfig) GetSamples() uint32 {
	if x != nil {
		return x.Samples
	}
	return 0
}

func (x *PTConfig) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *PTConfig) GetAsync() bool {
	if x != nil {
		return x.Async
	}
	return false
}

type PTTemperatures struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Temps []*PTTemperature `protobuf:"bytes,1,rep,name=temps,proto3" json:"temps,omitempty"`
}

func (x *PTTemperatures) Reset() {
	*x = PTTemperatures{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_pt100_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PTTemperatures) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PTTemperatures) ProtoMessage() {}

func (x *PTTemperatures) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_pt100_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PTTemperatures.ProtoReflect.Descriptor instead.
func (*PTTemperatures) Descriptor() ([]byte, []int) {
	return file_embeddedproto_pt100_proto_rawDescGZIP(), []int{2}
}

func (x *PTTemperatures) GetTemps() []*PTTemperature {
	if x != nil {
		return x.Temps
	}
	return nil
}

type PTTemperature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Readings []*PTReadings `protobuf:"bytes,1,rep,name=readings,proto3" json:"readings,omitempty"`
}

func (x *PTTemperature) Reset() {
	*x = PTTemperature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_pt100_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PTTemperature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PTTemperature) ProtoMessage() {}

func (x *PTTemperature) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_pt100_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PTTemperature.ProtoReflect.Descriptor instead.
func (*PTTemperature) Descriptor() ([]byte, []int) {
	return file_embeddedproto_pt100_proto_rawDescGZIP(), []int{3}
}

func (x *PTTemperature) GetReadings() []*PTReadings {
	if x != nil {
		return x.Readings
	}
	return nil
}

type PTReadings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID          string  `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Temperature float32 `protobuf:"fixed32,2,opt,name=Temperature,proto3" json:"Temperature,omitempty"`
	Average     float32 `protobuf:"fixed32,3,opt,name=Average,proto3" json:"Average,omitempty"`
	StampMillis int64   `protobuf:"varint,4,opt,name=StampMillis,proto3" json:"StampMillis,omitempty"`
	Error       string  `protobuf:"bytes,5,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *PTReadings) Reset() {
	*x = PTReadings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_pt100_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PTReadings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PTReadings) ProtoMessage() {}

func (x *PTReadings) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_pt100_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PTReadings.ProtoReflect.Descriptor instead.
func (*PTReadings) Descriptor() ([]byte, []int) {
	return file_embeddedproto_pt100_proto_rawDescGZIP(), []int{4}
}

func (x *PTReadings) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *PTReadings) GetTemperature() float32 {
	if x != nil {
		return x.Temperature
	}
	return 0
}

func (x *PTReadings) GetAverage() float32 {
	if x != nil {
		return x.Average
	}
	return 0
}

func (x *PTReadings) GetStampMillis() int64 {
	if x != nil {
		return x.StampMillis
	}
	return 0
}

func (x *PTReadings) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_embeddedproto_pt100_proto protoreflect.FileDescriptor

var file_embeddedproto_pt100_proto_rawDesc = []byte{
	0x0a, 0x19, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x70, 0x74, 0x31, 0x30, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x65, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x09, 0x50, 0x54, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x73, 0x12, 0x31, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x54, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x07,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x22, 0xbc, 0x01, 0x0a, 0x08, 0x50, 0x54, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x6f, 0x72, 0x72,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0a, 0x43, 0x6f,
	0x72, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x50, 0x6f, 0x6c, 0x6c,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c,
	0x50, 0x6f, 0x6c, 0x6c, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07,
	0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x53,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65,
	0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x41, 0x73, 0x79, 0x6e, 0x63, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x05, 0x41, 0x73, 0x79, 0x6e, 0x63, 0x22, 0x44, 0x0a, 0x0e, 0x50, 0x54, 0x54, 0x65, 0x6d, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x05, 0x74, 0x65, 0x6d, 0x70,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64,
	0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x54, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x05, 0x74, 0x65, 0x6d, 0x70, 0x73, 0x22, 0x46, 0x0a, 0x0d,
	0x50, 0x54, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x35, 0x0a,
	0x08, 0x72, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x54, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x73, 0x22, 0x90, 0x01, 0x0a, 0x0a, 0x50, 0x54, 0x52, 0x65, 0x61, 0x64, 0x69,
	0x6e, 0x67, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x07, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x4d, 0x69, 0x6c, 0x6c, 0x69,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x32, 0xd2, 0x01, 0x0a, 0x02, 0x50, 0x54, 0x12, 0x3b,
	0x0a, 0x05, 0x50, 0x54, 0x47, 0x65, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x18, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x54, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0b, 0x50,
	0x54, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x12, 0x17, 0x2e, 0x65, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x54, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x1a, 0x17, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x54, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x00, 0x12, 0x4c,
	0x0a, 0x11, 0x50, 0x54, 0x47, 0x65, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1d, 0x2e, 0x65, 0x6d,
	0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x54, 0x54, 0x65,
	0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x00, 0x42, 0x27, 0x50, 0x01,
	0x5a, 0x23, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65,
	0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x2f, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_embeddedproto_pt100_proto_rawDescOnce sync.Once
	file_embeddedproto_pt100_proto_rawDescData = file_embeddedproto_pt100_proto_rawDesc
)

func file_embeddedproto_pt100_proto_rawDescGZIP() []byte {
	file_embeddedproto_pt100_proto_rawDescOnce.Do(func() {
		file_embeddedproto_pt100_proto_rawDescData = protoimpl.X.CompressGZIP(file_embeddedproto_pt100_proto_rawDescData)
	})
	return file_embeddedproto_pt100_proto_rawDescData
}

var file_embeddedproto_pt100_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_embeddedproto_pt100_proto_goTypes = []interface{}{
	(*PTConfigs)(nil),      // 0: embeddedproto.PTConfigs
	(*PTConfig)(nil),       // 1: embeddedproto.PTConfig
	(*PTTemperatures)(nil), // 2: embeddedproto.PTTemperatures
	(*PTTemperature)(nil),  // 3: embeddedproto.PTTemperature
	(*PTReadings)(nil),     // 4: embeddedproto.PTReadings
	(*empty.Empty)(nil),    // 5: google.protobuf.Empty
}
var file_embeddedproto_pt100_proto_depIdxs = []int32{
	1, // 0: embeddedproto.PTConfigs.configs:type_name -> embeddedproto.PTConfig
	3, // 1: embeddedproto.PTTemperatures.temps:type_name -> embeddedproto.PTTemperature
	4, // 2: embeddedproto.PTTemperature.readings:type_name -> embeddedproto.PTReadings
	5, // 3: embeddedproto.PT.PTGet:input_type -> google.protobuf.Empty
	1, // 4: embeddedproto.PT.PTConfigure:input_type -> embeddedproto.PTConfig
	5, // 5: embeddedproto.PT.PTGetTemperatures:input_type -> google.protobuf.Empty
	0, // 6: embeddedproto.PT.PTGet:output_type -> embeddedproto.PTConfigs
	1, // 7: embeddedproto.PT.PTConfigure:output_type -> embeddedproto.PTConfig
	2, // 8: embeddedproto.PT.PTGetTemperatures:output_type -> embeddedproto.PTTemperatures
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_embeddedproto_pt100_proto_init() }
func file_embeddedproto_pt100_proto_init() {
	if File_embeddedproto_pt100_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_embeddedproto_pt100_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PTConfigs); i {
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
		file_embeddedproto_pt100_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PTConfig); i {
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
		file_embeddedproto_pt100_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PTTemperatures); i {
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
		file_embeddedproto_pt100_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PTTemperature); i {
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
		file_embeddedproto_pt100_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PTReadings); i {
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
			RawDescriptor: file_embeddedproto_pt100_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_embeddedproto_pt100_proto_goTypes,
		DependencyIndexes: file_embeddedproto_pt100_proto_depIdxs,
		MessageInfos:      file_embeddedproto_pt100_proto_msgTypes,
	}.Build()
	File_embeddedproto_pt100_proto = out.File
	file_embeddedproto_pt100_proto_rawDesc = nil
	file_embeddedproto_pt100_proto_goTypes = nil
	file_embeddedproto_pt100_proto_depIdxs = nil
}
