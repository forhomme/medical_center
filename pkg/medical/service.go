package medical

import (
	"context"
	"errors"
	"sync"
)

// simple CRUD
type Service interface {
	PostVisit(ctx context.Context, patientID string, visit Visit) error
	GetVisit(ctx context.Context, id string) (Visit, error)
	PostPatient(ctx context.Context, p Patient) error
	GetPatient(ctx context.Context, id string) (Patient, error)
}

// for connection to db
type Repository interface {
	PostVisit(ctx context.Context, patientID string, visit Visit) error
	GetVisit(ctx context.Context, id string) (Visit, error)
	PostPatient(ctx context.Context, p Patient) error
	GetPatient(ctx context.Context, id string) (Patient, error)
}

type Visit struct {
	ID      string  `json:"id"`
	Patient Patient `json:"patient,omitempty"`
	Day     string  `json:"day,omitempty"`
	Time    string  `json:"time,omitempty"`
}

type Patient struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	Sex  string `json:"sex,omitempty"`
	Age  int    `json:"age,omitempty"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type medicalService struct {
	mtx        sync.RWMutex
	repository Repository
}

func NewMedicalService(rep Repository) Service {
	return &medicalService{
		repository: rep,
	}
}

func (m *medicalService) PostVisit(ctx context.Context, patientId string, v Visit) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	var patient Patient
	patient, err := m.repository.GetPatient(ctx, patientId)
	if err != nil {
		return ErrNotFound
	}
	visitDetail := Visit{
		ID:      v.ID,
		Patient: patient,
		Day:     v.Day,
		Time:    v.Time,
	}
	if err := m.repository.PostVisit(ctx, patientId, visitDetail); err != nil {
		return ErrBadRouting
	}
	return nil
}

func (m *medicalService) GetVisit(ctx context.Context, id string) (Visit, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	var obj Visit
	obj, err := m.repository.GetVisit(ctx, id)
	if err != nil {
		return Visit{}, ErrBadRouting
	}
	return obj, nil
}

func (m *medicalService) PostPatient(ctx context.Context, p Patient) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if err := m.repository.PostPatient(ctx, p); err != nil {
		return ErrBadRouting
	}
	return nil
}

func (m *medicalService) GetPatient(ctx context.Context, id string) (Patient, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	var obj Patient
	obj, err := m.repository.GetPatient(ctx, id)
	if err != nil {
		return Patient{}, ErrNotFound
	}
	return obj, nil
}
