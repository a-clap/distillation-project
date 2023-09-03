// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: embeddedproto/ds18b20.proto

package embeddedproto

import (
	reflect "reflect"
	sync "sync"

	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DSConfigs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Configs []*DSConfig `protobuf:"bytes,1,rep,name=configs,proto3" json:"configs,omitempty"`
}

func (x *DSConfigs) Reset() {
	*x = DSConfigs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_ds18b20_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DSConfigs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DSConfigs) ProtoMessage() {}

func (x *DSConfigs) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_ds18b20_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DSConfigs.ProtoReflect.Descriptor instead.
func (*DSConfigs) Descriptor() ([]byte, []int) {
	return file_embeddedproto_ds18b20_proto_rawDescGZIP(), []int{0}
}

func (x *DSConfigs) GetConfigs() []*DSConfig {
	if x != nil {
		return x.Configs
	}
	return nil
}

type DSConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID           string  `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name         string  `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Correction   float32 `protobuf:"fixed32,3,opt,name=Correction,proto3" json:"Correction,omitempty"`
	Resolution   int32   `protobuf:"varint,4,opt,name=Resolution,proto3" json:"Resolution,omitempty"`
	PollInterval int32   `protobuf:"varint,5,opt,name=PollInterval,proto3" json:"PollInterval,omitempty"`
	Samples      uint32  `protobuf:"varint,6,opt,name=Samples,proto3" json:"Samples,omitempty"`
	Enabled      bool    `protobuf:"varint,7,opt,name=Enabled,proto3" json:"Enabled,omitempty"`
}

func (x *DSConfig) Reset() {
	*x = DSConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_ds18b20_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DSConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DSConfig) ProtoMessage() {}

func (x *DSConfig) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_ds18b20_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DSConfig.ProtoReflect.Descriptor instead.
func (*DSConfig) Descriptor() ([]byte, []int) {
	return file_embeddedproto_ds18b20_proto_rawDescGZIP(), []int{1}
}

func (x *DSConfig) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *DSConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DSConfig) GetCorrection() float32 {
	if x != nil {
		return x.Correction
	}
	return 0
}

func (x *DSConfig) GetResolution() int32 {
	if x != nil {
		return x.Resolution
	}
	return 0
}

func (x *DSConfig) GetPollInterval() int32 {
	if x != nil {
		return x.PollInterval
	}
	return 0
}

func (x *DSConfig) GetSamples() uint32 {
	if x != nil {
		return x.Samples
	}
	return 0
}

func (x *DSConfig) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

type DSTemperatures struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Temps []*DSTemperature `protobuf:"bytes,1,rep,name=temps,proto3" json:"temps,omitempty"`
}

func (x *DSTemperatures) Reset() {
	*x = DSTemperatures{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_ds18b20_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DSTemperatures) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DSTemperatures) ProtoMessage() {}

func (x *DSTemperatures) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_ds18b20_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DSTemperatures.ProtoReflect.Descriptor instead.
func (*DSTemperatures) Descriptor() ([]byte, []int) {
	return file_embeddedproto_ds18b20_proto_rawDescGZIP(), []int{2}
}

func (x *DSTemperatures) GetTemps() []*DSTemperature {
	if x != nil {
		return x.Temps
	}
	return nil
}

type DSTemperature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Readings []*DSReadings `protobuf:"bytes,1,rep,name=readings,proto3" json:"readings,omitempty"`
}

func (x *DSTemperature) Reset() {
	*x = DSTemperature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_ds18b20_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DSTemperature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DSTemperature) ProtoMessage() {}

func (x *DSTemperature) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_ds18b20_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DSTemperature.ProtoReflect.Descriptor instead.
func (*DSTemperature) Descriptor() ([]byte, []int) {
	return file_embeddedproto_ds18b20_proto_rawDescGZIP(), []int{3}
}

func (x *DSTemperature) GetReadings() []*DSReadings {
	if x != nil {
		return x.Readings
	}
	return nil
}

type DSReadings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID          string  `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Temperature float32 `protobuf:"fixed32,2,opt,name=Temperature,proto3" json:"Temperature,omitempty"`
	Average     float32 `protobuf:"fixed32,3,opt,name=Average,proto3" json:"Average,omitempty"`
	StampMillis int64   `protobuf:"varint,4,opt,name=StampMillis,proto3" json:"StampMillis,omitempty"`
	Error       string  `protobuf:"bytes,5,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *DSReadings) Reset() {
	*x = DSReadings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_embeddedproto_ds18b20_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DSReadings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DSReadings) ProtoMessage() {}

func (x *DSReadings) ProtoReflect() protoreflect.Message {
	mi := &file_embeddedproto_ds18b20_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DSReadings.ProtoReflect.Descriptor instead.
func (*DSReadings) Descriptor() ([]byte, []int) {
	return file_embeddedproto_ds18b20_proto_rawDescGZIP(), []int{4}
}

func (x *DSReadings) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *DSReadings) GetTemperature() float32 {
	if x != nil {
		return x.Temperature
	}
	return 0
}

func (x *DSReadings) GetAverage() float32 {
	if x != nil {
		return x.Average
	}
	return 0
}

func (x *DSReadings) GetStampMillis() int64 {
	if x != nil {
		return x.StampMillis
	}
	return 0
}

func (x *DSReadings) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_embeddedproto_ds18b20_proto protoreflect.FileDescriptor

var file_embeddedproto_ds18b20_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x64, 0x73, 0x31, 0x38, 0x62, 0x32, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x65,
	0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x09, 0x44, 0x53, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12, 0x31, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64,
	0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x53, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x22, 0xc6, 0x01, 0x0a, 0x08, 0x44, 0x53,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x6f,
	0x72, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0a,
	0x43, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x65,
	0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x52, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x50, 0x6f,
	0x6c, 0x6c, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x50, 0x6f, 0x6c, 0x6c, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x18,
	0x0a, 0x07, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x6e, 0x61, 0x62,
	0x6c, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x45, 0x6e, 0x61, 0x62, 0x6c,
	0x65, 0x64, 0x22, 0x44, 0x0a, 0x0e, 0x44, 0x53, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x05, 0x74, 0x65, 0x6d, 0x70, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x53, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x52, 0x05, 0x74, 0x65, 0x6d, 0x70, 0x73, 0x22, 0x46, 0x0a, 0x0d, 0x44, 0x53, 0x54, 0x65,
	0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x72, 0x65, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6d,
	0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x53, 0x52, 0x65,
	0x61, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x73,
	0x22, 0x90, 0x01, 0x0a, 0x0a, 0x44, 0x53, 0x52, 0x65, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12,
	0x20, 0x0a, 0x0b, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x07, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x53,
	0x74, 0x61, 0x6d, 0x70, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0b, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x32, 0xd2, 0x01, 0x0a, 0x02, 0x44, 0x53, 0x12, 0x3b, 0x0a, 0x05, 0x44, 0x53,
	0x47, 0x65, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x18, 0x2e, 0x65, 0x6d,
	0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x53, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x73, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0b, 0x44, 0x53, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x12, 0x17, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65,
	0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x53, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a,
	0x17, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x44, 0x53, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x11, 0x44, 0x53,
	0x47, 0x65, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1d, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64,
	0x65, 0x64, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x53, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x00, 0x42, 0x27, 0x50, 0x01, 0x5a, 0x23, 0x65, 0x6d,
	0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x6d, 0x62, 0x65, 0x64,
	0x64, 0x65, 0x64, 0x2f, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_embeddedproto_ds18b20_proto_rawDescOnce sync.Once
	file_embeddedproto_ds18b20_proto_rawDescData = file_embeddedproto_ds18b20_proto_rawDesc
)

func file_embeddedproto_ds18b20_proto_rawDescGZIP() []byte {
	file_embeddedproto_ds18b20_proto_rawDescOnce.Do(func() {
		file_embeddedproto_ds18b20_proto_rawDescData = protoimpl.X.CompressGZIP(file_embeddedproto_ds18b20_proto_rawDescData)
	})
	return file_embeddedproto_ds18b20_proto_rawDescData
}

var file_embeddedproto_ds18b20_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_embeddedproto_ds18b20_proto_goTypes = []interface{}{
	(*DSConfigs)(nil),      // 0: embeddedproto.DSConfigs
	(*DSConfig)(nil),       // 1: embeddedproto.DSConfig
	(*DSTemperatures)(nil), // 2: embeddedproto.DSTemperatures
	(*DSTemperature)(nil),  // 3: embeddedproto.DSTemperature
	(*DSReadings)(nil),     // 4: embeddedproto.DSReadings
	(*empty.Empty)(nil),    // 5: google.protobuf.Empty
}
var file_embeddedproto_ds18b20_proto_depIdxs = []int32{
	1, // 0: embeddedproto.DSConfigs.configs:type_name -> embeddedproto.DSConfig
	3, // 1: embeddedproto.DSTemperatures.temps:type_name -> embeddedproto.DSTemperature
	4, // 2: embeddedproto.DSTemperature.readings:type_name -> embeddedproto.DSReadings
	5, // 3: embeddedproto.DS.DSGet:input_type -> google.protobuf.Empty
	1, // 4: embeddedproto.DS.DSConfigure:input_type -> embeddedproto.DSConfig
	5, // 5: embeddedproto.DS.DSGetTemperatures:input_type -> google.protobuf.Empty
	0, // 6: embeddedproto.DS.DSGet:output_type -> embeddedproto.DSConfigs
	1, // 7: embeddedproto.DS.DSConfigure:output_type -> embeddedproto.DSConfig
	2, // 8: embeddedproto.DS.DSGetTemperatures:output_type -> embeddedproto.DSTemperatures
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_embeddedproto_ds18b20_proto_init() }
func file_embeddedproto_ds18b20_proto_init() {
	if File_embeddedproto_ds18b20_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_embeddedproto_ds18b20_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DSConfigs); i {
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
		file_embeddedproto_ds18b20_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DSConfig); i {
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
		file_embeddedproto_ds18b20_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DSTemperatures); i {
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
		file_embeddedproto_ds18b20_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DSTemperature); i {
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
		file_embeddedproto_ds18b20_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DSReadings); i {
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
			RawDescriptor: file_embeddedproto_ds18b20_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_embeddedproto_ds18b20_proto_goTypes,
		DependencyIndexes: file_embeddedproto_ds18b20_proto_depIdxs,
		MessageInfos:      file_embeddedproto_ds18b20_proto_msgTypes,
	}.Build()
	File_embeddedproto_ds18b20_proto = out.File
	file_embeddedproto_ds18b20_proto_rawDesc = nil
	file_embeddedproto_ds18b20_proto_goTypes = nil
	file_embeddedproto_ds18b20_proto_depIdxs = nil
}
