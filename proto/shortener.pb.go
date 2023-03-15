// Code generated by protoc-gen-go. DO NOT EDIT.
// source: shortener.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Empty simple is stub parameter.
type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

// The request message for ShortenURL.
type ShortenRequest struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShortenRequest) Reset()         { *m = ShortenRequest{} }
func (m *ShortenRequest) String() string { return proto.CompactTextString(m) }
func (*ShortenRequest) ProtoMessage()    {}
func (*ShortenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{1}
}

func (m *ShortenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShortenRequest.Unmarshal(m, b)
}
func (m *ShortenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShortenRequest.Marshal(b, m, deterministic)
}
func (m *ShortenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShortenRequest.Merge(m, src)
}
func (m *ShortenRequest) XXX_Size() int {
	return xxx_messageInfo_ShortenRequest.Size(m)
}
func (m *ShortenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShortenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShortenRequest proto.InternalMessageInfo

func (m *ShortenRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

// The response message for ShortenURL.
type ShortenResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Result               string   `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
	Success              bool     `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
}

func (m *ShortenResponse) Reset()         { *m = ShortenResponse{} }
func (m *ShortenResponse) String() string { return proto.CompactTextString(m) }
func (*ShortenResponse) ProtoMessage()    {}
func (*ShortenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{2}
}

func (m *ShortenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShortenResponse.Unmarshal(m, b)
}
func (m *ShortenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShortenResponse.Marshal(b, m, deterministic)
}
func (m *ShortenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShortenResponse.Merge(m, src)
}
func (m *ShortenResponse) XXX_Size() int {
	return xxx_messageInfo_ShortenResponse.Size(m)
}
func (m *ShortenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ShortenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ShortenResponse proto.InternalMessageInfo

func (m *ShortenResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *ShortenResponse) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func (m *ShortenResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

// The request message for GetStat.
type StatRequest struct {
	Shortuid             string   `protobuf:"bytes,1,opt,name=shortuid,proto3" json:"shortuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatRequest) Reset()         { *m = StatRequest{} }
func (m *StatRequest) String() string { return proto.CompactTextString(m) }
func (*StatRequest) ProtoMessage()    {}
func (*StatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{3}
}

func (m *StatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatRequest.Unmarshal(m, b)
}
func (m *StatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatRequest.Marshal(b, m, deterministic)
}
func (m *StatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatRequest.Merge(m, src)
}
func (m *StatRequest) XXX_Size() int {
	return xxx_messageInfo_StatRequest.Size(m)
}
func (m *StatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatRequest proto.InternalMessageInfo

func (m *StatRequest) GetShortuid() string {
	if m != nil {
		return m.Shortuid
	}
	return ""
}

// The response message for GetStat.
type StatResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	Shorturl             string   `protobuf:"bytes,1,opt,name=shorturl,proto3" json:"shorturl,omitempty"`
	Longurl              string   `protobuf:"bytes,2,opt,name=longurl,proto3" json:"longurl,omitempty"`
	CreatedAt            string   `protobuf:"bytes,3,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	XXX_unrecognized     []byte   `json:"-"`
	Usage                uint64   `protobuf:"varint,4,opt,name=usage,proto3" json:"usage,omitempty"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatResponse) Reset()         { *m = StatResponse{} }
func (m *StatResponse) String() string { return proto.CompactTextString(m) }
func (*StatResponse) ProtoMessage()    {}
func (*StatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{4}
}

func (m *StatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatResponse.Unmarshal(m, b)
}
func (m *StatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatResponse.Marshal(b, m, deterministic)
}
func (m *StatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatResponse.Merge(m, src)
}
func (m *StatResponse) XXX_Size() int {
	return xxx_messageInfo_StatResponse.Size(m)
}
func (m *StatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatResponse proto.InternalMessageInfo

func (m *StatResponse) GetShorturl() string {
	if m != nil {
		return m.Shorturl
	}
	return ""
}

func (m *StatResponse) GetLongurl() string {
	if m != nil {
		return m.Longurl
	}
	return ""
}

func (m *StatResponse) GetCreatedAt() string {
	if m != nil {
		return m.CreatedAt
	}
	return ""
}

func (m *StatResponse) GetUsage() uint64 {
	if m != nil {
		return m.Usage
	}
	return 0
}

// The response message for GetUsersURL.
type UsersURL struct {
	Urls                 []*UserURL `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *UsersURL) Reset()         { *m = UsersURL{} }
func (m *UsersURL) String() string { return proto.CompactTextString(m) }
func (*UsersURL) ProtoMessage()    {}
func (*UsersURL) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{5}
}

func (m *UsersURL) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UsersURL.Unmarshal(m, b)
}
func (m *UsersURL) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UsersURL.Marshal(b, m, deterministic)
}
func (m *UsersURL) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UsersURL.Merge(m, src)
}
func (m *UsersURL) XXX_Size() int {
	return xxx_messageInfo_UsersURL.Size(m)
}
func (m *UsersURL) XXX_DiscardUnknown() {
	xxx_messageInfo_UsersURL.DiscardUnknown(m)
}

var xxx_messageInfo_UsersURL proto.InternalMessageInfo

func (m *UsersURL) GetUrls() []*UserURL {
	if m != nil {
		return m.Urls
	}
	return nil
}

// The item for UsersURL.
type UserURL struct {
	ShortUrl             string   `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	LongUrl              string   `protobuf:"bytes,2,opt,name=long_url,json=longUrl,proto3" json:"long_url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserURL) Reset()         { *m = UserURL{} }
func (m *UserURL) String() string { return proto.CompactTextString(m) }
func (*UserURL) ProtoMessage()    {}
func (*UserURL) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{6}
}

func (m *UserURL) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserURL.Unmarshal(m, b)
}
func (m *UserURL) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserURL.Marshal(b, m, deterministic)
}
func (m *UserURL) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserURL.Merge(m, src)
}
func (m *UserURL) XXX_Size() int {
	return xxx_messageInfo_UserURL.Size(m)
}
func (m *UserURL) XXX_DiscardUnknown() {
	xxx_messageInfo_UserURL.DiscardUnknown(m)
}

var xxx_messageInfo_UserURL proto.InternalMessageInfo

func (m *UserURL) GetShortUrl() string {
	if m != nil {
		return m.ShortUrl
	}
	return ""
}

func (m *UserURL) GetLongUrl() string {
	if m != nil {
		return m.LongUrl
	}
	return ""
}

// The request message for ShortenURLBatch.
type ShortenBatchRequest struct {
	Urls                 []*ShortenBatchRequestItem `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ShortenBatchRequest) Reset()         { *m = ShortenBatchRequest{} }
func (m *ShortenBatchRequest) String() string { return proto.CompactTextString(m) }
func (*ShortenBatchRequest) ProtoMessage()    {}
func (*ShortenBatchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{7}
}

func (m *ShortenBatchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShortenBatchRequest.Unmarshal(m, b)
}
func (m *ShortenBatchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShortenBatchRequest.Marshal(b, m, deterministic)
}
func (m *ShortenBatchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShortenBatchRequest.Merge(m, src)
}
func (m *ShortenBatchRequest) XXX_Size() int {
	return xxx_messageInfo_ShortenBatchRequest.Size(m)
}
func (m *ShortenBatchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShortenBatchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShortenBatchRequest proto.InternalMessageInfo

func (m *ShortenBatchRequest) GetUrls() []*ShortenBatchRequestItem {
	if m != nil {
		return m.Urls
	}
	return nil
}

// The item ShortenBatchReques.
type ShortenBatchRequestItem struct {
	CorrelationId        string   `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	OriginalUrl          string   `protobuf:"bytes,2,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShortenBatchRequestItem) Reset()         { *m = ShortenBatchRequestItem{} }
func (m *ShortenBatchRequestItem) String() string { return proto.CompactTextString(m) }
func (*ShortenBatchRequestItem) ProtoMessage()    {}
func (*ShortenBatchRequestItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{8}
}

func (m *ShortenBatchRequestItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShortenBatchRequestItem.Unmarshal(m, b)
}
func (m *ShortenBatchRequestItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShortenBatchRequestItem.Marshal(b, m, deterministic)
}
func (m *ShortenBatchRequestItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShortenBatchRequestItem.Merge(m, src)
}
func (m *ShortenBatchRequestItem) XXX_Size() int {
	return xxx_messageInfo_ShortenBatchRequestItem.Size(m)
}
func (m *ShortenBatchRequestItem) XXX_DiscardUnknown() {
	xxx_messageInfo_ShortenBatchRequestItem.DiscardUnknown(m)
}

var xxx_messageInfo_ShortenBatchRequestItem proto.InternalMessageInfo

func (m *ShortenBatchRequestItem) GetCorrelationId() string {
	if m != nil {
		return m.CorrelationId
	}
	return ""
}

func (m *ShortenBatchRequestItem) GetOriginalUrl() string {
	if m != nil {
		return m.OriginalUrl
	}
	return ""
}

// The response message for ShortenURLBatch.
type ShortenBatchResponse struct {
	Urls                 []*ShortenBatchResponseItem `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *ShortenBatchResponse) Reset()         { *m = ShortenBatchResponse{} }
func (m *ShortenBatchResponse) String() string { return proto.CompactTextString(m) }
func (*ShortenBatchResponse) ProtoMessage()    {}
func (*ShortenBatchResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{9}
}

func (m *ShortenBatchResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShortenBatchResponse.Unmarshal(m, b)
}
func (m *ShortenBatchResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShortenBatchResponse.Marshal(b, m, deterministic)
}
func (m *ShortenBatchResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShortenBatchResponse.Merge(m, src)
}
func (m *ShortenBatchResponse) XXX_Size() int {
	return xxx_messageInfo_ShortenBatchResponse.Size(m)
}
func (m *ShortenBatchResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ShortenBatchResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ShortenBatchResponse proto.InternalMessageInfo

func (m *ShortenBatchResponse) GetUrls() []*ShortenBatchResponseItem {
	if m != nil {
		return m.Urls
	}
	return nil
}

// The item for ShortenBatchResponse.
type ShortenBatchResponseItem struct {
	CorrelationId        string   `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	ShortUrl             string   `protobuf:"bytes,2,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	Error                string   `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShortenBatchResponseItem) Reset()         { *m = ShortenBatchResponseItem{} }
func (m *ShortenBatchResponseItem) String() string { return proto.CompactTextString(m) }
func (*ShortenBatchResponseItem) ProtoMessage()    {}
func (*ShortenBatchResponseItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{10}
}

func (m *ShortenBatchResponseItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShortenBatchResponseItem.Unmarshal(m, b)
}
func (m *ShortenBatchResponseItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShortenBatchResponseItem.Marshal(b, m, deterministic)
}
func (m *ShortenBatchResponseItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShortenBatchResponseItem.Merge(m, src)
}
func (m *ShortenBatchResponseItem) XXX_Size() int {
	return xxx_messageInfo_ShortenBatchResponseItem.Size(m)
}
func (m *ShortenBatchResponseItem) XXX_DiscardUnknown() {
	xxx_messageInfo_ShortenBatchResponseItem.DiscardUnknown(m)
}

var xxx_messageInfo_ShortenBatchResponseItem proto.InternalMessageInfo

func (m *ShortenBatchResponseItem) GetCorrelationId() string {
	if m != nil {
		return m.CorrelationId
	}
	return ""
}

func (m *ShortenBatchResponseItem) GetShortUrl() string {
	if m != nil {
		return m.ShortUrl
	}
	return ""
}

func (m *ShortenBatchResponseItem) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

// The request message for ShortenDelete.
type ShortenDeleteRequest struct {
	Urls                 []string `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShortenDeleteRequest) Reset()         { *m = ShortenDeleteRequest{} }
func (m *ShortenDeleteRequest) String() string { return proto.CompactTextString(m) }
func (*ShortenDeleteRequest) ProtoMessage()    {}
func (*ShortenDeleteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{11}
}

func (m *ShortenDeleteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShortenDeleteRequest.Unmarshal(m, b)
}
func (m *ShortenDeleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShortenDeleteRequest.Marshal(b, m, deterministic)
}
func (m *ShortenDeleteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShortenDeleteRequest.Merge(m, src)
}
func (m *ShortenDeleteRequest) XXX_Size() int {
	return xxx_messageInfo_ShortenDeleteRequest.Size(m)
}
func (m *ShortenDeleteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShortenDeleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShortenDeleteRequest proto.InternalMessageInfo

func (m *ShortenDeleteRequest) GetUrls() []string {
	if m != nil {
		return m.Urls
	}
	return nil
}

// The response message for StatsInternalReponse.
type StatsInternalReponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	Urls                 int64    `protobuf:"varint,1,opt,name=urls,proto3" json:"urls,omitempty"`
	Users                int64    `protobuf:"varint,2,opt,name=users,proto3" json:"users,omitempty"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatsInternalReponse) Reset()         { *m = StatsInternalReponse{} }
func (m *StatsInternalReponse) String() string { return proto.CompactTextString(m) }
func (*StatsInternalReponse) ProtoMessage()    {}
func (*StatsInternalReponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{12}
}

func (m *StatsInternalReponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatsInternalReponse.Unmarshal(m, b)
}
func (m *StatsInternalReponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatsInternalReponse.Marshal(b, m, deterministic)
}
func (m *StatsInternalReponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatsInternalReponse.Merge(m, src)
}
func (m *StatsInternalReponse) XXX_Size() int {
	return xxx_messageInfo_StatsInternalReponse.Size(m)
}
func (m *StatsInternalReponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatsInternalReponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatsInternalReponse proto.InternalMessageInfo

func (m *StatsInternalReponse) GetUrls() int64 {
	if m != nil {
		return m.Urls
	}
	return 0
}

func (m *StatsInternalReponse) GetUsers() int64 {
	if m != nil {
		return m.Users
	}
	return 0
}

func init() {
	proto.RegisterType((*Empty)(nil), "github.com.alaleks.shortener.Empty")
	proto.RegisterType((*ShortenRequest)(nil), "github.com.alaleks.shortener.ShortenRequest")
	proto.RegisterType((*ShortenResponse)(nil), "github.com.alaleks.shortener.ShortenResponse")
	proto.RegisterType((*StatRequest)(nil), "github.com.alaleks.shortener.StatRequest")
	proto.RegisterType((*StatResponse)(nil), "github.com.alaleks.shortener.StatResponse")
	proto.RegisterType((*UsersURL)(nil), "github.com.alaleks.shortener.UsersURL")
	proto.RegisterType((*UserURL)(nil), "github.com.alaleks.shortener.UserURL")
	proto.RegisterType((*ShortenBatchRequest)(nil), "github.com.alaleks.shortener.ShortenBatchRequest")
	proto.RegisterType((*ShortenBatchRequestItem)(nil), "github.com.alaleks.shortener.ShortenBatchRequestItem")
	proto.RegisterType((*ShortenBatchResponse)(nil), "github.com.alaleks.shortener.ShortenBatchResponse")
	proto.RegisterType((*ShortenBatchResponseItem)(nil), "github.com.alaleks.shortener.ShortenBatchResponseItem")
	proto.RegisterType((*ShortenDeleteRequest)(nil), "github.com.alaleks.shortener.ShortenDeleteRequest")
	proto.RegisterType((*StatsInternalReponse)(nil), "github.com.alaleks.shortener.StatsInternalReponse")
}

func init() {
	proto.RegisterFile("shortener.proto", fileDescriptor_6a64040fb43d257f)
}

var fileDescriptor_6a64040fb43d257f = []byte{
	// 591 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xdf, 0x6f, 0xd3, 0x40,
	0x0c, 0x6e, 0x96, 0x6d, 0x6d, 0xdc, 0xfd, 0x40, 0xc7, 0x04, 0x21, 0xec, 0xa1, 0x1c, 0x0c, 0x75,
	0x13, 0x64, 0xa2, 0x08, 0x24, 0xde, 0xd8, 0xc4, 0x34, 0x15, 0xed, 0x29, 0x53, 0x1e, 0x40, 0x48,
	0x25, 0x4d, 0xad, 0xb6, 0xec, 0x9a, 0x94, 0xbb, 0x8b, 0x04, 0xfc, 0x11, 0xfc, 0xcd, 0x28, 0xc9,
	0xe5, 0x9a, 0x54, 0xac, 0x4d, 0x9f, 0x5a, 0xbb, 0x9f, 0xed, 0xcf, 0xf6, 0xe7, 0x2b, 0x1c, 0x8a,
	0x49, 0xcc, 0x25, 0x46, 0xc8, 0xdd, 0x39, 0x8f, 0x65, 0x4c, 0x8e, 0xc7, 0x53, 0x39, 0x49, 0x86,
	0x6e, 0x18, 0xcf, 0xdc, 0x80, 0x05, 0x0c, 0xef, 0x84, 0xab, 0x31, 0xb4, 0x09, 0x3b, 0x57, 0xb3,
	0xb9, 0xfc, 0x4d, 0x29, 0x1c, 0xdc, 0xe6, 0x5e, 0x0f, 0x7f, 0x26, 0x28, 0x24, 0x79, 0x00, 0x66,
	0xc2, 0x99, 0x6d, 0x74, 0x8c, 0xae, 0xe5, 0xa5, 0x5f, 0xe9, 0x17, 0x38, 0xd4, 0x18, 0x31, 0x8f,
	0x23, 0x81, 0xe4, 0x08, 0x76, 0x90, 0xf3, 0x98, 0x2b, 0x58, 0x6e, 0x90, 0x47, 0xb0, 0xcb, 0x51,
	0x24, 0x4c, 0xda, 0x5b, 0x99, 0x5b, 0x59, 0xc4, 0x86, 0xa6, 0x48, 0xc2, 0x10, 0x85, 0xb0, 0xcd,
	0x8e, 0xd1, 0x6d, 0x79, 0x85, 0x49, 0x4f, 0xa1, 0x7d, 0x2b, 0x03, 0x59, 0xd4, 0x76, 0xa0, 0x95,
	0x71, 0x4c, 0xa6, 0x23, 0x95, 0x59, 0xdb, 0xf4, 0x17, 0xec, 0xe5, 0x50, 0x45, 0x41, 0x63, 0x35,
	0x59, 0x6d, 0xa7, 0x05, 0x59, 0x1c, 0x8d, 0xd3, 0x9f, 0x72, 0x26, 0x85, 0x49, 0x8e, 0xc1, 0x0a,
	0x39, 0x06, 0x12, 0x47, 0x17, 0x32, 0x23, 0x63, 0x79, 0x0b, 0x47, 0xda, 0x56, 0x22, 0x82, 0x31,
	0xda, 0xdb, 0x1d, 0xa3, 0xbb, 0xed, 0xe5, 0x06, 0xbd, 0x82, 0x96, 0x2f, 0x90, 0x0b, 0xdf, 0xbb,
	0x21, 0x1f, 0x60, 0x3b, 0xe1, 0x4c, 0xd8, 0x46, 0xc7, 0xec, 0xb6, 0x7b, 0x27, 0xee, 0xaa, 0x29,
	0xbb, 0x69, 0x94, 0xef, 0xdd, 0x78, 0x59, 0x08, 0xbd, 0x80, 0xa6, 0x72, 0x90, 0xa7, 0x60, 0x65,
	0xa8, 0xc1, 0x32, 0x79, 0x9f, 0x33, 0xf2, 0x04, 0x5a, 0x29, 0xdb, 0xc1, 0x12, 0x7b, 0x9f, 0x33,
	0xfa, 0x1d, 0x1e, 0xaa, 0x4d, 0x5c, 0x06, 0x32, 0x9c, 0x14, 0x63, 0xeb, 0x57, 0x48, 0xbd, 0x5b,
	0x4d, 0xea, 0x3f, 0x09, 0xfa, 0x12, 0x67, 0x8a, 0x64, 0x08, 0x8f, 0xef, 0x01, 0x90, 0x13, 0x38,
	0x08, 0x63, 0xce, 0x91, 0x05, 0x72, 0x1a, 0x47, 0x03, 0xbd, 0xa2, 0xfd, 0x92, 0xb7, 0x3f, 0x22,
	0xcf, 0x60, 0x2f, 0xe6, 0xd3, 0xf1, 0x34, 0x0a, 0x58, 0xa9, 0x85, 0x76, 0xe1, 0x4b, 0xdb, 0x18,
	0xc2, 0x51, 0xb5, 0x88, 0x5a, 0xe9, 0xe7, 0x4a, 0x1f, 0xef, 0x37, 0xe9, 0x23, 0xcf, 0x50, 0x6a,
	0x44, 0x82, 0x7d, 0x1f, 0xa2, 0x6e, 0x27, 0x95, 0x2d, 0x6d, 0x2d, 0x6d, 0x49, 0x5f, 0x80, 0x59,
	0xba, 0x00, 0x7a, 0xa6, 0x3b, 0xfb, 0x84, 0x0c, 0x25, 0x16, 0x1b, 0x22, 0xa5, 0xce, 0x2c, 0xc5,
	0xf0, 0x23, 0x1c, 0xa5, 0x82, 0x16, 0xfd, 0x48, 0x22, 0x8f, 0x02, 0xe6, 0x61, 0x3e, 0x85, 0x05,
	0xd6, 0xe8, 0x9a, 0x39, 0x36, 0x17, 0x26, 0x72, 0x91, 0xd1, 0x30, 0xbd, 0xdc, 0xe8, 0xfd, 0xdd,
	0x01, 0xeb, 0xb6, 0x18, 0x08, 0xb9, 0x03, 0x50, 0x46, 0x2a, 0xb1, 0x57, 0xb5, 0xa6, 0xa7, 0xf8,
	0x39, 0xaf, 0x6b, 0xa2, 0xf3, 0x21, 0xd2, 0x06, 0x19, 0x42, 0xf3, 0x1a, 0x65, 0xca, 0x9f, 0x9c,
	0xae, 0x89, 0x5d, 0xdc, 0xb7, 0x73, 0x56, 0x07, 0xaa, 0x6b, 0x7c, 0x83, 0xf6, 0x35, 0x4a, 0x7d,
	0x7a, 0xcf, 0x57, 0x07, 0x67, 0xef, 0x99, 0xf3, 0x72, 0xfd, 0x45, 0xa6, 0xc9, 0x68, 0x83, 0xfc,
	0xd1, 0xaf, 0x9a, 0xef, 0xdd, 0x64, 0x1a, 0x21, 0x6f, 0x36, 0xbe, 0x1c, 0xa7, 0xb7, 0xb9, 0x48,
	0x69, 0x83, 0xfc, 0x80, 0xfd, 0x8a, 0x4c, 0x48, 0xbd, 0x34, 0x15, 0x4d, 0x39, 0x75, 0xe6, 0xa1,
	0x6a, 0x95, 0x65, 0x56, 0x6f, 0x8e, 0xbd, 0xf5, 0x9b, 0x5a, 0x16, 0x2e, 0x6d, 0x5c, 0xbe, 0xf8,
	0x4a, 0x17, 0x61, 0xe7, 0x2a, 0xec, 0x5c, 0x87, 0x9d, 0x67, 0x7f, 0x4d, 0xc3, 0xdd, 0xec, 0xe3,
	0xed, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x84, 0xff, 0x6a, 0xf5, 0xb4, 0x06, 0x00, 0x00,
}
