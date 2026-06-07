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
	roomRepo    repository.RoomRepository
	queueRepo   repository.QueueRepository
	patientRepo repository.PatientRepository
}

func NewStaffService(
	roomRepo repository.RoomRepository,
	queueRepo repository.QueueRepository,
	patientRepo repository.PatientRepository,
) StaffService {
	return &staffService{roomRepo: roomRepo, queueRepo: queueRepo, patientRepo: patientRepo}
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
			dp := dto.DisplayPatient{}
			if info, err := s.patientRepo.FindByLineID(ctx, e.LineID); err == nil && info != nil {
				dp.FirstName = info.FirstName
				dp.MaskedLastName = maskLastName(info.LastName)
			}
			patients = append(patients, dp)
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
		item := dto.StaffQueueItem{
			QueueID: e.QueueID,
			Status:  e.Status,
		}
		if info, err := s.patientRepo.FindByLineID(ctx, e.LineID); err == nil && info != nil {
			item.FirstName = info.FirstName
			item.MaskedLastName = maskLastName(info.LastName)
		}
		items = append(items, item)
	}
	return items, nil
}
