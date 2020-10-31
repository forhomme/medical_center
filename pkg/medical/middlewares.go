package medical

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type Middleware func(service Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostVisit(ctx context.Context, v Visit) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostVisit", "visitID", v.ID, "took", time.Since(begin), "error", err)
	}(time.Now())
	return mw.next.PostVisit(ctx, v)
}

func (mw loggingMiddleware) GetVisit(ctx context.Context, id string) (v Visit, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetVisit", "visitID", id, "took", time.Since(begin), "error", err)
	}(time.Now())
	return mw.next.GetVisit(ctx, id)
}

func (mw loggingMiddleware) PostPatient(ctx context.Context, p Patient) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostPatient", "patientID", p.ID, "took", time.Since(begin), "error", err)
	}(time.Now())
	return mw.next.PostPatient(ctx, p)
}

func (mw loggingMiddleware) GetPatient(ctx context.Context, patientID string) (p Patient, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetPatient", "patientID", patientID, "took", time.Since(begin), "error", err)
	}(time.Now())
	return mw.next.GetPatient(ctx, patientID)
}
