package service

import (
	"context"
	"fmt"
	"sort"

	"gorm.io/gorm"
	"lineoa-miniapp/domain"
	"lineoa-miniapp/dto"
	applog "lineoa-miniapp/pkg/logger"
	"lineoa-miniapp/pkg/mentalhealthcache"
	"lineoa-miniapp/repository"
)

type RouteService interface {
	ScanAfterPayment(ctx context.Context, lineID string, req dto.ScanAfterPaymentRequest) (*dto.RouteResponse, error)
	GetRoute(ctx context.Context, lineID string, eventID uint) (*dto.RouteResponse, error)
	CompletePsychologist(ctx context.Context, lineID string, eventID uint) (*dto.StationCompleteResponse, error)
	CompleteRightsTransfer(ctx context.Context, lineID string, eventID uint) (*dto.StationCompleteResponse, error)
}

type routeService struct {
	db          *gorm.DB
	patientRepo repository.PatientRepository
	mhCache     *mentalhealthcache.Cache
}

func NewRouteService(db *gorm.DB, patientRepo repository.PatientRepository, mhCache *mentalhealthcache.Cache) RouteService {
	return &routeService{db: db, patientRepo: patientRepo, mhCache: mhCache}
}

// ── Station templates ─────────────────────────────────────────────────────────

type stationTpl struct {
	Code     string
	Name     string
	Order    int
	Required bool
}

var routeTemplates = map[string][]stationTpl{
	domain.RouteA: {
		{domain.StationRegistration, "Registration", 1, true},
		{domain.StationPayment, "Payment", 2, true},
		{domain.StationDoctorConsultation, "Doctor consultation", 3, true},
		{domain.StationXray, "X-ray", 4, false},
	},
	domain.RouteB: {
		{domain.StationRegistration, "Registration", 1, true},
		{domain.StationPayment, "Payment", 2, true},
		{domain.StationRightsTransfer, "Healthcare rights transfer", 3, true},
		{domain.StationDoctorConsultation, "Doctor consultation", 4, true},
		{domain.StationXray, "X-ray", 5, false},
	},
	domain.RouteC: {
		{domain.StationRegistration, "Registration", 1, true},
		{domain.StationPayment, "Payment", 2, true},
		{domain.StationPsychologist, "Psychologist / mental health staff", 3, true},
		{domain.StationDoctorConsultation, "Doctor consultation", 4, true},
		{domain.StationXray, "X-ray", 5, false},
	},
	domain.RouteD: {
		{domain.StationRegistration, "Registration", 1, true},
		{domain.StationPayment, "Payment", 2, true},
		{domain.StationPsychologist, "Psychologist / mental health staff", 3, true},
		{domain.StationRightsTransfer, "Healthcare rights transfer", 4, true},
		{domain.StationDoctorConsultation, "Doctor consultation", 5, true},
		{domain.StationXray, "X-ray", 6, false},
	},
}

// ── ScanAfterPayment ──────────────────────────────────────────────────────────

func (s *routeService) ScanAfterPayment(ctx context.Context, lineID string, req dto.ScanAfterPaymentRequest) (*dto.RouteResponse, error) {
	applog.PaymentScanReceived(lineID, req.EventID)
	var resp *dto.RouteResponse

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		pc, err := s.patientRepo.FindCheckByLineAndEventTx(ctx, tx, lineID, req.EventID)
		if err != nil {
			applog.DBError("find_check_tx", err)
			return err
		}
		if pc == nil {
			return fmt.Errorf("patient not registered for this event")
		}
		if pc.IsPaid {
			resp = buildRouteResponse("Route already generated.", pc)
			return nil
		}

		lookupID := lineID
		patient, _ := s.patientRepo.FindByLineID(ctx, lineID)
		if patient != nil && patient.NationalID != "" {
			lookupID = patient.NationalID
		} else if patient != nil && patient.PassportID != "" {
			lookupID = patient.PassportID
		}

		needsTransfer := req.NeedsTransfer
		isPassportHolder := patient != nil && patient.NationalID == "" && patient.PassportID != ""
		if isPassportHolder {
			needsTransfer = false
			applog.TransferForcedFalse(lineID, req.EventID)
		}

		needsPsychologist := s.mhCache.NeedsPsychologist(lookupID)
		routeType := determineRouteType(needsPsychologist, needsTransfer)

		pc.IsPaid = true
		pc.NeedsTransfer = needsTransfer
		pc.NeedsPsychologist = needsPsychologist
		pc.RouteType = routeType
		if err := s.patientRepo.UpdateCheckTx(ctx, tx, pc); err != nil {
			applog.TxError("update_check_tx", err)
			return err
		}

		applog.RouteGenerated(lineID, req.EventID, routeType)
		resp = buildRouteResponse("Payment completed. Your route has been generated.", pc)
		return nil
	})

	return resp, err
}

// ── GetRoute ──────────────────────────────────────────────────────────────────

func (s *routeService) GetRoute(ctx context.Context, lineID string, eventID uint) (*dto.RouteResponse, error) {
	pc, err := s.patientRepo.FindCheckByLineAndEvent(ctx, lineID, eventID)
	if err != nil {
		return nil, err
	}
	if pc == nil {
		return nil, fmt.Errorf("patient not registered for this event")
	}
	if !pc.IsPaid {
		return nil, fmt.Errorf("payment has not been completed yet")
	}
	return buildRouteResponse("", pc), nil
}

// ── CompletePsychologist / CompleteRightsTransfer ─────────────────────────────

func (s *routeService) CompletePsychologist(ctx context.Context, lineID string, eventID uint) (*dto.StationCompleteResponse, error) {
	return s.completeStation(ctx, lineID, eventID, domain.StationPsychologist)
}

func (s *routeService) CompleteRightsTransfer(ctx context.Context, lineID string, eventID uint) (*dto.StationCompleteResponse, error) {
	return s.completeStation(ctx, lineID, eventID, domain.StationRightsTransfer)
}

func (s *routeService) completeStation(ctx context.Context, lineID string, eventID uint, stationCode string) (*dto.StationCompleteResponse, error) {
	var resp *dto.StationCompleteResponse

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		pc, err := s.patientRepo.FindCheckByLineAndEventTx(ctx, tx, lineID, eventID)
		if err != nil || pc == nil {
			return fmt.Errorf("patient not registered for this event")
		}
		if !pc.IsPaid {
			return fmt.Errorf("payment has not been completed yet")
		}

		steps := buildSteps(pc)

		if !hasStation(steps, stationCode) {
			applog.StationNotInRoute(lineID, eventID, stationCode)
			return fmt.Errorf("station '%s' is not part of this patient's route", stationCode)
		}
		if missing := firstMissingPrerequisite(steps, stationCode); missing != "" {
			applog.StationPrerequisiteNotMet(lineID, eventID, stationCode, missing)
			return fmt.Errorf("please complete the '%s' station first", missing)
		}

		switch stationCode {
		case domain.StationPsychologist:
			if pc.PsychologistDone {
				resp = stationAlreadyDone(pc)
				return nil
			}
			pc.PsychologistDone = true
		case domain.StationRightsTransfer:
			if pc.TransferCompleted {
				resp = stationAlreadyDone(pc)
				return nil
			}
			pc.TransferCompleted = true
		}

		if err := s.patientRepo.UpdateCheckTx(ctx, tx, pc); err != nil {
			applog.TxError("update_check_station", err)
			return err
		}

		applog.StationCompleted(lineID, eventID, stationCode)
		steps = buildSteps(pc)
		current, next := currentAndNext(steps)
		resp = &dto.StationCompleteResponse{
			Message:        fmt.Sprintf("Station '%s' completed.", stationCode),
			CurrentStation: current,
			NextStation:    next,
			Steps:          stepsToDTO(steps),
		}
		return nil
	})

	return resp, err
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func determineRouteType(needsPsychologist, needsTransfer bool) string {
	switch {
	case !needsPsychologist && !needsTransfer:
		return domain.RouteA
	case !needsPsychologist && needsTransfer:
		return domain.RouteB
	case needsPsychologist && !needsTransfer:
		return domain.RouteC
	default:
		return domain.RouteD
	}
}

type routeStep struct {
	Code      string
	Name      string
	Order     int
	Required  bool
	Completed bool
}

func buildSteps(pc *domain.PatientCheck) []routeStep {
	templates := routeTemplates[pc.RouteType]
	steps := make([]routeStep, len(templates))
	for i, t := range templates {
		var done bool
		switch t.Code {
		case domain.StationRegistration:
			done = true
		case domain.StationPayment:
			done = pc.IsPaid
		case domain.StationPsychologist:
			done = pc.PsychologistDone
		case domain.StationRightsTransfer:
			done = pc.TransferCompleted
		}
		steps[i] = routeStep{Code: t.Code, Name: t.Name, Order: t.Order, Required: t.Required, Completed: done}
	}
	return steps
}

func currentAndNext(steps []routeStep) (current, next string) {
	sort.Slice(steps, func(i, j int) bool { return steps[i].Order < steps[j].Order })
	for i, s := range steps {
		if !s.Completed {
			current = s.Code
			if i+1 < len(steps) {
				next = steps[i+1].Code
			}
			return
		}
	}
	current = "completed"
	return
}

func firstMissingPrerequisite(steps []routeStep, targetCode string) string {
	sort.Slice(steps, func(i, j int) bool { return steps[i].Order < steps[j].Order })
	targetOrder := 0
	for _, s := range steps {
		if s.Code == targetCode {
			targetOrder = s.Order
			break
		}
	}
	for _, s := range steps {
		if s.Order < targetOrder && s.Required && !s.Completed {
			return s.Code
		}
	}
	return ""
}

func hasStation(steps []routeStep, code string) bool {
	for _, s := range steps {
		if s.Code == code {
			return true
		}
	}
	return false
}

func stepsToDTO(steps []routeStep) []dto.StationDTO {
	sort.Slice(steps, func(i, j int) bool { return steps[i].Order < steps[j].Order })
	result := make([]dto.StationDTO, len(steps))
	for i, s := range steps {
		result[i] = dto.StationDTO{Order: s.Order, StationCode: s.Code, StationName: s.Name, Completed: s.Completed, Required: s.Required}
	}
	return result
}

func buildRouteResponse(msg string, pc *domain.PatientCheck) *dto.RouteResponse {
	steps := buildSteps(pc)
	current, next := currentAndNext(steps)
	return &dto.RouteResponse{
		RouteType:         pc.RouteType,
		NeedsPsychologist: pc.NeedsPsychologist,
		NeedsTransfer:     pc.NeedsTransfer,
		CurrentStation:    current,
		NextStation:       next,
		Steps:             stepsToDTO(steps),
		Message:           msg,
	}
}

func stationAlreadyDone(pc *domain.PatientCheck) *dto.StationCompleteResponse {
	steps := buildSteps(pc)
	current, next := currentAndNext(steps)
	return &dto.StationCompleteResponse{
		Message: "Station already completed.", CurrentStation: current, NextStation: next, Steps: stepsToDTO(steps),
	}
}
