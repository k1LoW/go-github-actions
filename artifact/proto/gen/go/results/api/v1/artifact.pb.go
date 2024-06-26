// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: results/api/v1/artifact.proto

package apiv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateArtifactRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkflowRunBackendId    string                 `protobuf:"bytes,1,opt,name=workflow_run_backend_id,json=workflowRunBackendId,proto3" json:"workflow_run_backend_id,omitempty"`
	WorkflowJobRunBackendId string                 `protobuf:"bytes,2,opt,name=workflow_job_run_backend_id,json=workflowJobRunBackendId,proto3" json:"workflow_job_run_backend_id,omitempty"`
	Name                    string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	ExpiresAt               *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
	Version                 int32                  `protobuf:"varint,5,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *CreateArtifactRequest) Reset() {
	*x = CreateArtifactRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_results_api_v1_artifact_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateArtifactRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateArtifactRequest) ProtoMessage() {}

func (x *CreateArtifactRequest) ProtoReflect() protoreflect.Message {
	mi := &file_results_api_v1_artifact_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateArtifactRequest.ProtoReflect.Descriptor instead.
func (*CreateArtifactRequest) Descriptor() ([]byte, []int) {
	return file_results_api_v1_artifact_proto_rawDescGZIP(), []int{0}
}

func (x *CreateArtifactRequest) GetWorkflowRunBackendId() string {
	if x != nil {
		return x.WorkflowRunBackendId
	}
	return ""
}

func (x *CreateArtifactRequest) GetWorkflowJobRunBackendId() string {
	if x != nil {
		return x.WorkflowJobRunBackendId
	}
	return ""
}

func (x *CreateArtifactRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateArtifactRequest) GetExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

func (x *CreateArtifactRequest) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

type CreateArtifactResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok              bool   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	SignedUploadUrl string `protobuf:"bytes,2,opt,name=signed_upload_url,json=signedUploadUrl,proto3" json:"signed_upload_url,omitempty"`
}

func (x *CreateArtifactResponse) Reset() {
	*x = CreateArtifactResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_results_api_v1_artifact_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateArtifactResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateArtifactResponse) ProtoMessage() {}

func (x *CreateArtifactResponse) ProtoReflect() protoreflect.Message {
	mi := &file_results_api_v1_artifact_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateArtifactResponse.ProtoReflect.Descriptor instead.
func (*CreateArtifactResponse) Descriptor() ([]byte, []int) {
	return file_results_api_v1_artifact_proto_rawDescGZIP(), []int{1}
}

func (x *CreateArtifactResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *CreateArtifactResponse) GetSignedUploadUrl() string {
	if x != nil {
		return x.SignedUploadUrl
	}
	return ""
}

type FinalizeArtifactRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkflowRunBackendId    string                  `protobuf:"bytes,1,opt,name=workflow_run_backend_id,json=workflowRunBackendId,proto3" json:"workflow_run_backend_id,omitempty"`
	WorkflowJobRunBackendId string                  `protobuf:"bytes,2,opt,name=workflow_job_run_backend_id,json=workflowJobRunBackendId,proto3" json:"workflow_job_run_backend_id,omitempty"`
	Name                    string                  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Size                    int64                   `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	Hash                    *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *FinalizeArtifactRequest) Reset() {
	*x = FinalizeArtifactRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_results_api_v1_artifact_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinalizeArtifactRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinalizeArtifactRequest) ProtoMessage() {}

func (x *FinalizeArtifactRequest) ProtoReflect() protoreflect.Message {
	mi := &file_results_api_v1_artifact_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinalizeArtifactRequest.ProtoReflect.Descriptor instead.
func (*FinalizeArtifactRequest) Descriptor() ([]byte, []int) {
	return file_results_api_v1_artifact_proto_rawDescGZIP(), []int{2}
}

func (x *FinalizeArtifactRequest) GetWorkflowRunBackendId() string {
	if x != nil {
		return x.WorkflowRunBackendId
	}
	return ""
}

func (x *FinalizeArtifactRequest) GetWorkflowJobRunBackendId() string {
	if x != nil {
		return x.WorkflowJobRunBackendId
	}
	return ""
}

func (x *FinalizeArtifactRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FinalizeArtifactRequest) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *FinalizeArtifactRequest) GetHash() *wrapperspb.StringValue {
	if x != nil {
		return x.Hash
	}
	return nil
}

type FinalizeArtifactResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok         bool  `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	ArtifactId int64 `protobuf:"varint,2,opt,name=artifact_id,json=artifactId,proto3" json:"artifact_id,omitempty"`
}

func (x *FinalizeArtifactResponse) Reset() {
	*x = FinalizeArtifactResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_results_api_v1_artifact_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinalizeArtifactResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinalizeArtifactResponse) ProtoMessage() {}

func (x *FinalizeArtifactResponse) ProtoReflect() protoreflect.Message {
	mi := &file_results_api_v1_artifact_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinalizeArtifactResponse.ProtoReflect.Descriptor instead.
func (*FinalizeArtifactResponse) Descriptor() ([]byte, []int) {
	return file_results_api_v1_artifact_proto_rawDescGZIP(), []int{3}
}

func (x *FinalizeArtifactResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *FinalizeArtifactResponse) GetArtifactId() int64 {
	if x != nil {
		return x.ArtifactId
	}
	return 0
}

var File_results_api_v1_artifact_proto protoreflect.FileDescriptor

var file_results_api_v1_artifact_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x1d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xf5, 0x01, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61,
	0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x17, 0x77, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x72, 0x75, 0x6e, 0x5f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e,
	0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x77, 0x6f, 0x72, 0x6b,
	0x66, 0x6c, 0x6f, 0x77, 0x52, 0x75, 0x6e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x49, 0x64,
	0x12, 0x3c, 0x0a, 0x1b, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x6a, 0x6f, 0x62,
	0x5f, 0x72, 0x75, 0x6e, 0x5f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x4a,
	0x6f, 0x62, 0x52, 0x75, 0x6e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x49, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x54, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f,
	0x6b, 0x12, 0x2a, 0x0a, 0x11, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x75, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x55, 0x72, 0x6c, 0x22, 0xe8, 0x01,
	0x0a, 0x17, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61,
	0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x17, 0x77, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x72, 0x75, 0x6e, 0x5f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e,
	0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x77, 0x6f, 0x72, 0x6b,
	0x66, 0x6c, 0x6f, 0x77, 0x52, 0x75, 0x6e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x49, 0x64,
	0x12, 0x3c, 0x0a, 0x1b, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x6a, 0x6f, 0x62,
	0x5f, 0x72, 0x75, 0x6e, 0x5f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x4a,
	0x6f, 0x62, 0x52, 0x75, 0x6e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x49, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x22, 0x4b, 0x0a, 0x18, 0x46, 0x69, 0x6e, 0x61,
	0x6c, 0x69, 0x7a, 0x65, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x02, 0x6f, 0x6b, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x61, 0x72, 0x74, 0x69, 0x66,
	0x61, 0x63, 0x74, 0x49, 0x64, 0x32, 0x96, 0x02, 0x0a, 0x0f, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61,
	0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x7d, 0x0a, 0x0e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x12, 0x34, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x35, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x83, 0x01, 0x0a, 0x10, 0x46, 0x69, 0x6e,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x12, 0x36, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69,
	0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x37, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x41, 0x72,
	0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x9a,
	0x02, 0x0a, 0x21, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6b, 0x31, 0x4c, 0x6f, 0x57, 0x2f, 0x67, 0x6f, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2d, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x72, 0x74, 0x69, 0x66, 0x61,
	0x63, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x61,
	0x70, 0x69, 0x76, 0x31, 0xa2, 0x02, 0x04, 0x47, 0x41, 0x52, 0x41, 0xaa, 0x02, 0x1d, 0x47, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x2e, 0x41, 0x70, 0x69, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1d, 0x47, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x5c, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5c, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x5c, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x29, 0x47, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x5c, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5c, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x5c, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x21, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x3a, 0x3a, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3a, 0x3a, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x73, 0x3a, 0x3a, 0x41, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_results_api_v1_artifact_proto_rawDescOnce sync.Once
	file_results_api_v1_artifact_proto_rawDescData = file_results_api_v1_artifact_proto_rawDesc
)

func file_results_api_v1_artifact_proto_rawDescGZIP() []byte {
	file_results_api_v1_artifact_proto_rawDescOnce.Do(func() {
		file_results_api_v1_artifact_proto_rawDescData = protoimpl.X.CompressGZIP(file_results_api_v1_artifact_proto_rawDescData)
	})
	return file_results_api_v1_artifact_proto_rawDescData
}

var file_results_api_v1_artifact_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_results_api_v1_artifact_proto_goTypes = []interface{}{
	(*CreateArtifactRequest)(nil),    // 0: github.actions.results.api.v1.CreateArtifactRequest
	(*CreateArtifactResponse)(nil),   // 1: github.actions.results.api.v1.CreateArtifactResponse
	(*FinalizeArtifactRequest)(nil),  // 2: github.actions.results.api.v1.FinalizeArtifactRequest
	(*FinalizeArtifactResponse)(nil), // 3: github.actions.results.api.v1.FinalizeArtifactResponse
	(*timestamppb.Timestamp)(nil),    // 4: google.protobuf.Timestamp
	(*wrapperspb.StringValue)(nil),   // 5: google.protobuf.StringValue
}
var file_results_api_v1_artifact_proto_depIdxs = []int32{
	4, // 0: github.actions.results.api.v1.CreateArtifactRequest.expires_at:type_name -> google.protobuf.Timestamp
	5, // 1: github.actions.results.api.v1.FinalizeArtifactRequest.hash:type_name -> google.protobuf.StringValue
	0, // 2: github.actions.results.api.v1.ArtifactService.CreateArtifact:input_type -> github.actions.results.api.v1.CreateArtifactRequest
	2, // 3: github.actions.results.api.v1.ArtifactService.FinalizeArtifact:input_type -> github.actions.results.api.v1.FinalizeArtifactRequest
	1, // 4: github.actions.results.api.v1.ArtifactService.CreateArtifact:output_type -> github.actions.results.api.v1.CreateArtifactResponse
	3, // 5: github.actions.results.api.v1.ArtifactService.FinalizeArtifact:output_type -> github.actions.results.api.v1.FinalizeArtifactResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_results_api_v1_artifact_proto_init() }
func file_results_api_v1_artifact_proto_init() {
	if File_results_api_v1_artifact_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_results_api_v1_artifact_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateArtifactRequest); i {
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
		file_results_api_v1_artifact_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateArtifactResponse); i {
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
		file_results_api_v1_artifact_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinalizeArtifactRequest); i {
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
		file_results_api_v1_artifact_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinalizeArtifactResponse); i {
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
			RawDescriptor: file_results_api_v1_artifact_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_results_api_v1_artifact_proto_goTypes,
		DependencyIndexes: file_results_api_v1_artifact_proto_depIdxs,
		MessageInfos:      file_results_api_v1_artifact_proto_msgTypes,
	}.Build()
	File_results_api_v1_artifact_proto = out.File
	file_results_api_v1_artifact_proto_rawDesc = nil
	file_results_api_v1_artifact_proto_goTypes = nil
	file_results_api_v1_artifact_proto_depIdxs = nil
}
