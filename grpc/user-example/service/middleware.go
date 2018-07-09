package service

import (
	"context"

	"git.heytea.com/go-kit-example/user-service/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(UserService) UserService

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next UserService) UserService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   UserService
}

func (mw loggingMiddleware) ListUsers(ctx context.Context, query string, pageNumber int64, pageSize int64) (resp []*model.User, err error) {
	defer func() {
		mw.logger.Log("method", "ListUsers", "query", query, "page", pageNumber, "size", pageSize, "err", err)
	}()
	return mw.next.ListUsers(ctx, query, pageNumber, pageSize)
}

func (mw loggingMiddleware) GetUser(ctx context.Context, id int64) (resp *model.User, err error) {
	defer func() {
		mw.logger.Log("method", "GetUser", "id", id, "err", err)
	}()
	return mw.next.GetUser(ctx, id)
}

func (mw loggingMiddleware) NewUser(ctx context.Context, user *model.User) (resp *model.User, err error) {
	defer func() {
		mw.logger.Log("method", "NewUser", "username", user.Name, "err", err)
	}()
	return mw.next.NewUser(ctx, user)
}

func (mw loggingMiddleware) DelUser(ctx context.Context, id int64) (resp *model.User, err error) {
	defer func() {
		mw.logger.Log("method", "DelUser", "id", id, "err", err)
	}()
	return mw.next.DelUser(ctx, id)
}

func (mw loggingMiddleware) PutUser(ctx context.Context, user *model.User) (resp *model.User, err error) {
	defer func() {
		mw.logger.Log("method", "UpdateUser", "username", user.Name, "err", err)
	}()
	return mw.next.PutUser(ctx, user)
}

// InstrumentingMiddleware returns a service middleware that instruments
// the number of integers summed and characters concatenated over the lifetime of
// the service.
func InstrumentingMiddleware(ints, chars metrics.Counter) Middleware {
	return func(next UserService) UserService {
		return instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
	}
}

type instrumentingMiddleware struct {
	ints            metrics.Counter
	chars           metrics.Counter
	countListUsers  metrics.Counter
	countGetUser    metrics.Counter
	countNewUser    metrics.Counter
	countDelUser    metrics.Counter
	countUpdateUser metrics.Counter
	next            UserService
}

func (mw instrumentingMiddleware) ListUsers(ctx context.Context, query string, pageNumber int64, pageSize int64) ([]*model.User, error) {
	v, err := mw.next.ListUsers(ctx, query, pageNumber, pageSize)
	mw.chars.Add(1)
	return v, err
}

func (mw instrumentingMiddleware) GetUser(ctx context.Context, id int64) (*model.User, error) {
	v, err := mw.next.GetUser(ctx, id)
	mw.chars.Add(1)
	return v, err
}

func (mw instrumentingMiddleware) NewUser(ctx context.Context, user *model.User) (*model.User, error) {
	v, err := mw.next.NewUser(ctx, user)
	mw.chars.Add(1)
	return v, err
}

func (mw instrumentingMiddleware) DelUser(ctx context.Context, id int64) (*model.User, error) {
	v, err := mw.next.DelUser(ctx, id)
	mw.chars.Add(1)
	return v, err
}

func (mw instrumentingMiddleware) PutUser(ctx context.Context, user *model.User) (*model.User, error) {
	v, err := mw.next.PutUser(ctx, user)
	mw.chars.Add(1)
	return v, err
}
