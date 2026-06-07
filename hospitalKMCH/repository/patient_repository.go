package repository

import (
	"context"

	"gorm.io/gorm"
	"lineoa-miniapp/domain"
)

type PatientRepository interface {
	// PatientInfo
	FindByLineID(ctx context.Context, lineID string) (*domain.PatientInfo, error)
	UpsertPatient(ctx context.Context, p *domain.PatientInfo) error

	// PatientCheck
	FindCheckByLineAndEvent(ctx context.Context, lineID string, eventID uint) (*domain.PatientCheck, error)
	FindCheckByLineAndEventTx(ctx context.Context, tx *gorm.DB, lineID string, eventID uint) (*domain.PatientCheck, error)
	CreateCheck(ctx context.Context, pc *domain.PatientCheck) error
	UpdateCheck(ctx context.Context, pc *domain.PatientCheck) error
	UpdateCheckTx(ctx context.Context, tx *gorm.DB, pc *domain.PatientCheck) error
}
