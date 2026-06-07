package service

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"lineoa-miniapp/domain"
	"lineoa-miniapp/dto"
	applog "lineoa-miniapp/pkg/logger"
	"lineoa-miniapp/repository"
)

const (
	MaxPatientsPerRoom = 4
	LowQueueThreshold  = 5
)

type QueueService interface {
	ScanDoctorQueue(ctx context.Context, lineID string, eventID uint) error
	GetQueueStatus(ctx context.Context, lineID string, eventID uint) (*dto.QueueStatusResponse, error)
	CompleteDoctorConsultation(ctx context.Context, lineID string, eventID uint) error
	SkipQueue(ctx context.Context, queueID uint) error
	AutoFillRooms(ctx context.Context) error
}

type queueService struct {
	db          *gorm.DB
	patientRepo repository.PatientRepository
	queueRepo   repository.QueueRepository
	roomRepo    repository.RoomRepository
}

func NewQueueService(
	db *gorm.DB,
	patientRepo repository.PatientRepository,
	queueRepo repository.QueueRepository,
	roomRepo repository.RoomRepository,
) QueueService {
	return &queueService{db: db, patientRepo: patientRepo, queueRepo: queueRepo, roomRepo: roomRepo}
}

func (s *queueService) ScanDoctorQueue(ctx context.Context, lineID string, eventID uint) error {
	pc, err := s.patientRepo.FindCheckByLineAndEvent(ctx, lineID, eventID)
	if err != nil {
		applog.DBError("find_check_by_line_event", err)
		return err
	}
	if pc == nil {
		return fmt.Errorf("patient not registered for this event")
	}
	if !pc.IsPaid {
		return fmt.Errorf("payment has not been completed yet")
	}

	steps := buildSteps(pc)
	if missing := firstMissingPrerequisite(steps, domain.StationDoctorConsultation); missing != "" {
		applog.StationPrerequisiteNotMet(lineID, eventID, domain.StationDoctorConsultation, missing)
		return fmt.Errorf("please complete the '%s' station first", missing)
	}

	existing, err := s.queueRepo.FindActiveByLineAndEvent(ctx, lineID, eventID)
	if err != nil {
		return err
	}
	if existing != nil {
		applog.QueueDuplicate(lineID, eventID)
		return nil
	}

	now := time.Now()
	q := &domain.PatientQueue{
		LineID:     lineID,
		EventID:    eventID,
		Status:     domain.QueueWaiting,
		Station:    domain.QStationQueue,
		QStartTime: &now,
	}
	if err := s.queueRepo.Create(ctx, q); err != nil {
		applog.DBError("create_queue_entry", err)
		return fmt.Errorf("create queue entry: %w", err)
	}
	applog.QueueJoined(lineID, eventID)

	go func() {
		if err := s.AutoFillRooms(context.Background()); err != nil {
			applog.UnhandledError("auto_fill_after_join", err)
		}
	}()

	return nil
}

func (s *queueService) GetQueueStatus(ctx context.Context, lineID string, eventID uint) (*dto.QueueStatusResponse, error) {
	q, err := s.queueRepo.FindByLineAndEvent(ctx, lineID, eventID)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return &dto.QueueStatusResponse{
			Status:  "not_in_queue",
			Message: "You have not joined the doctor consultation queue yet.",
		}, nil
	}

	switch q.Status {
	case domain.QueueWaiting:
		ahead, _ := s.queueRepo.CountWaitingAhead(ctx, q.Queue)
		total, _ := s.queueRepo.CountWaiting(ctx)
		return &dto.QueueStatusResponse{
			Status:       domain.QueueWaiting,
			TotalWaiting: total,
			Ahead:        ahead,
			Message:      fmt.Sprintf("You are in the queue. %d patient(s) ahead of you.", ahead),
		}, nil

	case domain.QueueAssigned:
		out := &dto.QueueStatusResponse{
			Status:  domain.QueueAssigned,
			Ahead:   0,
			Message: "Your queue has been called.",
		}
		if q.Room != nil {
			out.RoomName = q.Room.RoomName
			out.Message = fmt.Sprintf("Please go to %s.", q.Room.RoomName)
		}
		return out, nil

	case domain.QueueCompleted:
		return &dto.QueueStatusResponse{
			Status:      domain.QueueCompleted,
			NextStation: domain.StationXray,
			Message:     "Doctor consultation completed. Please proceed to X-ray.",
		}, nil

	case domain.QueueSkip:
		return &dto.QueueStatusResponse{
			Status:  domain.QueueSkip,
			Message: "Your queue was skipped. Please contact staff.",
		}, nil
	}

	return &dto.QueueStatusResponse{Status: q.Status, Message: "Unknown status."}, nil
}

func (s *queueService) CompleteDoctorConsultation(ctx context.Context, lineID string, eventID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		active, err := s.queueRepo.FindActiveByLineAndEvent(ctx, lineID, eventID)
		if err != nil || active == nil {
			return fmt.Errorf("no active queue entry found")
		}
		if active.Status != domain.QueueAssigned {
			return fmt.Errorf("patient is not currently assigned to a room")
		}

		q, err := s.queueRepo.FindByQueueIDTx(ctx, tx, active.QueueID)
		if err != nil || q == nil {
			return fmt.Errorf("queue entry not found")
		}

		now := time.Now()
		q.Status = domain.QueueCompleted
		q.Station = domain.QStationCompleted
		q.QEndTime = &now
		if err := s.queueRepo.UpdateTx(ctx, tx, q); err != nil {
			applog.TxError("complete_consultation", err)
			return err
		}

		roomName := ""
		if q.RoomID != nil {
			room, _ := s.roomRepo.FindByID(ctx, *q.RoomID)
			if room != nil {
				roomName = room.RoomName
				room.LastCompletedAt = &now
				_ = s.roomRepo.UpdateTx(ctx, tx, room)
			}
		}

		applog.ConsultationCompleted(lineID, eventID, roomName)
		return s.fillRoomsTx(ctx, tx)
	})
}

func (s *queueService) SkipQueue(ctx context.Context, queueID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		q, err := s.queueRepo.FindByQueueIDTx(ctx, tx, queueID)
		if err != nil || q == nil {
			return fmt.Errorf("queue entry not found")
		}
		if q.Status != domain.QueueWaiting && q.Status != domain.QueueAssigned {
			return fmt.Errorf("cannot skip a completed entry")
		}
		now := time.Now()
		q.Status = domain.QueueSkip
		q.QEndTime = &now
		if err := s.queueRepo.UpdateTx(ctx, tx, q); err != nil {
			applog.TxError("skip_queue", err)
			return err
		}
		applog.PatientSkipped(queueID)
		return s.fillRoomsTx(ctx, tx)
	})
}

func (s *queueService) AutoFillRooms(ctx context.Context) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.fillRoomsTx(ctx, tx)
	})
}

func (s *queueService) fillRoomsTx(ctx context.Context, tx *gorm.DB) error {
	for {
		next, err := s.queueRepo.NextWaitingTx(ctx, tx)
		if err != nil || next == nil {
			break
		}

		totalWaiting, _ := s.queueRepo.CountWaiting(ctx)

		var room *domain.DoctorRoom
		if totalWaiting > LowQueueThreshold {
			room, err = s.roomRepo.NextRoundRobinRoomTx(ctx, tx)
		} else {
			room, err = s.roomRepo.LeastLoadedRoomTx(ctx, tx)
		}
		if err != nil || room == nil {
			break
		}

		activeCount, _ := s.roomRepo.CountAssignedInRoom(ctx, tx, room.RoomID)
		if activeCount >= MaxPatientsPerRoom {
			break
		}

		now := time.Now()
		next.Status = domain.QueueAssigned
		next.RoomID = &room.RoomID
		if err := s.queueRepo.UpdateTx(ctx, tx, next); err != nil {
			applog.TxError("fill_rooms_assign", err)
			return err
		}

		room.LastAssignedAt = &now
		if err := s.roomRepo.UpdateTx(ctx, tx, room); err != nil {
			applog.TxError("fill_rooms_update_room", err)
			return err
		}

		applog.RoomRefillTriggered(room.RoomName, next.LineID)
		applog.QueueAssigned(next.LineID, next.EventID, room.RoomName)
	}
	return nil
}
