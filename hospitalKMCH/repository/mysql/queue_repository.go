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

var _ repository.QueueRepository = (*MySQLQueueRepository)(nil)

type MySQLQueueRepository struct{ db *gorm.DB }

func NewMySQLQueueRepository(db *gorm.DB) *MySQLQueueRepository {
	return &MySQLQueueRepository{db: db}
}

func (r *MySQLQueueRepository) Create(ctx context.Context, q *domain.PatientQueue) error {
	now := time.Now()
	q.CreatedAt = now
	q.UpdatedAt = now
	// Internal FIFO key: nanosecond timestamp string — never returned to client
	q.Queue = now.Format("20060102150405.000000000")
	return r.db.WithContext(ctx).Create(q).Error
}

func (r *MySQLQueueRepository) FindActiveByLineAndEvent(ctx context.Context, lineID string, eventID uint) (*domain.PatientQueue, error) {
	var q domain.PatientQueue
	err := r.db.WithContext(ctx).
		Where("line_id = ? AND event_id = ? AND status IN (?,?)",
			lineID, eventID, domain.QueueWaiting, domain.QueueAssigned).
		First(&q).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &q, err
}

func (r *MySQLQueueRepository) FindByLineAndEvent(ctx context.Context, lineID string, eventID uint) (*domain.PatientQueue, error) {
	var q domain.PatientQueue
	err := r.db.WithContext(ctx).
		Preload("Room").
		Where("line_id = ? AND event_id = ?", lineID, eventID).
		Order("queue_id DESC").
		First(&q).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &q, err
}

func (r *MySQLQueueRepository) FindByQueueIDTx(ctx context.Context, tx *gorm.DB, queueID uint) (*domain.PatientQueue, error) {
	var q domain.PatientQueue
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("queue_id = ?", queueID).
		First(&q).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &q, err
}

// NextWaitingTx picks the oldest WAITING entry with no room assigned — FOR UPDATE SKIP LOCKED.
func (r *MySQLQueueRepository) NextWaitingTx(ctx context.Context, tx *gorm.DB) (*domain.PatientQueue, error) {
	var q domain.PatientQueue
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("status = ? AND room_id IS NULL", domain.QueueWaiting).
		Order("queue ASC").
		Limit(1).
		First(&q).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &q, err
}

func (r *MySQLQueueRepository) UpdateTx(ctx context.Context, tx *gorm.DB, q *domain.PatientQueue) error {
	q.UpdatedAt = time.Now()
	return tx.WithContext(ctx).Save(q).Error
}

func (r *MySQLQueueRepository) CountWaiting(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.PatientQueue{}).
		Where("status = ?", domain.QueueWaiting).
		Count(&count).Error
	return count, err
}

func (r *MySQLQueueRepository) CountWaitingAhead(ctx context.Context, beforeQueue string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.PatientQueue{}).
		Where("status = ? AND `queue` < ?", domain.QueueWaiting, beforeQueue).
		Count(&count).Error
	return count, err
}

func (r *MySQLQueueRepository) ListWaiting(ctx context.Context) ([]domain.PatientQueue, error) {
	var entries []domain.PatientQueue
	err := r.db.WithContext(ctx).
		Where("status = ?", domain.QueueWaiting).
		Order("queue ASC").
		Find(&entries).Error
	return entries, err
}

func (r *MySQLQueueRepository) ListAssignedInRoom(ctx context.Context, roomID string) ([]domain.PatientQueue, error) {
	var entries []domain.PatientQueue
	err := r.db.WithContext(ctx).
		Where("room_id = ? AND status = ?", roomID, domain.QueueAssigned).
		Order("queue ASC").
		Find(&entries).Error
	return entries, err
}
