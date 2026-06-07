package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"lineoa-miniapp/domain"
	"lineoa-miniapp/repository"
)

var _ repository.EventRepository = (*MySQLEventRepository)(nil)

type MySQLEventRepository struct{ db *gorm.DB }

func NewMySQLEventRepository(db *gorm.DB) *MySQLEventRepository {
	return &MySQLEventRepository{db: db}
}

func (r *MySQLEventRepository) FindByID(ctx context.Context, eventID uint) (*domain.EventInfo, error) {
	var e domain.EventInfo
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).First(&e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &e, err
}

func (r *MySQLEventRepository) FindAll(ctx context.Context) ([]domain.EventInfo, error) {
	var events []domain.EventInfo
	err := r.db.WithContext(ctx).Order("event_id ASC").Find(&events).Error
	return events, err
}
