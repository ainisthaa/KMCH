package service

import (
	"context"

	"lineoa-miniapp/dto"
	"lineoa-miniapp/repository"
)

const maxPerRoom = MaxPatientsPerRoom

type StaffService interface {
	GetDashboard(ctx context.Context) (*dto.StaffDashboardResponse, error)
	GetWaitingQueue(ctx context.Context) ([]dto.StaffQueueItem, error)
}

type staffService struct {
	roomRepo  repository.RoomRepository
	queueRepo repository.QueueRepository
}

func NewStaffService(roomRepo repository.RoomRepository, queueRepo repository.QueueRepository) StaffService {
	return &staffService{roomRepo: roomRepo, queueRepo: queueRepo}
}

func (s *staffService) GetDashboard(ctx context.Context) (*dto.StaffDashboardResponse, error) {
	rooms, err := s.roomRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var totalAssigned int64
	roomStatuses := make([]dto.StaffRoomStatus, 0, len(rooms))
	for _, room := range rooms {
		active, err := s.queueRepo.ListAssignedInRoom(ctx, room.RoomID)
		if err != nil {
			return nil, err
		}
		count := int64(len(active))
		totalAssigned += count

		patients := make([]dto.DisplayPatient, 0, len(active))
		for _, e := range active {
			masked := e.LineID
			if len(masked) > 6 {
				masked = masked[:6] + "..."
			}
			patients = append(patients, dto.DisplayPatient{FirstName: masked})
		}

		roomStatuses = append(roomStatuses, dto.StaffRoomStatus{
			RoomID:         room.RoomID,
			RoomName:       room.RoomName,
			ActiveCount:    count,
			AvailableSlots: maxPerRoom - count,
			Patients:       patients,
		})
	}

	waiting, _ := s.queueRepo.CountWaiting(ctx)

	return &dto.StaffDashboardResponse{
		WaitingCount:  waiting,
		AssignedCount: totalAssigned,
		Rooms:         roomStatuses,
	}, nil
}

func (s *staffService) GetWaitingQueue(ctx context.Context) ([]dto.StaffQueueItem, error) {
	entries, err := s.queueRepo.ListWaiting(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]dto.StaffQueueItem, 0, len(entries))
	for _, e := range entries {
		items = append(items, dto.StaffQueueItem{
			QueueID: e.QueueID,
			Status:  e.Status,
		})
	}
	return items, nil
}
