package endpoint

import (
	"context"
	"time"

	"golang.org/x/time/rate"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"

	"git.heytea.com/go-kit-example/user-service/model"
	"git.heytea.com/go-kit-example/user-service/pkg/service"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	ListUsersEndpoint endpoint.Endpoint
	GetUserEndpoint   endpoint.Endpoint
	NewUserEndpoint   endpoint.Endpoint
	DelUserEndpoint   endpoint.Endpoint
	PutUserEndpoint   endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc service.UserService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) Set {
	var listUsersEndpoint endpoint.Endpoint
	{
		listUsersEndpoint = MakeListUsersEndpoint(svc)
		listUsersEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(listUsersEndpoint)
		listUsersEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(listUsersEndpoint)
		listUsersEndpoint = opentracing.TraceServer(otTracer, "ListUsers")(listUsersEndpoint)
		listUsersEndpoint = zipkin.TraceEndpoint(zipkinTracer, "ListUsers")(listUsersEndpoint)
		listUsersEndpoint = LoggingMiddleware(log.With(logger, "method", "ListUsers"))(listUsersEndpoint)
		listUsersEndpoint = InstrumentingMiddleware(duration.With("method", "ListUsers"))(listUsersEndpoint)
	}

	return Set{
		ListUsersEndpoint: listUsersEndpoint,
		GetUserEndpoint:   MakeGetUserEndpoint(svc),
		NewUserEndpoint:   MakeNewlUserEndpoint(svc),
		DelUserEndpoint:   MakeDelUserEndpoint(svc),
		PutUserEndpoint:   MakePutUserEndpoint(svc),
	}
}

// GetUser(ctx context.Context, id int64) (*User, error)
// 	NewUser(ctx context.Context, user *User) (*User, error)
// 	DelUser(ctx context.Context, id int64) (*User, error)
// 	PutUser(ctx context.Context, user *User) (*User, error)

// This is primarily useful in the context of a client library.
func (s Set) ListUsers(ctx context.Context, query string, pageNumber int64, pageSize int64) ([]*model.User, error) {
	resp, err := s.ListUsersEndpoint(ctx, ListUsersRequest{
		Query:      query,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(ListUsersResponse)
	return response.Data, response.Err
}

func (s Set) GetUser(ctx context.Context, id int64) (*model.User, error) {
	resp, err := s.GetUserEndpoint(ctx, GetUserReqeust{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(GetUserResponse)
	return response.Data, response.Err
}

func (s Set) NewUser(ctx context.Context, user *model.User) (*model.User, error) {
	resp, err := s.NewUserEndpoint(ctx, NewUserRequest{
		User: user,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(NewUserResponse)
	return nil, response.Err
}

func (s Set) DelUser(ctx context.Context, id int64) (*model.User, error) {
	resp, err := s.DelUserEndpoint(ctx, DelUserRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(DelUserResponse)
	return nil, response.Err
}

func (s Set) PutUser(ctx context.Context, user *model.User) (*model.User, error) {
	resp, err := s.PutUserEndpoint(ctx, PutUserRequest{
		User: user,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(PutUserResponse)
	return nil, response.Err
}

func MakeListUsersEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListUsersRequest)
		v, err := s.ListUsers(ctx, req.Query, req.PageNumber, req.PageSize)
		return ListUsersResponse{
			Code:    200,
			Message: "list user success",
			Data:    v,
		}, nil
	}
}

func MakeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserReqeust)
		v, err := s.GetUser(ctx, req.Id)
		return GetUserResponse{
			Code:    200,
			Message: "get user success",
			Data:    v,
		}, err
	}
}

func MakeNewlUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(NewUserRequest)
		_, err = s.NewUser(ctx, req.User)
		return NewUserResponse{
			Code:    200,
			Message: "create user success",
		}, err
	}
}

func MakeDelUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DelUserRequest)
		_, err = s.DelUser(ctx, req.Id)
		return DelUserResponse{
			Code:    200,
			Message: "delete user success",
		}, err
	}
}

func MakePutUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PutUserRequest)
		_, err = s.PutUser(ctx, req.User)
		return PutUserResponse{
			Code:    200,
			Message: "update user success",
		}, err
	}
}

type ListUsersRequest struct {
	Query      string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	PageNumber int64  `protobuf:"varint,2,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	PageSize   int64  `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

type GetUserReqeust struct {
	Id int64 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

type NewUserRequest struct {
	User *model.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

type DelUserRequest struct {
	Id int64 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}
type PutUserRequest struct {
	User *model.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

type ListUsersResponse struct {
	Code    int64         `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string        `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    []*model.User `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Err     error         `protobuf:"bytes,4,opt,name=err,proto3" json:"err,omitempty"`
}

type GetUserResponse struct {
	Code    int64       `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string      `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *model.User `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Err     error       `protobuf:"bytes,4,opt,name=err,proto3" json:"err,omitempty"`
}

type NewUserResponse struct {
	Code    int64  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Err     error  `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
}

type DelUserResponse struct {
	Code    int64  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Err     error  `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
}

type PutUserResponse struct {
	Code    int64  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Err     error  `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
}

func (r ListUsersResponse) Failed() error { return r.Err }
func (r GetUserResponse) Failed() error   { return r.Err }
func (r NewUserResponse) Failed() error   { return r.Err }
func (r DelUserResponse) Failed() error   { return r.Err }
func (r PutUserResponse) Failed() error   { return r.Err }

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = ListUsersResponse{}
	_ endpoint.Failer = GetUserResponse{}
	_ endpoint.Failer = NewUserResponse{}
	_ endpoint.Failer = DelUserResponse{}
	_ endpoint.Failer = PutUserResponse{}
)
