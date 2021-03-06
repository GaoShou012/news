// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/news/news.proto

package news

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for NewsService service

func NewNewsServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for NewsService service

type NewsService interface {
	Sub(ctx context.Context, in *SubReq, opts ...client.CallOption) (*SubRsp, error)
	Cancel(ctx context.Context, in *CancelReq, opts ...client.CallOption) (*CancelRsp, error)
	GetSubList(ctx context.Context, in *GetSubListReq, opts ...client.CallOption) (*GetSubListRsp, error)
	GetNews(ctx context.Context, in *GetNewsReq, opts ...client.CallOption) (*GetNewsRsp, error)
}

type newsService struct {
	c    client.Client
	name string
}

func NewNewsService(name string, c client.Client) NewsService {
	return &newsService{
		c:    c,
		name: name,
	}
}

func (c *newsService) Sub(ctx context.Context, in *SubReq, opts ...client.CallOption) (*SubRsp, error) {
	req := c.c.NewRequest(c.name, "NewsService.Sub", in)
	out := new(SubRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsService) Cancel(ctx context.Context, in *CancelReq, opts ...client.CallOption) (*CancelRsp, error) {
	req := c.c.NewRequest(c.name, "NewsService.Cancel", in)
	out := new(CancelRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsService) GetSubList(ctx context.Context, in *GetSubListReq, opts ...client.CallOption) (*GetSubListRsp, error) {
	req := c.c.NewRequest(c.name, "NewsService.GetSubList", in)
	out := new(GetSubListRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsService) GetNews(ctx context.Context, in *GetNewsReq, opts ...client.CallOption) (*GetNewsRsp, error) {
	req := c.c.NewRequest(c.name, "NewsService.GetNews", in)
	out := new(GetNewsRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for NewsService service

type NewsServiceHandler interface {
	Sub(context.Context, *SubReq, *SubRsp) error
	Cancel(context.Context, *CancelReq, *CancelRsp) error
	GetSubList(context.Context, *GetSubListReq, *GetSubListRsp) error
	GetNews(context.Context, *GetNewsReq, *GetNewsRsp) error
}

func RegisterNewsServiceHandler(s server.Server, hdlr NewsServiceHandler, opts ...server.HandlerOption) error {
	type newsService interface {
		Sub(ctx context.Context, in *SubReq, out *SubRsp) error
		Cancel(ctx context.Context, in *CancelReq, out *CancelRsp) error
		GetSubList(ctx context.Context, in *GetSubListReq, out *GetSubListRsp) error
		GetNews(ctx context.Context, in *GetNewsReq, out *GetNewsRsp) error
	}
	type NewsService struct {
		newsService
	}
	h := &newsServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&NewsService{h}, opts...))
}

type newsServiceHandler struct {
	NewsServiceHandler
}

func (h *newsServiceHandler) Sub(ctx context.Context, in *SubReq, out *SubRsp) error {
	return h.NewsServiceHandler.Sub(ctx, in, out)
}

func (h *newsServiceHandler) Cancel(ctx context.Context, in *CancelReq, out *CancelRsp) error {
	return h.NewsServiceHandler.Cancel(ctx, in, out)
}

func (h *newsServiceHandler) GetSubList(ctx context.Context, in *GetSubListReq, out *GetSubListRsp) error {
	return h.NewsServiceHandler.GetSubList(ctx, in, out)
}

func (h *newsServiceHandler) GetNews(ctx context.Context, in *GetNewsReq, out *GetNewsRsp) error {
	return h.NewsServiceHandler.GetNews(ctx, in, out)
}
