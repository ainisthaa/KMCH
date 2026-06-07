package service

import (
	"context"
	"fmt"
	"time"

	"lineoa-miniapp/domain"
	"lineoa-miniapp/dto"
	"lineoa-miniapp/pkg/idutil"
	applog "lineoa-miniapp/pkg/logger"
	"lineoa-miniapp/pkg/mentalhealthcache"
	"lineoa-miniapp/repository"
	excelrepo "lineoa-miniapp/repository/excel"
)

type RegisterService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
}

type registerService struct {
	patientRepo      repository.PatientRepository
	registrationRepo excelrepo.RegistrationRepository
	mhCache          *mentalhealthcache.Cache
}

func NewRegisterService(
	patientRepo repository.PatientRepository,
	registrationRepo excelrepo.RegistrationRepository,
	mhCache *mentalhealthcache.Cache,
) RegisterService {
	return &registerService{
		patientRepo:      patientRepo,
		registrationRepo: registrationRepo,
		mhCache:          mhCache,
	}
}

func (s *registerService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	nationalID, passportID := idutil.ClassifyID(req.ID)
	lookupID := req.ID
	if lookupID == "" {
		return nil, fmt.Errorf("id is required")
	}

	found, err := s.registrationRepo.Exists(ctx, lookupID)
	if err != nil {
		applog.UnhandledError("registration_excel_check", err)
		return nil, fmt.Errorf("registration check: %w", err)
	}
	if !found {
		applog.VerificationFailed(req.LineID, req.EventID, "not in hospital records")
		return nil, fmt.Errorf("patient not found in hospital records")
	}

	existing, err := s.patientRepo.FindCheckByLineAndEvent(ctx, req.LineID, req.EventID)
	if err != nil {
		applog.DBError("find_check_by_line_event", err)
		return nil, err
	}
	if existing != nil {
		applog.AlreadyRegistered(req.LineID, req.EventID)
		return &dto.RegisterResponse{
			LineID:     req.LineID,
			EventID:    req.EventID,
			RouteReady: false,
			Message:    "Already registered for this event.",
		}, nil
	}

	now := time.Now()
	patient := &domain.PatientInfo{
		LineID:       req.LineID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		TelNo:        req.TelNo,
		NationalID:   nationalID,
		PassportID:   passportID,
		RegisterDate: &now,
		StudentID:    req.StudentID,
		EmployeeID:   req.EmployeeID,
	}
	if err := s.patientRepo.UpsertPatient(ctx, patient); err != nil {
		applog.DBError("upsert_patient", err)
		return nil, fmt.Errorf("upsert patient: %w", err)
	}

	psyevalForm := s.mhCache.HasCompletedScreening(lookupID)
	isSV := s.mhCache.HasIssue(lookupID)

	if !psyevalForm {
		applog.MentalHealthNotCompleted(req.LineID, req.EventID)
	}

	pc := &domain.PatientCheck{
		LineID:      req.LineID,
		EventID:     req.EventID,
		PsyevalForm: psyevalForm,
		IsSV:        isSV,
	}
	if err := s.patientRepo.CreateCheck(ctx, pc); err != nil {
		applog.DBError("create_check", err)
		return nil, fmt.Errorf("create check: %w", err)
	}

	applog.Registered(req.LineID, req.EventID)
	return &dto.RegisterResponse{
		LineID:     req.LineID,
		EventID:    req.EventID,
		RouteReady: false,
		Message:    "Registration successful. Please proceed to the payment counter.",
	}, nil
}
