package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"

	"github.com/go-kit/kit/endpoint"
)

type Middleware func(endpoint.Endpoint) endpoint.Endpoint

type loggingMiddleware struct {
	logger log.Logger
	next   StringService
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           StringService
}

func (mw loggingMiddleware) Uppercase(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.next.Uppercase(ctx, s)
	return
}

func (mw loggingMiddleware) Count(ctx context.Context, s string) (n int) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())
	n = mw.next.Count(ctx, s)
	return
}

func (mw instrumentingMiddleware) Uppercase(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{
			"method", "uppercase",
			"error", fmt.Sprint(err != nil),
		}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.next.Uppercase(ctx, s)
	return
}

func (mw instrumentingMiddleware) Count(ctx context.Context, s string) (n int) {
	defer func(begin time.Time) {
		lvs := []string{
			"method", "count",
			"error", "false",
		}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.countResult.Observe(float64(n))
	}(time.Now())
	n = mw.next.Count(ctx, s)
	return
}
