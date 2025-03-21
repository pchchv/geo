// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.12.4
// source: vector_tile.proto

package vectortile

import (
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

type Tile_GeomType int32

const (
	Tile_UNKNOWN    Tile_GeomType = 0
	Tile_POINT      Tile_GeomType = 1
	Tile_LINESTRING Tile_GeomType = 2
	Tile_POLYGON    Tile_GeomType = 3
)

// Enum value maps for Tile_GeomType.
var (
	Tile_GeomType_name = map[int32]string{
		0: "UNKNOWN",
		1: "POINT",
		2: "LINESTRING",
		3: "POLYGON",
	}
	Tile_GeomType_value = map[string]int32{
		"UNKNOWN":    0,
		"POINT":      1,
		"LINESTRING": 2,
		"POLYGON":    3,
	}
)

func (x Tile_GeomType) Enum() *Tile_GeomType {
	p := new(Tile_GeomType)
	*p = x
	return p
}

func (x Tile_GeomType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Tile_GeomType) Descriptor() protoreflect.EnumDescriptor {
	return file_vector_tile_proto_enumTypes[0].Descriptor()
}

func (Tile_GeomType) Type() protoreflect.EnumType {
	return &file_vector_tile_proto_enumTypes[0]
}

func (x Tile_GeomType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Tile_GeomType.Descriptor instead.
func (Tile_GeomType) EnumDescriptor() ([]byte, []int) {
	return file_vector_tile_proto_rawDescGZIP(), []int{0, 0}
}

type Tile struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Layers        []*Tile_Layer          `protobuf:"bytes,3,rep,name=layers,proto3" json:"layers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Tile) Reset() {
	*x = Tile{}
	mi := &file_vector_tile_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tile) ProtoMessage() {}

func (x *Tile) ProtoReflect() protoreflect.Message {
	mi := &file_vector_tile_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tile.ProtoReflect.Descriptor instead.
func (*Tile) Descriptor() ([]byte, []int) {
	return file_vector_tile_proto_rawDescGZIP(), []int{0}
}

func (x *Tile) GetLayers() []*Tile_Layer {
	if x != nil {
		return x.Layers
	}
	return nil
}

type Tile_Value struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Value:
	//
	//	*Tile_Value_StringValue
	//	*Tile_Value_FloatValue
	//	*Tile_Value_DoubleValue
	//	*Tile_Value_IntValue
	//	*Tile_Value_UintValue
	//	*Tile_Value_SintValue
	//	*Tile_Value_BoolValue
	Value         isTile_Value_Value `protobuf_oneof:"value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Tile_Value) Reset() {
	*x = Tile_Value{}
	mi := &file_vector_tile_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tile_Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tile_Value) ProtoMessage() {}

func (x *Tile_Value) ProtoReflect() protoreflect.Message {
	mi := &file_vector_tile_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tile_Value.ProtoReflect.Descriptor instead.
func (*Tile_Value) Descriptor() ([]byte, []int) {
	return file_vector_tile_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Tile_Value) GetValue() isTile_Value_Value {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Tile_Value) GetStringValue() string {
	if x != nil {
		if x, ok := x.Value.(*Tile_Value_StringValue); ok {
			return x.StringValue
		}
	}
	return ""
}

func (x *Tile_Value) GetFloatValue() float32 {
	if x != nil {
		if x, ok := x.Value.(*Tile_Value_FloatValue); ok {
			return x.FloatValue
		}
	}
	return 0
}

func (x *Tile_Value) GetDoubleValue() float64 {
	if x != nil {
		if x, ok := x.Value.(*Tile_Value_DoubleValue); ok {
			return x.DoubleValue
		}
	}
	return 0
}

func (x *Tile_Value) GetIntValue() int64 {
	if x != nil {
		if x, ok := x.Value.(*Tile_Value_IntValue); ok {
			return x.IntValue
		}
	}
	return 0
}

func (x *Tile_Value) GetUintValue() uint64 {
	if x != nil {
		if x, ok := x.Value.(*Tile_Value_UintValue); ok {
			return x.UintValue
		}
	}
	return 0
}

func (x *Tile_Value) GetSintValue() int64 {
	if x != nil {
		if x, ok := x.Value.(*Tile_Value_SintValue); ok {
			return x.SintValue
		}
	}
	return 0
}

func (x *Tile_Value) GetBoolValue() bool {
	if x != nil {
		if x, ok := x.Value.(*Tile_Value_BoolValue); ok {
			return x.BoolValue
		}
	}
	return false
}

type isTile_Value_Value interface {
	isTile_Value_Value()
}

type Tile_Value_StringValue struct {
	StringValue string `protobuf:"bytes,1,opt,name=string_value,json=stringValue,proto3,oneof"`
}

type Tile_Value_FloatValue struct {
	FloatValue float32 `protobuf:"fixed32,2,opt,name=float_value,json=floatValue,proto3,oneof"`
}

type Tile_Value_DoubleValue struct {
	DoubleValue float64 `protobuf:"fixed64,3,opt,name=double_value,json=doubleValue,proto3,oneof"`
}

type Tile_Value_IntValue struct {
	IntValue int64 `protobuf:"varint,4,opt,name=int_value,json=intValue,proto3,oneof"`
}

type Tile_Value_UintValue struct {
	UintValue uint64 `protobuf:"varint,5,opt,name=uint_value,json=uintValue,proto3,oneof"`
}

type Tile_Value_SintValue struct {
	SintValue int64 `protobuf:"zigzag64,6,opt,name=sint_value,json=sintValue,proto3,oneof"`
}

type Tile_Value_BoolValue struct {
	BoolValue bool `protobuf:"varint,7,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

func (*Tile_Value_StringValue) isTile_Value_Value() {}

func (*Tile_Value_FloatValue) isTile_Value_Value() {}

func (*Tile_Value_DoubleValue) isTile_Value_Value() {}

func (*Tile_Value_IntValue) isTile_Value_Value() {}

func (*Tile_Value_UintValue) isTile_Value_Value() {}

func (*Tile_Value_SintValue) isTile_Value_Value() {}

func (*Tile_Value_BoolValue) isTile_Value_Value() {}

type Tile_Feature struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Tags          []uint32               `protobuf:"varint,2,rep,packed,name=tags,proto3" json:"tags,omitempty"`
	Type          Tile_GeomType          `protobuf:"varint,3,opt,name=type,proto3,enum=vector_tile.Tile_GeomType" json:"type,omitempty"`
	Geometry      []uint32               `protobuf:"varint,4,rep,packed,name=geometry,proto3" json:"geometry,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

const Default_Tile_Layer_Version uint32 = 1
const Default_Tile_Layer_Extent uint32 = 4096

func (x *Tile_Feature) Reset() {
	*x = Tile_Feature{}
	mi := &file_vector_tile_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tile_Feature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tile_Feature) ProtoMessage() {}

func (x *Tile_Feature) ProtoReflect() protoreflect.Message {
	mi := &file_vector_tile_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tile_Feature.ProtoReflect.Descriptor instead.
func (*Tile_Feature) Descriptor() ([]byte, []int) {
	return file_vector_tile_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Tile_Feature) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Tile_Feature) GetTags() []uint32 {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Tile_Feature) GetType() Tile_GeomType {
	if x != nil {
		return x.Type
	}
	return Tile_UNKNOWN
}

func (x *Tile_Feature) GetGeometry() []uint32 {
	if x != nil {
		return x.Geometry
	}
	return nil
}

type Tile_Layer struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Version       uint32                 `protobuf:"varint,15,opt,name=version,proto3" json:"version,omitempty"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Features      []*Tile_Feature        `protobuf:"bytes,2,rep,name=features,proto3" json:"features,omitempty"`
	Keys          []string               `protobuf:"bytes,3,rep,name=keys,proto3" json:"keys,omitempty"`
	Values        []*Tile_Value          `protobuf:"bytes,4,rep,name=values,proto3" json:"values,omitempty"`
	Extent        uint32                 `protobuf:"varint,5,opt,name=extent,proto3" json:"extent,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Tile_Layer) Reset() {
	*x = Tile_Layer{}
	mi := &file_vector_tile_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tile_Layer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tile_Layer) ProtoMessage() {}

func (x *Tile_Layer) ProtoReflect() protoreflect.Message {
	mi := &file_vector_tile_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tile_Layer.ProtoReflect.Descriptor instead.
func (*Tile_Layer) Descriptor() ([]byte, []int) {
	return file_vector_tile_proto_rawDescGZIP(), []int{0, 2}
}

func (x *Tile_Layer) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Tile_Layer) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tile_Layer) GetFeatures() []*Tile_Feature {
	if x != nil {
		return x.Features
	}
	return nil
}

func (x *Tile_Layer) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *Tile_Layer) GetValues() []*Tile_Value {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *Tile_Layer) GetExtent() uint32 {
	if x != nil {
		return x.Extent
	}
	return 0
}

var File_vector_tile_proto protoreflect.FileDescriptor

var file_vector_tile_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x74, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x74, 0x69, 0x6c, 0x65,
	0x22, 0xca, 0x05, 0x0a, 0x04, 0x54, 0x69, 0x6c, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x5f, 0x74, 0x69, 0x6c, 0x65, 0x2e, 0x54, 0x69, 0x6c, 0x65, 0x2e, 0x4c, 0x61, 0x79,
	0x65, 0x72, 0x52, 0x06, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x1a, 0xff, 0x01, 0x0a, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x0b, 0x66, 0x6c, 0x6f,
	0x61, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x48, 0x00,
	0x52, 0x0a, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c,
	0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x1d, 0x0a, 0x09, 0x69, 0x6e, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x1f, 0x0a, 0x0a, 0x75, 0x69, 0x6e, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x09, 0x75, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x1f, 0x0a, 0x0a, 0x73, 0x69, 0x6e, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x12, 0x48, 0x00, 0x52, 0x09, 0x73, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x1f, 0x0a, 0x0a, 0x62, 0x6f, 0x6f, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x81, 0x01, 0x0a,
	0x07, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0d, 0x42, 0x02, 0x10, 0x01, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73,
	0x12, 0x2e, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a,
	0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x74, 0x69, 0x6c, 0x65, 0x2e, 0x54, 0x69, 0x6c,
	0x65, 0x2e, 0x47, 0x65, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x1e, 0x0a, 0x08, 0x67, 0x65, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0d, 0x42, 0x02, 0x10, 0x01, 0x52, 0x08, 0x67, 0x65, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79,
	0x1a, 0xc9, 0x01, 0x0a, 0x05, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x66, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x5f, 0x74, 0x69, 0x6c, 0x65, 0x2e, 0x54, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x65,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b,
	0x65, 0x79, 0x73, 0x12, 0x2f, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x74, 0x69, 0x6c,
	0x65, 0x2e, 0x54, 0x69, 0x6c, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x3f, 0x0a, 0x08,
	0x47, 0x65, 0x6f, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x4f, 0x49, 0x4e, 0x54, 0x10, 0x01,
	0x12, 0x0e, 0x0a, 0x0a, 0x4c, 0x49, 0x4e, 0x45, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x10, 0x02,
	0x12, 0x0b, 0x0a, 0x07, 0x50, 0x4f, 0x4c, 0x59, 0x47, 0x4f, 0x4e, 0x10, 0x03, 0x42, 0x3c, 0x48,
	0x03, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x63,
	0x68, 0x63, 0x68, 0x76, 0x2f, 0x67, 0x65, 0x6f, 0x2f, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e,
	0x67, 0x2f, 0x6d, 0x76, 0x74, 0x2f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x74, 0x69, 0x6c, 0x65,
	0x3b, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_vector_tile_proto_rawDescOnce sync.Once
	file_vector_tile_proto_rawDescData []byte
)

func file_vector_tile_proto_rawDescGZIP() []byte {
	file_vector_tile_proto_rawDescOnce.Do(func() {
		file_vector_tile_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_vector_tile_proto_rawDesc), len(file_vector_tile_proto_rawDesc)))
	})
	return file_vector_tile_proto_rawDescData
}

var file_vector_tile_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_vector_tile_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_vector_tile_proto_goTypes = []any{
	(Tile_GeomType)(0),   // 0: vector_tile.Tile.GeomType
	(*Tile)(nil),         // 1: vector_tile.Tile
	(*Tile_Value)(nil),   // 2: vector_tile.Tile.Value
	(*Tile_Feature)(nil), // 3: vector_tile.Tile.Feature
	(*Tile_Layer)(nil),   // 4: vector_tile.Tile.Layer
}
var file_vector_tile_proto_depIdxs = []int32{
	4, // 0: vector_tile.Tile.layers:type_name -> vector_tile.Tile.Layer
	0, // 1: vector_tile.Tile.Feature.type:type_name -> vector_tile.Tile.GeomType
	3, // 2: vector_tile.Tile.Layer.features:type_name -> vector_tile.Tile.Feature
	2, // 3: vector_tile.Tile.Layer.values:type_name -> vector_tile.Tile.Value
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_vector_tile_proto_init() }
func file_vector_tile_proto_init() {
	if File_vector_tile_proto != nil {
		return
	}
	file_vector_tile_proto_msgTypes[1].OneofWrappers = []any{
		(*Tile_Value_StringValue)(nil),
		(*Tile_Value_FloatValue)(nil),
		(*Tile_Value_DoubleValue)(nil),
		(*Tile_Value_IntValue)(nil),
		(*Tile_Value_UintValue)(nil),
		(*Tile_Value_SintValue)(nil),
		(*Tile_Value_BoolValue)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_vector_tile_proto_rawDesc), len(file_vector_tile_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_vector_tile_proto_goTypes,
		DependencyIndexes: file_vector_tile_proto_depIdxs,
		EnumInfos:         file_vector_tile_proto_enumTypes,
		MessageInfos:      file_vector_tile_proto_msgTypes,
	}.Build()
	File_vector_tile_proto = out.File
	file_vector_tile_proto_goTypes = nil
	file_vector_tile_proto_depIdxs = nil
}
