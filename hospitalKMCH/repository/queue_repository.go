package repository

import (
	"context"

	"gorm.io/gorm"
	"lineoa-miniapp/domain"
)

type QueueRepository interface {
	Create(ctx context.Context, q *domain.PatientQueue) error
	FindActiveByLineAndEvent(ctx context.Context, lineID string, eventID uint) (*domain.PatientQueue, error)
	FindByLineAndEvent(ctx context.Context, lineID string, eventID uint) (*domain.PatientQueue, error)
	FindByQueueIDTx(ctx context.Context, tx *gorm.DB, queueID uint) (*domain.PatientQueue, error)
	// FOR UPDATE SKIP LOCKED — picks next WAITING entry for room assignment
	NextWaitingTx(ctx context.Context, tx *gorm.DB) (*domain.PatientQueue, error)
	UpdateTx(ctx context.Context, tx *gorm.DB, q *domain.PatientQueue) error
	CountWaiting(ctx context.Context) (int64, error)
	CountWaitingAhead(ctx context.Context, beforeQueue string) (int64, error)
	ListWaiting(ctx context.Context) ([]domain.PatientQueue, error)
	ListAssignedInRoom(ctx context.Context, roomID string) ([]domain.PatientQueue, error)
}
