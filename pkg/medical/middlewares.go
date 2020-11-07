package medical

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"time"
)

type Middleware func(service Service) Service

type Log struct {
	ID       int           `json:"id" gorm:"primaryKey"`
	Method   string        `json:"method"`
	CalledID int           `json:"called_id"`
	Created  time.Time     `json:"created"`
	Took     time.Duration `json:"took"`
	Error    error         `json:"error"`
}

func LoggingMiddleware(logger log.Logger, db *gorm.DB) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
			db:     db,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
	db     *gorm.DB
}

func (mw loggingMiddleware) PostVisit(ctx context.Context, v Visit) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostVisit", "visitID", v.ID, "took", time.Since(begin), "error", err)
		logDetail := Log{
			Method:   "PostVisit",
			CalledID: v.ID,
			Created:  begin,
			Took:     time.Since(begin),
			Error:    err,
		}
		mw.db.Create(&logDetail)
	}(time.Now())
	return mw.next.PostVisit(ctx, v)
}

func (mw loggingMiddleware) GetVisit(ctx context.Context, id string) (v Visit, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetVisit", "visitID", id, "took", time.Since(begin), "error", err)
		logDetail := Log{
			Method:   "GetVisit",
			CalledID: v.ID,
			Created:  begin,
			Took:     time.Since(begin),
			Error:    err,
		}
		mw.db.Create(&logDetail)
	}(time.Now())
	return mw.next.GetVisit(ctx, id)
}

func (mw loggingMiddleware) PostPatient(ctx context.Context, p Patient) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostPatient", "patientID", p.ID, "took", time.Since(begin), "error", err)
		logDetail := Log{
			Method:   "PostPatient",
			CalledID: p.ID,
			Created:  begin,
			Took:     time.Since(begin),
			Error:    err,
		}
		mw.db.Create(&logDetail)
	}(time.Now())
	return mw.next.PostPatient(ctx, p)
}

func (mw loggingMiddleware) GetPatient(ctx context.Context, patientID string) (p Patient, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetPatient", "patientID", patientID, "took", time.Since(begin), "error", err)
		logDetail := Log{
			Method:   "GetPatient",
			CalledID: p.ID,
			Created:  begin,
			Took:     time.Since(begin),
			Error:    err,
		}
		mw.db.Create(&logDetail)
	}(time.Now())
	return mw.next.GetPatient(ctx, patientID)
}
