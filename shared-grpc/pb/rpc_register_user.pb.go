// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.12.4
// source: rpc_register_user.proto

package pb

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type RegisterUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fullname           string `protobuf:"bytes,1,opt,name=fullname,proto3" json:"fullname,omitempty"`
	Email              string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Password           string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	PaydUsername       string `protobuf:"bytes,4,opt,name=payd_username,json=paydUsername,proto3" json:"payd_username,omitempty"`
	PaydPasswordApiKey string `protobuf:"bytes,5,opt,name=payd_password_api_key,json=paydPasswordApiKey,proto3" json:"payd_password_api_key,omitempty"`
	PaydUsernameApiKey string `protobuf:"bytes,6,opt,name=payd_username_api_key,json=paydUsernameApiKey,proto3" json:"payd_username_api_key,omitempty"`
}

func (x *RegisterUserRequest) Reset() {
	*x = RegisterUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_register_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserRequest) ProtoMessage() {}

func (x *RegisterUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_register_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserRequest.ProtoReflect.Descriptor instead.
func (*RegisterUserRequest) Descriptor() ([]byte, []int) {
	return file_rpc_register_user_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterUserRequest) GetFullname() string {
	if x != nil {
		return x.Fullname
	}
	return ""
}

func (x *RegisterUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RegisterUserRequest) GetPaydUsername() string {
	if x != nil {
		return x.PaydUsername
	}
	return ""
}

func (x *RegisterUserRequest) GetPaydPasswordApiKey() string {
	if x != nil {
		return x.PaydPasswordApiKey
	}
	return ""
}

func (x *RegisterUserRequest) GetPaydUsernameApiKey() string {
	if x != nil {
		return x.PaydUsernameApiKey
	}
	return ""
}

type RegisterUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fullname  string               `protobuf:"bytes,1,opt,name=fullname,proto3" json:"fullname,omitempty"`
	Email     string               `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *RegisterUserResponse) Reset() {
	*x = RegisterUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_register_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserResponse) ProtoMessage() {}

func (x *RegisterUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_register_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserResponse.ProtoReflect.Descriptor instead.
func (*RegisterUserResponse) Descriptor() ([]byte, []int) {
	return file_rpc_register_user_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterUserResponse) GetFullname() string {
	if x != nil {
		return x.Fullname
	}
	return ""
}

func (x *RegisterUserResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterUserResponse) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

var File_rpc_register_user_proto protoreflect.FileDescriptor

var file_rpc_register_user_proto_rawDesc = []byte{
	0x0a, 0x17, 0x72, 0x70, 0x63, 0x5f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xee,
	0x01, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x61, 0x79, 0x64, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61, 0x79,
	0x64, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x15, 0x70, 0x61, 0x79,
	0x64, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x61, 0x70, 0x69, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x61, 0x79, 0x64, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x41, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x12, 0x31, 0x0a, 0x15,
	0x70, 0x61, 0x79, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x61, 0x70,
	0x69, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x61, 0x79,
	0x64, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x41, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x22,
	0x83, 0x01, 0x0a, 0x14, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6c, 0x6c,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x45, 0x6d, 0x69, 0x6c, 0x69, 0x6f, 0x43, 0x6c, 0x69, 0x66, 0x66, 0x2f,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x70, 0x6f, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x2d,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2d, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_register_user_proto_rawDescOnce sync.Once
	file_rpc_register_user_proto_rawDescData = file_rpc_register_user_proto_rawDesc
)

func file_rpc_register_user_proto_rawDescGZIP() []byte {
	file_rpc_register_user_proto_rawDescOnce.Do(func() {
		file_rpc_register_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_register_user_proto_rawDescData)
	})
	return file_rpc_register_user_proto_rawDescData
}

var file_rpc_register_user_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_register_user_proto_goTypes = []interface{}{
	(*RegisterUserRequest)(nil),  // 0: pb.RegisterUserRequest
	(*RegisterUserResponse)(nil), // 1: pb.RegisterUserResponse
	(*timestamp.Timestamp)(nil),  // 2: google.protobuf.Timestamp
}
var file_rpc_register_user_proto_depIdxs = []int32{
	2, // 0: pb.RegisterUserResponse.created_at:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_register_user_proto_init() }
func file_rpc_register_user_proto_init() {
	if File_rpc_register_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_register_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterUserRequest); i {
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
		file_rpc_register_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterUserResponse); i {
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
			RawDescriptor: file_rpc_register_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_register_user_proto_goTypes,
		DependencyIndexes: file_rpc_register_user_proto_depIdxs,
		MessageInfos:      file_rpc_register_user_proto_msgTypes,
	}.Build()
	File_rpc_register_user_proto = out.File
	file_rpc_register_user_proto_rawDesc = nil
	file_rpc_register_user_proto_goTypes = nil
	file_rpc_register_user_proto_depIdxs = nil
}