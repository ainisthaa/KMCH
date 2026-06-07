package mysql

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lineoa-miniapp/domain"
	"lineoa-miniapp/repository"
)

var _ repository.PatientRepository = (*MySQLPatientRepository)(nil)

type MySQLPatientRepository struct{ db *gorm.DB }

func NewMySQLPatientRepository(db *gorm.DB) *MySQLPatientRepository {
	return &MySQLPatientRepository{db: db}
}

func (r *MySQLPatientRepository) FindByLineID(ctx context.Context, lineID string) (*domain.PatientInfo, error) {
	var p domain.PatientInfo
	err := r.db.WithContext(ctx).Where("line_id = ?", lineID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *MySQLPatientRepository) UpsertPatient(ctx context.Context, p *domain.PatientInfo) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "line_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"first_name", "last_name", "tel_no", "national_id",
				"passport_id", "register_date", "student_id", "employee_id",
			}),
		}).
		Create(p).Error
}

func (r *MySQLPatientRepository) FindCheckByLineAndEvent(ctx context.Context, lineID string, eventID uint) (*domain.PatientCheck, error) {
	var pc domain.PatientCheck
	err := r.db.WithContext(ctx).
		Where("line_id = ? AND event_id = ?", lineID, eventID).
		First(&pc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pc, err
}

func (r *MySQLPatientRepository) FindCheckByLineAndEventTx(ctx context.Context, tx *gorm.DB, lineID string, eventID uint) (*domain.PatientCheck, error) {
	var pc domain.PatientCheck
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("line_id = ? AND event_id = ?", lineID, eventID).
		First(&pc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pc, err
}

func (r *MySQLPatientRepository) CreateCheck(ctx context.Context, pc *domain.PatientCheck) error {
	now := time.Now()
	pc.CreatedAt = now
	pc.UpdatedAt = now
	return r.db.WithContext(ctx).Create(pc).Error
}

func (r *MySQLPatientRepository) UpdateCheck(ctx context.Context, pc *domain.PatientCheck) error {
	pc.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(pc).Error
}

func (r *MySQLPatientRepository) UpdateCheckTx(ctx context.Context, tx *gorm.DB, pc *domain.PatientCheck) error {
	pc.UpdatedAt = time.Now()
	return tx.WithContext(ctx).Save(pc).Error
}
