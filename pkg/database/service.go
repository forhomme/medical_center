package database

import (
	"context"
	"database/sql"
	"errors"
	"test/pkg/medical"
)

var (
	ErrDatabase   = errors.New("Unable to handle Database Request")
	ErrIdNotFound = errors.New("Id not found")
)

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) (medical.Repository, error) {
	return &repo{
		db: db,
	}, nil
}

func (repo *repo) PostVisit(ctx context.Context, patientID string, v medical.Visit) error {
	patient := medical.Patient{}
	err := repo.db.QueryRowContext(ctx, "SELECT * FROM patient WHERE id = ?", patientID).Scan(&patient)
	if err != nil {
		return ErrIdNotFound
	}
	_, errs := repo.db.ExecContext(ctx, "INSERT INTO visit(id, patient, day, time) VALUES (?, ?, ?, ?)",
		v.ID, patient, v.Day, v.Time)
	if errs != nil {
		return ErrDatabase
	}
	return nil
}

func (repo *repo) GetVisit(ctx context.Context, visitID string) (medical.Visit, error) {
	visit := medical.Visit{}
	err := repo.db.QueryRowContext(ctx, "SELECT * FROM patient WHERE id = ?", visitID).Scan(&visit)
	if err != nil {
		return medical.Visit{}, ErrIdNotFound
	}
	return visit, nil
}

func (repo *repo) PostPatient(ctx context.Context, p medical.Patient) error {
	_, errs := repo.db.ExecContext(ctx, "INSERT INTO visit(id, name, sex, age) VALUES (?, ?, ?, ?)",
		p.ID, p.Name, p.Sex, p.Age)
	if errs != nil {
		return ErrDatabase
	}
	return nil
}

func (repo *repo) GetPatient(ctx context.Context, patientID string) (medical.Patient, error) {
	patient := medical.Patient{}
	err := repo.db.QueryRowContext(ctx, "SELECT * FROM patient WHERE id = ?", patientID).Scan(&patient)
	if err != nil {
		return medical.Patient{}, ErrIdNotFound
	}
	return patient, nil
}
