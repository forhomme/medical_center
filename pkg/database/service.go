package database

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"test/pkg/medical"
)

var (
	ErrDatabase   = errors.New("Unable to handle Database Request")
	ErrIdNotFound = errors.New("Id not found")
)

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) (medical.Repository, error) {
	return &repo{
		db: db,
	}, nil
}

func (repo *repo) PostVisit(_ context.Context, v medical.Visit) error {
	err := repo.db.Create(&v).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetVisit(_ context.Context, visitID string) (medical.Visit, error) {
	visit := medical.Visit{}
	err := repo.db.Preload("Patient").First(&visit, visitID).Error
	if err != nil {
		return medical.Visit{}, ErrIdNotFound
	}
	return visit, nil
}

func (repo *repo) PostPatient(_ context.Context, p medical.Patient) error {
	err := repo.db.Create(&p).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetPatient(_ context.Context, patientID string) (medical.Patient, error) {
	patient := medical.Patient{}
	err := repo.db.First(&patient, patientID).Error
	if err != nil {
		return medical.Patient{}, ErrIdNotFound
	}
	return patient, nil
}
