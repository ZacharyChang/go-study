package transport

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	oldcontext "golang.org/x/net/context"
	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"git.heytea.com/go-kit-example/user-service/model"
	pb "git.heytea.com/go-kit-example/user-service/pb"
	"git.heytea.com/go-kit-example/user-service/pkg/endpoint"
	"git.heytea.com/go-kit-example/user-service/pkg/service"
)

type grpcServer struct {
	listUsers grpctransport.Handler
	getUser   grpctransport.Handler
	newUser   grpctransport.Handler
	delUser   grpctransport.Handler
	putUser   grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints endpoint.Set, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) pb.UserServiceServer {
	// Zipkin GRPC Server Trace can either be instantiated per gRPC method with a
	// provided operation name or a global tracing service can be instantiated
	// without an operation name and fed to each Go kit gRPC server as a
	// ServerOption.
	// In the latter case, the operation name will be the endpoint's grpc method
	// path if used in combination with the Go kit gRPC Interceptor.
	//
	// In this example, we demonstrate a global Zipkin tracing service with
	// Go kit gRPC Interceptor.
	zipkinServer := zipkin.GRPCServerTrace(zipkinTracer)

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
		zipkinServer,
	}

	return &grpcServer{
		listUsers: grpctransport.NewServer(
			endpoints.ListUsersEndpoint,
			decodeGRPCListUsersRequest,
			encodeGRPCListUsersResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "listUsers", logger)))...,
		),
	}
}

func (s *grpcServer) ListUsers(ctx oldcontext.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	_, rep, err := s.listUsers.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListUsersResponse), nil
}

func (s *grpcServer) GetUser(ctx context.Context, req *pb.GetUserReqeust) (*pb.GetUserResponse, error) {
	_, rep, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetUserResponse), nil
}

func (s *grpcServer) NewUser(ctx context.Context, req *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	_, rep, err := s.newUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.NewUserResponse), nil
}

func (s *grpcServer) DelUser(ctx context.Context, req *pb.DelUserRequest) (*pb.DelUserResponse, error) {
	_, rep, err := s.delUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.DelUserResponse), nil
}

func (s *grpcServer) PutUser(ctx context.Context, req *pb.PutUserRequest) (*pb.PutUserResponse, error) {
	_, rep, err := s.putUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.PutUserResponse), nil
}

// NewGRPCClient returns an AddService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.UserService {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// Zipkin GRPC Client Trace can either be instantiated per gRPC method with a
	// provided operation name or a global tracing client can be instantiated
	// without an operation name and fed to each Go kit client as ClientOption.
	// In the latter case, the operation name will be the endpoint's grpc method
	// path.
	//
	// In this example, we demonstrace a global tracing client.
	zipkinClient := zipkin.GRPCClientTrace(zipkinTracer)

	// global client middlewares
	options := []grpctransport.ClientOption{
		zipkinClient,
	}

	// Each individual endpoint is an grpc/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var listUsersEndpoint kitendpoint.Endpoint
	{
		listUsersEndpoint = grpctransport.NewClient(
			conn,
			"pb.UserService",
			"ListUsers",
			encodeGRPCListUsersRequest,
			decodeGRPCListUsersResponse,
			pb.ListUsersResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		listUsersEndpoint = opentracing.TraceClient(otTracer, "ListUsers")(listUsersEndpoint)
		listUsersEndpoint = limiter(listUsersEndpoint)
		listUsersEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "ListUsers",
			Timeout: 30 * time.Second,
		}))(listUsersEndpoint)
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return endpoint.Set{
		ListUsersEndpoint: listUsersEndpoint,
	}
}

func decodeGRPCListUsersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListUsersRequest)
	return endpoint.ListUsersRequest{
		Query:      req.Query,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	}, nil
}

func decodeGRPCListUsersResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	resp := grpcReply.(*pb.ListUsersResponse)
	return endpoint.ListUsersResponse{
		Code:    resp.Code,
		Message: resp.Message,
		Data:    pbUserArray2Model(resp.Data),
		Err:     str2err(resp.Err),
	}, nil
}

func encodeGRPCListUsersRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoint.ListUsersRequest)
	return &pb.ListUsersRequest{
		Query:      req.Query,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	}, nil
}

func encodeGRPCListUsersResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.ListUsersResponse)
	return &pb.ListUsersResponse{
		Code:    resp.Code,
		Message: resp.Message,
		Data:    serviceUserArray2PB(resp.Data),
		Err:     err2str(resp.Err),
	}, nil
}

// These annoying helper functions are required to translate Go error types to
// and from strings, which is the type we use in our IDLs to represent errors.
// There is special casing to treat empty strings as nil errors.

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func pbUser2Service(u *pb.User) *model.User {
	return &model.User{
		Id:       u.Id,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
		Birthday: u.Birthday,
		IsVip:    u.IsVip,
	}
}

func pbUserArray2Model(users []*pb.User) []*model.User {
	result := make([]*model.User, len(users))
	for i, u := range users {
		result[i] = &model.User{
			Id:       u.Id,
			Name:     u.Name,
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
			Birthday: u.Birthday,
			IsVip:    u.IsVip,
		}
	}
	return result
}

func serviceUser2PB(u *model.User) *pb.User {
	return &pb.User{
		Id:       u.Id,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
		Birthday: u.Birthday,
		IsVip:    u.IsVip,
	}
}

func serviceUserArray2PB(users []*model.User) []*pb.User {
	result := make([]*pb.User, len(users))
	for i, u := range users {
		result[i] = &pb.User{
			Id:       u.Id,
			Name:     u.Name,
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
			Birthday: u.Birthday,
			IsVip:    u.IsVip,
		}
	}
	return result
}
