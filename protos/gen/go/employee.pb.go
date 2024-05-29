// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.24.4
// source: employee.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Employee struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName  string                 `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName   string                 `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	PositionId string                 `protobuf:"bytes,4,opt,name=position_id,json=positionId,proto3" json:"position_id,omitempty"`
	CreatedAt  *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt  *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Employee) Reset() {
	*x = Employee{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Employee) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Employee) ProtoMessage() {}

func (x *Employee) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Employee.ProtoReflect.Descriptor instead.
func (*Employee) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{0}
}

func (x *Employee) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Employee) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *Employee) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *Employee) GetPositionId() string {
	if x != nil {
		return x.PositionId
	}
	return ""
}

func (x *Employee) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Employee) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type CreateEmployeeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FirstName    string `protobuf:"bytes,1,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName     string `protobuf:"bytes,2,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	PositionName string `protobuf:"bytes,3,opt,name=position_name,json=positionName,proto3" json:"position_name,omitempty"`
	Salary       int32  `protobuf:"varint,4,opt,name=salary,proto3" json:"salary,omitempty"`
}

func (x *CreateEmployeeRequest) Reset() {
	*x = CreateEmployeeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateEmployeeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateEmployeeRequest) ProtoMessage() {}

func (x *CreateEmployeeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateEmployeeRequest.ProtoReflect.Descriptor instead.
func (*CreateEmployeeRequest) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{1}
}

func (x *CreateEmployeeRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *CreateEmployeeRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *CreateEmployeeRequest) GetPositionName() string {
	if x != nil {
		return x.PositionName
	}
	return ""
}

func (x *CreateEmployeeRequest) GetSalary() int32 {
	if x != nil {
		return x.Salary
	}
	return 0
}

type GetEmployeeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EmployeeId string `protobuf:"bytes,1,opt,name=employee_id,json=employeeId,proto3" json:"employee_id,omitempty"`
}

func (x *GetEmployeeRequest) Reset() {
	*x = GetEmployeeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEmployeeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEmployeeRequest) ProtoMessage() {}

func (x *GetEmployeeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEmployeeRequest.ProtoReflect.Descriptor instead.
func (*GetEmployeeRequest) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{2}
}

func (x *GetEmployeeRequest) GetEmployeeId() string {
	if x != nil {
		return x.EmployeeId
	}
	return ""
}

type GetEmployeeListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cursor string `protobuf:"bytes,1,opt,name=cursor,proto3" json:"cursor,omitempty"`
}

func (x *GetEmployeeListRequest) Reset() {
	*x = GetEmployeeListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEmployeeListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEmployeeListRequest) ProtoMessage() {}

func (x *GetEmployeeListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEmployeeListRequest.ProtoReflect.Descriptor instead.
func (*GetEmployeeListRequest) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{3}
}

func (x *GetEmployeeListRequest) GetCursor() string {
	if x != nil {
		return x.Cursor
	}
	return ""
}

type GetEmployeeListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cursor    string      `protobuf:"bytes,1,opt,name=cursor,proto3" json:"cursor,omitempty"`
	Employees []*Employee `protobuf:"bytes,2,rep,name=employees,proto3" json:"employees,omitempty"`
}

func (x *GetEmployeeListResponse) Reset() {
	*x = GetEmployeeListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEmployeeListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEmployeeListResponse) ProtoMessage() {}

func (x *GetEmployeeListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEmployeeListResponse.ProtoReflect.Descriptor instead.
func (*GetEmployeeListResponse) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{4}
}

func (x *GetEmployeeListResponse) GetCursor() string {
	if x != nil {
		return x.Cursor
	}
	return ""
}

func (x *GetEmployeeListResponse) GetEmployees() []*Employee {
	if x != nil {
		return x.Employees
	}
	return nil
}

type UpdateEmployeeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EmployeeId string `protobuf:"bytes,1,opt,name=employee_id,json=employeeId,proto3" json:"employee_id,omitempty"`
	PositionId string `protobuf:"bytes,2,opt,name=position_id,json=positionId,proto3" json:"position_id,omitempty"`
	FirstName  string `protobuf:"bytes,3,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName   string `protobuf:"bytes,4,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Salary     int32  `protobuf:"varint,5,opt,name=salary,proto3" json:"salary,omitempty"`
}

func (x *UpdateEmployeeRequest) Reset() {
	*x = UpdateEmployeeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateEmployeeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateEmployeeRequest) ProtoMessage() {}

func (x *UpdateEmployeeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateEmployeeRequest.ProtoReflect.Descriptor instead.
func (*UpdateEmployeeRequest) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateEmployeeRequest) GetEmployeeId() string {
	if x != nil {
		return x.EmployeeId
	}
	return ""
}

func (x *UpdateEmployeeRequest) GetPositionId() string {
	if x != nil {
		return x.PositionId
	}
	return ""
}

func (x *UpdateEmployeeRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *UpdateEmployeeRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *UpdateEmployeeRequest) GetSalary() int32 {
	if x != nil {
		return x.Salary
	}
	return 0
}

type DeleteEmployeeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EmployeeId string `protobuf:"bytes,1,opt,name=employee_id,json=employeeId,proto3" json:"employee_id,omitempty"`
}

func (x *DeleteEmployeeRequest) Reset() {
	*x = DeleteEmployeeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteEmployeeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteEmployeeRequest) ProtoMessage() {}

func (x *DeleteEmployeeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteEmployeeRequest.ProtoReflect.Descriptor instead.
func (*DeleteEmployeeRequest) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteEmployeeRequest) GetEmployeeId() string {
	if x != nil {
		return x.EmployeeId
	}
	return ""
}

type DeleteEmployeeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteEmployeeResponse) Reset() {
	*x = DeleteEmployeeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteEmployeeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteEmployeeResponse) ProtoMessage() {}

func (x *DeleteEmployeeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteEmployeeResponse.ProtoReflect.Descriptor instead.
func (*DeleteEmployeeResponse) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{7}
}

type SignInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EmployeeId string `protobuf:"bytes,1,opt,name=employee_id,json=employeeId,proto3" json:"employee_id,omitempty"`
}

func (x *SignInRequest) Reset() {
	*x = SignInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignInRequest) ProtoMessage() {}

func (x *SignInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignInRequest.ProtoReflect.Descriptor instead.
func (*SignInRequest) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{8}
}

func (x *SignInRequest) GetEmployeeId() string {
	if x != nil {
		return x.EmployeeId
	}
	return ""
}

type SignInResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *SignInResponse) Reset() {
	*x = SignInResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_employee_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignInResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignInResponse) ProtoMessage() {}

func (x *SignInResponse) ProtoReflect() protoreflect.Message {
	mi := &file_employee_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignInResponse.ProtoReflect.Descriptor instead.
func (*SignInResponse) Descriptor() ([]byte, []int) {
	return file_employee_proto_rawDescGZIP(), []int{9}
}

func (x *SignInResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

var File_employee_proto protoreflect.FileDescriptor

var file_employee_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xed,
	0x01, 0x0a, 0x08, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66,
	0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c,
	0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x90,
	0x01, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73,
	0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69,
	0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x61, 0x6c,
	0x61, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x61, 0x6c, 0x61, 0x72,
	0x79, 0x22, 0x35, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x6d, 0x70, 0x6c, 0x6f,
	0x79, 0x65, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6d,
	0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x49, 0x64, 0x22, 0x30, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x45,
	0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x22, 0x66, 0x0a, 0x17, 0x47, 0x65,
	0x74, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x12, 0x33, 0x0a,
	0x09, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x45,
	0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x52, 0x09, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65,
	0x65, 0x73, 0x22, 0xad, 0x01, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x70,
	0x6c, 0x6f, 0x79, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b,
	0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x49, 0x64, 0x12, 0x1f, 0x0a,
	0x0b, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1d,
	0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x61,
	0x6c, 0x61, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x61, 0x6c, 0x61,
	0x72, 0x79, 0x22, 0x38, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x6d, 0x70, 0x6c,
	0x6f, 0x79, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x65,
	0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x49, 0x64, 0x22, 0x18, 0x0a, 0x16,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x30, 0x0a, 0x0d, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x6d, 0x70, 0x6c, 0x6f,
	0x79, 0x65, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6d,
	0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x49, 0x64, 0x22, 0x26, 0x0a, 0x0e, 0x53, 0x69, 0x67, 0x6e,
	0x49, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x32, 0xee, 0x03, 0x0a, 0x0f, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x4b, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x6d,
	0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x12, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f,
	0x76, 0x69, 0x65, 0x77, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x70, 0x6c, 0x6f,
	0x79, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x72, 0x65, 0x73,
	0x75, 0x6d, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65,
	0x65, 0x12, 0x45, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65,
	0x12, 0x1f, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x47,
	0x65, 0x74, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x15, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e,
	0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x12, 0x5c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x45,
	0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x23, 0x2e, 0x72, 0x65,
	0x73, 0x75, 0x6d, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x6d, 0x70,
	0x6c, 0x6f, 0x79, 0x65, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x24, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x47,
	0x65, 0x74, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4b, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x12, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6d,
	0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x6d, 0x70,
	0x6c, 0x6f, 0x79, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x72,
	0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x45, 0x6d, 0x70, 0x6c, 0x6f,
	0x79, 0x65, 0x65, 0x12, 0x59, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x6d, 0x70,
	0x6c, 0x6f, 0x79, 0x65, 0x65, 0x12, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x76,
	0x69, 0x65, 0x77, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x6d, 0x70, 0x6c, 0x6f, 0x79,
	0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x72, 0x65, 0x73, 0x75,
	0x6d, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x6d,
	0x70, 0x6c, 0x6f, 0x79, 0x65, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41,
	0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x12, 0x1a, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6d,
	0x65, 0x5f, 0x76, 0x69, 0x65, 0x77, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x76, 0x69,
	0x65, 0x77, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x56, 0x65, 0x72, 0x63, 0x65, 0x31, 0x31, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x2d,
	0x76, 0x69, 0x65, 0x77, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_employee_proto_rawDescOnce sync.Once
	file_employee_proto_rawDescData = file_employee_proto_rawDesc
)

func file_employee_proto_rawDescGZIP() []byte {
	file_employee_proto_rawDescOnce.Do(func() {
		file_employee_proto_rawDescData = protoimpl.X.CompressGZIP(file_employee_proto_rawDescData)
	})
	return file_employee_proto_rawDescData
}

var file_employee_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_employee_proto_goTypes = []interface{}{
	(*Employee)(nil),                // 0: resume_view.Employee
	(*CreateEmployeeRequest)(nil),   // 1: resume_view.CreateEmployeeRequest
	(*GetEmployeeRequest)(nil),      // 2: resume_view.GetEmployeeRequest
	(*GetEmployeeListRequest)(nil),  // 3: resume_view.GetEmployeeListRequest
	(*GetEmployeeListResponse)(nil), // 4: resume_view.GetEmployeeListResponse
	(*UpdateEmployeeRequest)(nil),   // 5: resume_view.UpdateEmployeeRequest
	(*DeleteEmployeeRequest)(nil),   // 6: resume_view.DeleteEmployeeRequest
	(*DeleteEmployeeResponse)(nil),  // 7: resume_view.DeleteEmployeeResponse
	(*SignInRequest)(nil),           // 8: resume_view.SignInRequest
	(*SignInResponse)(nil),          // 9: resume_view.SignInResponse
	(*timestamppb.Timestamp)(nil),   // 10: google.protobuf.Timestamp
}
var file_employee_proto_depIdxs = []int32{
	10, // 0: resume_view.Employee.created_at:type_name -> google.protobuf.Timestamp
	10, // 1: resume_view.Employee.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 2: resume_view.GetEmployeeListResponse.employees:type_name -> resume_view.Employee
	1,  // 3: resume_view.EmployeeService.CreateEmployee:input_type -> resume_view.CreateEmployeeRequest
	2,  // 4: resume_view.EmployeeService.GetEmployee:input_type -> resume_view.GetEmployeeRequest
	3,  // 5: resume_view.EmployeeService.GetEmployeeList:input_type -> resume_view.GetEmployeeListRequest
	5,  // 6: resume_view.EmployeeService.UpdateEmployee:input_type -> resume_view.UpdateEmployeeRequest
	6,  // 7: resume_view.EmployeeService.DeleteEmployee:input_type -> resume_view.DeleteEmployeeRequest
	8,  // 8: resume_view.EmployeeService.SignIn:input_type -> resume_view.SignInRequest
	0,  // 9: resume_view.EmployeeService.CreateEmployee:output_type -> resume_view.Employee
	0,  // 10: resume_view.EmployeeService.GetEmployee:output_type -> resume_view.Employee
	4,  // 11: resume_view.EmployeeService.GetEmployeeList:output_type -> resume_view.GetEmployeeListResponse
	0,  // 12: resume_view.EmployeeService.UpdateEmployee:output_type -> resume_view.Employee
	7,  // 13: resume_view.EmployeeService.DeleteEmployee:output_type -> resume_view.DeleteEmployeeResponse
	9,  // 14: resume_view.EmployeeService.SignIn:output_type -> resume_view.SignInResponse
	9,  // [9:15] is the sub-list for method output_type
	3,  // [3:9] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_employee_proto_init() }
func file_employee_proto_init() {
	if File_employee_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_employee_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Employee); i {
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
		file_employee_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateEmployeeRequest); i {
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
		file_employee_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEmployeeRequest); i {
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
		file_employee_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEmployeeListRequest); i {
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
		file_employee_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEmployeeListResponse); i {
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
		file_employee_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateEmployeeRequest); i {
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
		file_employee_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteEmployeeRequest); i {
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
		file_employee_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteEmployeeResponse); i {
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
		file_employee_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignInRequest); i {
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
		file_employee_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignInResponse); i {
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
			RawDescriptor: file_employee_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_employee_proto_goTypes,
		DependencyIndexes: file_employee_proto_depIdxs,
		MessageInfos:      file_employee_proto_msgTypes,
	}.Build()
	File_employee_proto = out.File
	file_employee_proto_rawDesc = nil
	file_employee_proto_goTypes = nil
	file_employee_proto_depIdxs = nil
}
