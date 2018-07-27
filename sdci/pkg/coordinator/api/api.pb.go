// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api.proto

/*
	Package api is a generated protocol buffer package.

	It is generated from these files:
		api.proto

	It has these top-level messages:
		Job
		Recipe
		Build
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/golang/protobuf/ptypes/timestamp"

import time "time"

import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"

import types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Job struct {
	Recipe *Recipe `protobuf:"bytes,1,opt,name=recipe" json:"recipe,omitempty"`
	Build  *Build  `protobuf:"bytes,2,opt,name=build" json:"build,omitempty"`
}

func (m *Job) Reset()                    { *m = Job{} }
func (m *Job) String() string            { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()               {}
func (*Job) Descriptor() ([]byte, []int) { return fileDescriptorApi, []int{0} }

func (m *Job) GetRecipe() *Recipe {
	if m != nil {
		return m.Recipe
	}
	return nil
}

func (m *Job) GetBuild() *Build {
	if m != nil {
		return m.Build
	}
	return nil
}

type Recipe struct {
	Concurrency  int64  `protobuf:"varint,1,opt,name=concurrency,proto3" json:"concurrency,omitempty" yaml:"concurrency" db:"concurrency"`
	Clone        string `protobuf:"bytes,2,opt,name=clone,proto3" json:"clone,omitempty" yaml:"clone" db:"clone"`
	SlackWebhook string `protobuf:"bytes,3,opt,name=slackWebhook,proto3" json:"slackWebhook,omitempty" yaml:"slack_webhook" db:"slack_webhook"`
	GithubSecret string `protobuf:"bytes,4,opt,name=githubSecret,proto3" json:"githubSecret,omitempty" yaml:"github_secret" db:"github_secret"`
	Environment  string `protobuf:"bytes,5,opt,name=environment,proto3" json:"environment,omitempty" yaml:"environment" db:"environment"`
	Commands     string `protobuf:"bytes,6,opt,name=commands,proto3" json:"commands,omitempty" yaml:"commands" db:"commands"`
}

func (m *Recipe) Reset()                    { *m = Recipe{} }
func (m *Recipe) String() string            { return proto.CompactTextString(m) }
func (*Recipe) ProtoMessage()               {}
func (*Recipe) Descriptor() ([]byte, []int) { return fileDescriptorApi, []int{1} }

func (m *Recipe) GetConcurrency() int64 {
	if m != nil {
		return m.Concurrency
	}
	return 0
}

func (m *Recipe) GetClone() string {
	if m != nil {
		return m.Clone
	}
	return ""
}

func (m *Recipe) GetSlackWebhook() string {
	if m != nil {
		return m.SlackWebhook
	}
	return ""
}

func (m *Recipe) GetGithubSecret() string {
	if m != nil {
		return m.GithubSecret
	}
	return ""
}

func (m *Recipe) GetEnvironment() string {
	if m != nil {
		return m.Environment
	}
	return ""
}

func (m *Recipe) GetCommands() string {
	if m != nil {
		return m.Commands
	}
	return ""
}

type Build struct {
	ID            int64      `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty" yaml:"id" db:"id"`
	RepoFullName  string     `protobuf:"bytes,2,opt,name=repoFullName,proto3" json:"repoFullName,omitempty" yaml:"repo_full_name" db:"repo_full_name"`
	CommitHash    string     `protobuf:"bytes,3,opt,name=commitHash,proto3" json:"commitHash,omitempty" yaml:"commit_hash" db:"commit_hash"`
	CommitMessage string     `protobuf:"bytes,4,opt,name=commitMessage,proto3" json:"commitMessage,omitempty" yaml:"commit_message" db:"commit_message"`
	StartedAt     *time.Time `protobuf:"bytes,5,opt,name=startedAt,stdtime" json:"startedAt,omitempty" yaml:"started_at" db:"started_at"`
	Success       bool       `protobuf:"varint,6,opt,name=success,proto3" json:"success,omitempty" yaml:"success" db:"success"`
	Log           string     `protobuf:"bytes,7,opt,name=log,proto3" json:"log,omitempty" yaml:"log" db:"log"`
	CompletedAt   *time.Time `protobuf:"bytes,8,opt,name=completedAt,stdtime" json:"completedAt,omitempty" yaml:"completed_at" db:"completed_at"`
}

func (m *Build) Reset()                    { *m = Build{} }
func (m *Build) String() string            { return proto.CompactTextString(m) }
func (*Build) ProtoMessage()               {}
func (*Build) Descriptor() ([]byte, []int) { return fileDescriptorApi, []int{2} }

func (m *Build) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Build) GetRepoFullName() string {
	if m != nil {
		return m.RepoFullName
	}
	return ""
}

func (m *Build) GetCommitHash() string {
	if m != nil {
		return m.CommitHash
	}
	return ""
}

func (m *Build) GetCommitMessage() string {
	if m != nil {
		return m.CommitMessage
	}
	return ""
}

func (m *Build) GetStartedAt() *time.Time {
	if m != nil {
		return m.StartedAt
	}
	return nil
}

func (m *Build) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *Build) GetLog() string {
	if m != nil {
		return m.Log
	}
	return ""
}

func (m *Build) GetCompletedAt() *time.Time {
	if m != nil {
		return m.CompletedAt
	}
	return nil
}

func init() {
	proto.RegisterType((*Job)(nil), "api.Job")
	proto.RegisterType((*Recipe)(nil), "api.Recipe")
	proto.RegisterType((*Build)(nil), "api.Build")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Runner service

type RunnerClient interface {
	Run(ctx context.Context, opts ...grpc.CallOption) (Runner_RunClient, error)
}

type runnerClient struct {
	cc *grpc.ClientConn
}

func NewRunnerClient(cc *grpc.ClientConn) RunnerClient {
	return &runnerClient{cc}
}

func (c *runnerClient) Run(ctx context.Context, opts ...grpc.CallOption) (Runner_RunClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Runner_serviceDesc.Streams[0], c.cc, "/api.Runner/Run", opts...)
	if err != nil {
		return nil, err
	}
	x := &runnerRunClient{stream}
	return x, nil
}

type Runner_RunClient interface {
	Send(*Job) error
	Recv() (*Job, error)
	grpc.ClientStream
}

type runnerRunClient struct {
	grpc.ClientStream
}

func (x *runnerRunClient) Send(m *Job) error {
	return x.ClientStream.SendMsg(m)
}

func (x *runnerRunClient) Recv() (*Job, error) {
	m := new(Job)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Runner service

type RunnerServer interface {
	Run(Runner_RunServer) error
}

func RegisterRunnerServer(s *grpc.Server, srv RunnerServer) {
	s.RegisterService(&_Runner_serviceDesc, srv)
}

func _Runner_Run_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RunnerServer).Run(&runnerRunServer{stream})
}

type Runner_RunServer interface {
	Send(*Job) error
	Recv() (*Job, error)
	grpc.ServerStream
}

type runnerRunServer struct {
	grpc.ServerStream
}

func (x *runnerRunServer) Send(m *Job) error {
	return x.ServerStream.SendMsg(m)
}

func (x *runnerRunServer) Recv() (*Job, error) {
	m := new(Job)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Runner_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Runner",
	HandlerType: (*RunnerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Run",
			Handler:       _Runner_Run_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api.proto",
}

func (m *Job) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Job) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Recipe != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintApi(dAtA, i, uint64(m.Recipe.Size()))
		n1, err := m.Recipe.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Build != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintApi(dAtA, i, uint64(m.Build.Size()))
		n2, err := m.Build.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *Recipe) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Recipe) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Concurrency != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintApi(dAtA, i, uint64(m.Concurrency))
	}
	if len(m.Clone) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintApi(dAtA, i, uint64(len(m.Clone)))
		i += copy(dAtA[i:], m.Clone)
	}
	if len(m.SlackWebhook) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintApi(dAtA, i, uint64(len(m.SlackWebhook)))
		i += copy(dAtA[i:], m.SlackWebhook)
	}
	if len(m.GithubSecret) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintApi(dAtA, i, uint64(len(m.GithubSecret)))
		i += copy(dAtA[i:], m.GithubSecret)
	}
	if len(m.Environment) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintApi(dAtA, i, uint64(len(m.Environment)))
		i += copy(dAtA[i:], m.Environment)
	}
	if len(m.Commands) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintApi(dAtA, i, uint64(len(m.Commands)))
		i += copy(dAtA[i:], m.Commands)
	}
	return i, nil
}

func (m *Build) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Build) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintApi(dAtA, i, uint64(m.ID))
	}
	if len(m.RepoFullName) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintApi(dAtA, i, uint64(len(m.RepoFullName)))
		i += copy(dAtA[i:], m.RepoFullName)
	}
	if len(m.CommitHash) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintApi(dAtA, i, uint64(len(m.CommitHash)))
		i += copy(dAtA[i:], m.CommitHash)
	}
	if len(m.CommitMessage) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintApi(dAtA, i, uint64(len(m.CommitMessage)))
		i += copy(dAtA[i:], m.CommitMessage)
	}
	if m.StartedAt != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintApi(dAtA, i, uint64(types.SizeOfStdTime(*m.StartedAt)))
		n3, err := types.StdTimeMarshalTo(*m.StartedAt, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.Success {
		dAtA[i] = 0x30
		i++
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.Log) > 0 {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintApi(dAtA, i, uint64(len(m.Log)))
		i += copy(dAtA[i:], m.Log)
	}
	if m.CompletedAt != nil {
		dAtA[i] = 0x42
		i++
		i = encodeVarintApi(dAtA, i, uint64(types.SizeOfStdTime(*m.CompletedAt)))
		n4, err := types.StdTimeMarshalTo(*m.CompletedAt, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	return i, nil
}

func encodeVarintApi(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Job) Size() (n int) {
	var l int
	_ = l
	if m.Recipe != nil {
		l = m.Recipe.Size()
		n += 1 + l + sovApi(uint64(l))
	}
	if m.Build != nil {
		l = m.Build.Size()
		n += 1 + l + sovApi(uint64(l))
	}
	return n
}

func (m *Recipe) Size() (n int) {
	var l int
	_ = l
	if m.Concurrency != 0 {
		n += 1 + sovApi(uint64(m.Concurrency))
	}
	l = len(m.Clone)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	l = len(m.SlackWebhook)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	l = len(m.GithubSecret)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	l = len(m.Environment)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	l = len(m.Commands)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	return n
}

func (m *Build) Size() (n int) {
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovApi(uint64(m.ID))
	}
	l = len(m.RepoFullName)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	l = len(m.CommitHash)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	l = len(m.CommitMessage)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	if m.StartedAt != nil {
		l = types.SizeOfStdTime(*m.StartedAt)
		n += 1 + l + sovApi(uint64(l))
	}
	if m.Success {
		n += 2
	}
	l = len(m.Log)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	if m.CompletedAt != nil {
		l = types.SizeOfStdTime(*m.CompletedAt)
		n += 1 + l + sovApi(uint64(l))
	}
	return n
}

func sovApi(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozApi(x uint64) (n int) {
	return sovApi(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Job) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Job: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Job: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipe", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Recipe == nil {
				m.Recipe = &Recipe{}
			}
			if err := m.Recipe.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Build", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Build == nil {
				m.Build = &Build{}
			}
			if err := m.Build.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Recipe) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Recipe: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Recipe: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Concurrency", wireType)
			}
			m.Concurrency = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Concurrency |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Clone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Clone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlackWebhook", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SlackWebhook = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GithubSecret", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GithubSecret = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Environment", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Environment = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Commands", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Commands = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Build) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Build: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Build: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RepoFullName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RepoFullName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommitHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CommitHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommitMessage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CommitMessage = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.StartedAt == nil {
				m.StartedAt = new(time.Time)
			}
			if err := types.StdTimeUnmarshal(m.StartedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Success", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Success = bool(v != 0)
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Log", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Log = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CompletedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CompletedAt == nil {
				m.CompletedAt = new(time.Time)
			}
			if err := types.StdTimeUnmarshal(m.CompletedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthApi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipApi(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowApi
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowApi
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowApi
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthApi
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowApi
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipApi(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthApi = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowApi   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("api.proto", fileDescriptorApi) }

var fileDescriptorApi = []byte{
	// 438 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xae, 0xe3, 0xc6, 0x4d, 0x26, 0x80, 0xaa, 0x15, 0x07, 0xcb, 0x07, 0xa7, 0x32, 0x1c, 0xca,
	0x01, 0x17, 0x85, 0x3b, 0x12, 0x51, 0x85, 0x68, 0x05, 0x1c, 0x16, 0x24, 0xce, 0xeb, 0xcd, 0xe0,
	0x58, 0xb5, 0x77, 0xac, 0xdd, 0x35, 0x12, 0x6f, 0xc1, 0x91, 0xd7, 0xe0, 0x25, 0x10, 0x47, 0xde,
	0x00, 0x14, 0x5e, 0x04, 0x79, 0x37, 0x6e, 0x9d, 0x1b, 0xb7, 0xfd, 0x7e, 0x66, 0xbe, 0x19, 0x8f,
	0x61, 0x2e, 0xda, 0x2a, 0x6f, 0x35, 0x59, 0x62, 0xa1, 0x68, 0xab, 0xe4, 0x69, 0x59, 0xd9, 0x6d,
	0x57, 0xe4, 0x92, 0x9a, 0x8b, 0x92, 0x4a, 0xba, 0x70, 0x5a, 0xd1, 0x7d, 0x72, 0xc8, 0x01, 0xf7,
	0xf2, 0x35, 0xc9, 0xb2, 0x24, 0x2a, 0x6b, 0xbc, 0x73, 0xd9, 0xaa, 0x41, 0x63, 0x45, 0xd3, 0x7a,
	0x43, 0xf6, 0x06, 0xc2, 0x6b, 0x2a, 0xd8, 0x23, 0x88, 0x34, 0xca, 0xaa, 0xc5, 0x38, 0x38, 0x0b,
	0xce, 0x17, 0xab, 0x45, 0xde, 0xe7, 0x72, 0x47, 0xf1, 0xbd, 0xc4, 0xce, 0x60, 0x5a, 0x74, 0x55,
	0xbd, 0x89, 0x27, 0xce, 0x03, 0xce, 0xb3, 0xee, 0x19, 0xee, 0x85, 0xec, 0x47, 0x00, 0x11, 0x1f,
	0xcc, 0x0b, 0x49, 0x4a, 0x76, 0x5a, 0xa3, 0x92, 0x5f, 0x5c, 0xdb, 0x90, 0x8f, 0x29, 0xf6, 0x10,
	0xa6, 0xb2, 0x26, 0x85, 0xae, 0xdd, 0x9c, 0x7b, 0xc0, 0x32, 0xb8, 0x67, 0x6a, 0x21, 0x6f, 0x3e,
	0x62, 0xb1, 0x25, 0xba, 0x89, 0x43, 0x27, 0x1e, 0x70, 0xbd, 0xc7, 0x7f, 0x86, 0xf7, 0x28, 0x35,
	0xda, 0xf8, 0xd8, 0x7b, 0xc6, 0x5c, 0x9f, 0x8f, 0xea, 0x73, 0xa5, 0x49, 0x35, 0xa8, 0x6c, 0x3c,
	0x75, 0x96, 0x31, 0xc5, 0x12, 0x98, 0x49, 0x6a, 0x1a, 0xa1, 0x36, 0x26, 0x8e, 0x9c, 0x7c, 0x8b,
	0xb3, 0xef, 0x13, 0x98, 0xba, 0xcd, 0xd8, 0x03, 0x98, 0x5c, 0x5d, 0xee, 0xc7, 0x9f, 0x5c, 0x5d,
	0xf6, 0xd9, 0x1a, 0x5b, 0x7a, 0xd5, 0xd5, 0xf5, 0x3b, 0xd1, 0x0c, 0xc3, 0x1f, 0x70, 0x2c, 0x05,
	0xe8, 0x3b, 0x55, 0xf6, 0xb5, 0x30, 0xdb, 0xfd, 0x06, 0x23, 0x86, 0x3d, 0x86, 0xfb, 0x1e, 0xbd,
	0x45, 0x63, 0x44, 0x89, 0xfb, 0x05, 0x0e, 0x49, 0xf6, 0x02, 0xe6, 0xc6, 0x0a, 0x6d, 0x71, 0xf3,
	0xd2, 0xcf, 0xbf, 0x58, 0x25, 0xb9, 0xbf, 0x67, 0x3e, 0xdc, 0x33, 0xff, 0x30, 0xdc, 0x73, 0x7d,
	0xfc, 0xf5, 0xf7, 0x32, 0xe0, 0x77, 0x25, 0x2c, 0x86, 0x13, 0xd3, 0x49, 0x89, 0xc6, 0xaf, 0x37,
	0xe3, 0x03, 0x64, 0xa7, 0x10, 0xd6, 0x54, 0xc6, 0x27, 0x2e, 0xb5, 0x7f, 0xb2, 0x75, 0x7f, 0xad,
	0xa6, 0xad, 0xd1, 0xa7, 0xcd, 0xfe, 0x33, 0x6d, 0x5c, 0xb4, 0x7a, 0x02, 0x11, 0xef, 0x94, 0x42,
	0xcd, 0x96, 0x10, 0xf2, 0x4e, 0xb1, 0x99, 0xfb, 0x41, 0xae, 0xa9, 0x48, 0x6e, 0x5f, 0xd9, 0xd1,
	0x79, 0xf0, 0x2c, 0x58, 0x9f, 0xfe, 0xdc, 0xa5, 0xc1, 0xaf, 0x5d, 0x1a, 0xfc, 0xd9, 0xa5, 0xc1,
	0xb7, 0xbf, 0xe9, 0x51, 0x11, 0xb9, 0x8c, 0xe7, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x38, 0x4a,
	0xf5, 0x97, 0xf0, 0x02, 0x00, 0x00,
}
