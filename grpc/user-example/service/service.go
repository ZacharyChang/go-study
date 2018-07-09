package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"

	db "git.heytea.com/go-kit-example/user-service/db"
	"git.heytea.com/go-kit-example/user-service/model"
)

// service接口定义
type UserService interface {
	ListUsers(ctx context.Context, query string, pageNumber int64, pageSize int64) ([]*model.User, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	NewUser(ctx context.Context, user *model.User) (*model.User, error)
	DelUser(ctx context.Context, id int64) (*model.User, error)
	PutUser(ctx context.Context, user *model.User) (*model.User, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
// 在此添加中间件
func NewUserService(logger log.Logger, ints, chars metrics.Counter) UserService {
	var svc UserService
	{
		svc = NewBasicService()
		svc = LoggingMiddleware(logger)(svc)
		svc = InstrumentingMiddleware(ints, chars)(svc)
	}
	return svc
}

var (
	// ErrTwoZeroes is an arbitrary business rule for the Add method.
	ErrTwoZeroes = errors.New("can't sum two zeroes")

	// ErrIntOverflow protects the Add method. We've decided that this error
	// indicates a misbehaving service and should count against e.g. circuit
	// breakers. So, we return it directly in endpoints, to illustrate the
	// difference. In a real service, this probably wouldn't be the case.
	ErrIntOverflow = errors.New("integer overflow")

	// ErrMaxSizeExceeded protects the Concat method.
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")

	ErrNoSuchElement = errors.New("no such element")
)

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() UserService {
	return basicService{}
}

type basicService struct{}

// 实现业务接口
func (s basicService) ListUsers(ctx context.Context, query string, pageNumber int64, pageSize int64) (resp []*model.User, err error) {
	fmt.Println("ListUsers called")
	fmt.Println(db.DefaultConn)
	var users []*model.User
	db.DefaultConn.Find(&users)
	return users, nil
}

func (s basicService) GetUser(ctx context.Context, id int64) (resp *model.User, err error) {
	fmt.Println("GetUser called")
	return &model.User{}, nil
}

func (s basicService) NewUser(ctx context.Context, user *model.User) (resp *model.User, err error) {
	fmt.Println("NewUser called")
	return &model.User{}, nil
}

func (s basicService) DelUser(ctx context.Context, id int64) (resp *model.User, err error) {
	fmt.Println("DelUser called")
	return &model.User{}, nil
}

func (s basicService) PutUser(context.Context, *model.User) (resp *model.User, err error) {
	fmt.Println("PutUser called")
	return &model.User{}, nil
}
