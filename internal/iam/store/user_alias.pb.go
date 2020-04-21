// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: user_alias.proto

package store

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

type UserAlias struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @inject_tag: gorm:"primary_key"
	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" gorm:"primary_key"`
	// create_time from the RDBMS
	// @inject_tag: `gorm:"default:current_timestamp"`
	CreateTime *Timestamp `protobuf:"bytes,2,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty" gorm:"default:current_timestamp"`
	// update_time from the RDBMS
	// @inject_tag: `gorm:"default:current_timestamp"`
	UpdateTime *Timestamp `protobuf:"bytes,3,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty" gorm:"default:current_timestamp"`
	// public_id is used to access the UserAlias via an API
	PublicId string `protobuf:"bytes,4,opt,name=public_id,json=publicId,proto3" json:"public_id,omitempty"`
	// FriendlyName is the optional friendly name used to
	// access the UserAlias via an API
	// @inject_tag: `gorm:"default:null"`
	FriendlyName string `protobuf:"bytes,5,opt,name=friendly_name,json=friendlyName,proto3" json:"friendly_name,omitempty" gorm:"default:null"`
	Name         string `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	// primary scope id of the PrimaryScope
	// @inject_tag: `gorm:"default:null"`
	PrimaryScopeId uint32 `protobuf:"varint,7,opt,name=primary_scope_id,json=primaryScopeId,proto3" json:"primary_scope_id,omitempty" gorm:"default:null"`
	// @inject_tag: gorm:"foreignkey:ScopeId"
	PrimaryScope *Scope `protobuf:"bytes,8,opt,name=primary_scope,json=primaryScope,proto3" json:"primary_scope,omitempty" gorm:"foreignkey:ScopeId"`
	// user_id is the id (primary key in the user aggregate) of the user that
	// owns this UserAlias
	// @inject_tag: `gorm:"default:0"`
	UserId uint32 `protobuf:"varint,9,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" gorm:"default:0"`
	// auth_method_id is the AuthMethod of this UserAlias.
	AuthMethodId uint32 `protobuf:"varint,10,opt,name=auth_method_id,json=authMethodId,proto3" json:"auth_method_id,omitempty"`
	// disabled is by default false and allows a UserAlias to be marked disabled.
	// which means among other things, the UserAlias cannot authen/authz
	Disabled bool `protobuf:"varint,11,opt,name=disabled,proto3" json:"disabled,omitempty"`
}

func (x *UserAlias) Reset() {
	*x = UserAlias{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_alias_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserAlias) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserAlias) ProtoMessage() {}

func (x *UserAlias) ProtoReflect() protoreflect.Message {
	mi := &file_user_alias_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserAlias.ProtoReflect.Descriptor instead.
func (*UserAlias) Descriptor() ([]byte, []int) {
	return file_user_alias_proto_rawDescGZIP(), []int{0}
}

func (x *UserAlias) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserAlias) GetCreateTime() *Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *UserAlias) GetUpdateTime() *Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

func (x *UserAlias) GetPublicId() string {
	if x != nil {
		return x.PublicId
	}
	return ""
}

func (x *UserAlias) GetFriendlyName() string {
	if x != nil {
		return x.FriendlyName
	}
	return ""
}

func (x *UserAlias) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserAlias) GetPrimaryScopeId() uint32 {
	if x != nil {
		return x.PrimaryScopeId
	}
	return 0
}

func (x *UserAlias) GetPrimaryScope() *Scope {
	if x != nil {
		return x.PrimaryScope
	}
	return nil
}

func (x *UserAlias) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserAlias) GetAuthMethodId() uint32 {
	if x != nil {
		return x.AuthMethodId
	}
	return 0
}

func (x *UserAlias) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

var File_user_alias_proto protoreflect.FileDescriptor

var file_user_alias_proto_rawDesc = []byte{
	0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x2c, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2e, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x74, 0x6f, 0x77, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x6c, 0x65, 0x72, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x1a, 0x0f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0b, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84,
	0x04, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x41, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x58, 0x0a, 0x0b,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x37, 0x2e, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2e, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x74, 0x6f, 0x77, 0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x6c, 0x65, 0x72, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x58, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x68, 0x61,
	0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2e, 0x77, 0x61, 0x74, 0x63, 0x68, 0x74, 0x6f, 0x77,
	0x65, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x69, 0x61,
	0x6d, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x49, 0x64, 0x12, 0x23, 0x0a,
	0x0d, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72,
	0x79, 0x5f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0e, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x49, 0x64,
	0x12, 0x58, 0x0a, 0x0d, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x5f, 0x73, 0x63, 0x6f, 0x70,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63,
	0x6f, 0x72, 0x70, 0x2e, 0x77, 0x61, 0x74, 0x63, 0x68, 0x74, 0x6f, 0x77, 0x65, 0x72, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x52, 0x0c, 0x70, 0x72,
	0x69, 0x6d, 0x61, 0x72, 0x79, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x61, 0x75, 0x74,
	0x68, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73,
	0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x73,
	0x61, 0x62, 0x6c, 0x65, 0x64, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2f, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x74, 0x6f, 0x77, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x3b, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_alias_proto_rawDescOnce sync.Once
	file_user_alias_proto_rawDescData = file_user_alias_proto_rawDesc
)

func file_user_alias_proto_rawDescGZIP() []byte {
	file_user_alias_proto_rawDescOnce.Do(func() {
		file_user_alias_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_alias_proto_rawDescData)
	})
	return file_user_alias_proto_rawDescData
}

var file_user_alias_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_user_alias_proto_goTypes = []interface{}{
	(*UserAlias)(nil), // 0: hashicorp.watchtower.controller.iam.store.v1.UserAlias
	(*Timestamp)(nil), // 1: hashicorp.watchtower.controller.iam.store.v1.Timestamp
	(*Scope)(nil),     // 2: hashicorp.watchtower.controller.iam.store.v1.Scope
}
var file_user_alias_proto_depIdxs = []int32{
	1, // 0: hashicorp.watchtower.controller.iam.store.v1.UserAlias.create_time:type_name -> hashicorp.watchtower.controller.iam.store.v1.Timestamp
	1, // 1: hashicorp.watchtower.controller.iam.store.v1.UserAlias.update_time:type_name -> hashicorp.watchtower.controller.iam.store.v1.Timestamp
	2, // 2: hashicorp.watchtower.controller.iam.store.v1.UserAlias.primary_scope:type_name -> hashicorp.watchtower.controller.iam.store.v1.Scope
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_user_alias_proto_init() }
func file_user_alias_proto_init() {
	if File_user_alias_proto != nil {
		return
	}
	file_timestamp_proto_init()
	file_scope_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_user_alias_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserAlias); i {
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
			RawDescriptor: file_user_alias_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_user_alias_proto_goTypes,
		DependencyIndexes: file_user_alias_proto_depIdxs,
		MessageInfos:      file_user_alias_proto_msgTypes,
	}.Build()
	File_user_alias_proto = out.File
	file_user_alias_proto_rawDesc = nil
	file_user_alias_proto_goTypes = nil
	file_user_alias_proto_depIdxs = nil
}
