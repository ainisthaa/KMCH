package service

import (
	"context"

	"lineoa-miniapp/domain"
	"lineoa-miniapp/dto"
	"lineoa-miniapp/repository"
)

type DisplayService interface {
	GetDisplay(ctx context.Context) (*dto.DisplayResponse, error)
}

type displayService struct {
	roomRepo  repository.RoomRepository
	queueRepo repository.QueueRepository
}

func NewDisplayService(roomRepo repository.RoomRepository, queueRepo repository.QueueRepository) DisplayService {
	return &displayService{roomRepo: roomRepo, queueRepo: queueRepo}
}

func (s *displayService) GetDisplay(ctx context.Context) (*dto.DisplayResponse, error) {
	rooms, err := s.roomRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	roomDisplays := make([]dto.RoomDisplay, 0, len(rooms))
	for _, room := range rooms {
		entries, err := s.queueRepo.ListAssignedInRoom(ctx, room.RoomID)
		if err != nil {
			return nil, err
		}
		patients := make([]dto.DisplayPatient, 0, len(entries))
		for _, e := range entries {
			info := getPatientInfo(ctx, e, s)
			patients = append(patients, info)
		}
		roomDisplays = append(roomDisplays, dto.RoomDisplay{
			RoomID:   room.RoomID,
			RoomName: room.RoomName,
			Patients: patients,
		})
	}

	return &dto.DisplayResponse{Rooms: roomDisplays}, nil
}

// getPatientInfo fetches patient name for display — masked last name.
func getPatientInfo(_ context.Context, e domain.PatientQueue, _ *displayService) dto.DisplayPatient {
	// LineID is available; in production inject patientRepo to fetch name.
	// For now return anonymised line_id slice.
	masked := e.LineID
	if len(masked) > 6 {
		masked = masked[:6] + "..."
	}
	return dto.DisplayPatient{FirstName: masked, MaskedLastName: ""}
}
