package medical

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)

// simple CRUD
type Service interface {
	PostVisit(ctx context.Context, v Visit) error
	GetVisit(ctx context.Context, id string) (Visit, error)
	PostPatient(ctx context.Context, p Patient) error
	GetPatient(ctx context.Context, id string) (Patient, error)
}

type Visit struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	PatientID int       `json:"patient_id"`
	Patient   Patient   `json:"patient" gorm:"foreignKey:PatientID"`
	Schedule  time.Time `json:"schedule"`
}

type Patient struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	Sex  string `json:"sex"`
	Age  int    `json:"age"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type medicalService struct {
	mtx sync.RWMutex
	db  *gorm.DB
}

func NewMedicalService(db *gorm.DB) Service {
	return &medicalService{
		db: db,
	}
}

func nestedLoop(h []int, m []int, v Visit) Visit {
	var visit Visit
	for _, a := range h {
		for _, b := range m {
			if v.Schedule.Hour() == a && v.Schedule.Minute() == b {
				visit = v
				return visit
			}
		}
	}
	return visit
}

func checkSchedule(v Visit) (Visit, error) {
	if v.Schedule.IsZero() {
		return Visit{}, ErrNotFound
	}
	day := v.Schedule.Weekday()
	schedule := Visit{}
	switch day {
	case time.Monday:
		hour := []int{8, 9, 10, 14, 15, 20, 21}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		schedule = nestedLoop(hour, minute, v)
	case time.Tuesday:
		hour := []int{10, 11, 15, 16, 17}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		schedule = nestedLoop(hour, minute, v)
	case time.Wednesday:
		hour := []int{13, 14, 15, 16, 17}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		schedule = nestedLoop(hour, minute, v)
	case time.Thursday:
		hour := []int{8, 9, 10, 14, 15, 16, 17}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		schedule = nestedLoop(hour, minute, v)
	case time.Friday:
		hour := []int{14, 15, 16, 17}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		schedule = nestedLoop(hour, minute, v)
	case time.Saturday:
		hour := []int{8, 9, 10}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		schedule = nestedLoop(hour, minute, v)
	case time.Sunday:
		hour := []int{20, 21, 22}
		minute := make([]int, 60)
		for i := range minute {
			minute[i] = i
		}
		schedule = nestedLoop(hour, minute, v)
	}
	return schedule, nil
}

func (m *medicalService) PostVisit(ctx context.Context, v Visit) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	visit, errs := checkSchedule(v)
	if errs != nil {
		return errs
	}
	visitDetail := Visit{
		Patient:  visit.Patient,
		Schedule: visit.Schedule,
	}
	if err := m.db.Create(&visitDetail).Error; err != nil {
		return err
	}
	return nil
}

func (m *medicalService) GetVisit(ctx context.Context, id string) (Visit, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	var obj Visit
	err := m.db.Preload("Patient").First(&obj, id).Error
	if err != nil {
		return Visit{}, err
	}
	return obj, nil
}

func (m *medicalService) PostPatient(ctx context.Context, p Patient) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if err := m.db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func (m *medicalService) GetPatient(ctx context.Context, id string) (Patient, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	var obj Patient
	if err := m.db.First(&obj, id).Error; err != nil {
		return Patient{}, err
	}
	return obj, nil
}
