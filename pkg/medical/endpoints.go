package medical

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	PostVisitEndpoints   endpoint.Endpoint
	GetVisitEndpoints    endpoint.Endpoint
	PostPatientEndpoints endpoint.Endpoint
	GetPatientEndpoints  endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostVisitEndpoints:   MakePostVisitEndpoints(s),
		GetVisitEndpoints:    MakeGetVisitEndpoints(s),
		PostPatientEndpoints: MakePostPatientEndpoints(s),
		GetPatientEndpoints:  MakeGetPatientEndpoints(s),
	}
}

func MakePostVisitEndpoints(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postVisit)
		e := s.PostVisit(ctx, req.Visit)
		return postVisitResponse{Err: e}, nil
	}
}

func MakeGetVisitEndpoints(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getVisit)
		v, e := s.GetVisit(ctx, req.ID)
		return getVisitResponse{Visit: v, Err: e}, nil
	}
}

func MakePostPatientEndpoints(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postPatient)
		e := s.PostPatient(ctx, req.Patient)
		return postPatientResponse{Err: e}, nil
	}
}

func MakeGetPatientEndpoints(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPatient)
		p, e := s.GetPatient(ctx, req.ID)
		return getPatientResponse{Patient: p, Err: e}, nil
	}
}

// make struct for return of endpoints
type postVisit struct {
	Visit Visit
}

type postVisitResponse struct {
	Err error `json:"err,omitempty"`
}

func (v postVisitResponse) error() error { return v.Err }

type getVisit struct {
	ID string
}

type getVisitResponse struct {
	Visit Visit `json:"visit,omitempty"`
	Err   error `json:"err,omitempty"`
}

func (v getVisitResponse) error() error { return v.Err }

type postPatient struct {
	Patient Patient
}

type postPatientResponse struct {
	Err error `json:"err,omitempty"`
}

func (p postPatientResponse) error() error { return p.Err }

type getPatient struct {
	ID string
}

type getPatientResponse struct {
	Patient Patient `json:"patient,omitempty"`
	Err     error   `json:"err,omitempty"`
}

func (p getPatientResponse) error() error { return p.Err }
