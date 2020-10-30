package medical

import (
	"context"
	"errors"
	"sync"
	"time"
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
	ID       string    `json:"id"`
	Patient  Patient   `json:"patient,omitempty"`
	Schedule time.Time `json:"schedule,omitempty"`
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

func checkSchedule(v Visit) (Visit, error) {
	if v.Schedule.IsZero() {
		return Visit{}, ErrBadRouting
	}
	day := v.Schedule.Weekday()
	schedule := Visit{}
	switch day {
	case time.Monday:
		hour := [7]int{8, 9, 10, 14, 15, 20, 21}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		for _, a := range hour {
			for _, b := range minute {
				if a == v.Schedule.Hour() && b == v.Schedule.Minute() {
					schedule = v
					break
				} else {
					schedule = Visit{}
				}
			}
		}
	case time.Tuesday:
		hour := [5]int{10, 11, 15, 16, 17}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		for _, a := range hour {
			for _, b := range minute {
				if a == v.Schedule.Hour() && b == v.Schedule.Minute() {
					schedule = v
					break
				} else {
					schedule = Visit{}
				}
			}
		}
	case time.Wednesday:
		hour := [5]int{13, 14, 15, 16, 17}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		for _, a := range hour {
			for _, b := range minute {
				if a == v.Schedule.Hour() && b == v.Schedule.Minute() {
					schedule = v
					break
				} else {
					schedule = Visit{}
				}
			}
		}
	case time.Thursday:
		hour := [7]int{8, 9, 10, 14, 15, 16, 17}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		for _, a := range hour {
			for _, b := range minute {
				if a == v.Schedule.Hour() && b == v.Schedule.Minute() {
					schedule = v
					break
				} else {
					schedule = Visit{}
				}
			}
		}
	case time.Friday:
		hour := [4]int{14, 15, 16, 17}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		for _, a := range hour {
			for _, b := range minute {
				if a == v.Schedule.Hour() && b == v.Schedule.Minute() {
					schedule = v
					break
				} else {
					schedule = Visit{}
				}
			}
		}
	case time.Saturday:
		hour := [3]int{8, 9, 10}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		for _, a := range hour {
			for _, b := range minute {
				if a == v.Schedule.Hour() && b == v.Schedule.Minute() {
					schedule = v
					break
				} else {
					schedule = Visit{}
				}
			}
		}
	case time.Sunday:
		hour := [3]int{20, 21, 22}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		for _, a := range hour {
			for _, b := range minute {
				if a == v.Schedule.Hour() && b == v.Schedule.Minute() {
					schedule = v
					break
				} else {
					schedule = Visit{}
				}
			}
		}

	}
	return schedule, nil
}

func (m *medicalService) PostVisit(ctx context.Context, patientId string, v Visit) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	var patient Patient
	patient, err := m.repository.GetPatient(ctx, patientId)
	if err != nil {
		return ErrNotFound
	}
	visit, errs := checkSchedule(v)
	if errs != nil {
		return errs
	}
	visitDetail := Visit{
		ID:       visit.ID,
		Patient:  patient,
		Schedule: visit.Schedule,
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
		return Visit{}, err
	}
	return obj, nil
}

func (m *medicalService) PostPatient(ctx context.Context, p Patient) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if err := m.repository.PostPatient(ctx, p); err != nil {
		return err
	}
	return nil
}

func (m *medicalService) GetPatient(ctx context.Context, id string) (Patient, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	var obj Patient
	obj, err := m.repository.GetPatient(ctx, id)
	if err != nil {
		return Patient{}, err
	}
	return obj, nil
}
