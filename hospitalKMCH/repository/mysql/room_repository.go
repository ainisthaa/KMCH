package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lineoa-miniapp/domain"
	"lineoa-miniapp/repository"
)

var _ repository.RoomRepository = (*MySQLRoomRepository)(nil)

type MySQLRoomRepository struct{ db *gorm.DB }

func NewMySQLRoomRepository(db *gorm.DB) *MySQLRoomRepository {
	return &MySQLRoomRepository{db: db}
}

func (r *MySQLRoomRepository) FindAll(ctx context.Context) ([]domain.DoctorRoom, error) {
	var rooms []domain.DoctorRoom
	err := r.db.WithContext(ctx).Order("room_id ASC").Find(&rooms).Error
	return rooms, err
}

func (r *MySQLRoomRepository) FindByID(ctx context.Context, roomID string) (*domain.DoctorRoom, error) {
	var room domain.DoctorRoom
	err := r.db.WithContext(ctx).Where("room_id = ?", roomID).First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &room, err
}

func (r *MySQLRoomRepository) UpdateTx(ctx context.Context, tx *gorm.DB, room *domain.DoctorRoom) error {
	return tx.WithContext(ctx).Model(room).
		Updates(map[string]interface{}{
			"last_assigned_at":  room.LastAssignedAt,
			"last_completed_at": room.LastCompletedAt,
		}).Error
}

// LeastLoadedRoomTx: room with fewest ASSIGNED patients.
// Tie-break 1: oldest last_completed_at (NULL first). Tie-break 2: smallest room_id.
// Uses SELECT ... FOR UPDATE to lock rows during assignment.
func (r *MySQLRoomRepository) LeastLoadedRoomTx(ctx context.Context, tx *gorm.DB) (*domain.DoctorRoom, error) {
	var room domain.DoctorRoom
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Raw(`
			SELECT dr.*
			FROM doctor_room dr
			LEFT JOIN (
				SELECT room_id, COUNT(*) AS active_count
				FROM patient_queue
				WHERE status = 'assigned'
				GROUP BY room_id
			) pq ON pq.room_id = dr.room_id
			ORDER BY
				COALESCE(pq.active_count, 0) ASC,
				COALESCE(dr.last_completed_at, '1970-01-01 00:00:00') ASC,
				dr.room_id ASC
			LIMIT 1
		`).
		Scan(&room).Error
	if err != nil {
		return nil, err
	}
	if room.RoomID == "" {
		return nil, nil
	}
	return &room, nil
}

// NextRoundRobinRoomTx: room with oldest last_assigned_at (NULL = never assigned = first).
func (r *MySQLRoomRepository) NextRoundRobinRoomTx(ctx context.Context, tx *gorm.DB) (*domain.DoctorRoom, error) {
	var room domain.DoctorRoom
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Order("COALESCE(last_assigned_at, '1970-01-01 00:00:00') ASC, room_id ASC").
		First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &room, err
}

func (r *MySQLRoomRepository) CountAssignedInRoom(ctx context.Context, tx *gorm.DB, roomID string) (int64, error) {
	var count int64
	db := r.db
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Model(&domain.PatientQueue{}).
		Where("room_id = ? AND status = ?", roomID, domain.QueueAssigned).
		Count(&count).Error
	return count, err
}
