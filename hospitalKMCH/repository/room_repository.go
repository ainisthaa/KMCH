package repository

import (
	"context"

	"gorm.io/gorm"
	"lineoa-miniapp/domain"
)

type RoomRepository interface {
	FindAll(ctx context.Context) ([]domain.DoctorRoom, error)
	FindByID(ctx context.Context, roomID string) (*domain.DoctorRoom, error)
	UpdateTx(ctx context.Context, tx *gorm.DB, room *domain.DoctorRoom) error
	// Returns room with fewest ASSIGNED patients; tie-break by oldest last_completed_at then room_id
	LeastLoadedRoomTx(ctx context.Context, tx *gorm.DB) (*domain.DoctorRoom, error)
	// Returns room with oldest last_assigned_at (round-robin order)
	NextRoundRobinRoomTx(ctx context.Context, tx *gorm.DB) (*domain.DoctorRoom, error)
	CountAssignedInRoom(ctx context.Context, tx *gorm.DB, roomID string) (int64, error)
}
