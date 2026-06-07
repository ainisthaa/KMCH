package service

import (
	"context"
	"unicode/utf8"

	"lineoa-miniapp/dto"
	"lineoa-miniapp/repository"
)

type DisplayService interface {
	GetDisplay(ctx context.Context) (*dto.DisplayResponse, error)
}

type displayService struct {
	roomRepo    repository.RoomRepository
	queueRepo   repository.QueueRepository
	patientRepo repository.PatientRepository
}

func NewDisplayService(
	roomRepo repository.RoomRepository,
	queueRepo repository.QueueRepository,
	patientRepo repository.PatientRepository,
) DisplayService {
	return &displayService{roomRepo: roomRepo, queueRepo: queueRepo, patientRepo: patientRepo}
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
			dp := dto.DisplayPatient{FirstName: e.LineID, MaskedLastName: ""}
			if info, err := s.patientRepo.FindByLineID(ctx, e.LineID); err == nil && info != nil {
				dp.FirstName = info.FirstName
				dp.MaskedLastName = maskLastName(info.LastName)
			}
			patients = append(patients, dp)
		}
		roomDisplays = append(roomDisplays, dto.RoomDisplay{
			RoomID:   room.RoomID,
			RoomName: room.RoomName,
			Patients: patients,
		})
	}

	return &dto.DisplayResponse{Rooms: roomDisplays}, nil
}

// maskLastName shows only the first character of the last name for the public display board.
func maskLastName(name string) string {
	if name == "" {
		return ""
	}
	_, size := utf8.DecodeRuneInString(name)
	return name[:size] + "."
}
