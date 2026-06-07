package logger

// All helpers enforce privacy rules:
//   - national_id / passport_id: NEVER logged
//   - full patient name: NEVER logged
//   - internal queue value: NEVER logged

// ── Registration ──────────────────────────────────────────────────────────────

func Registered(lineID string, eventID uint) {
	Log.Info().
		Str("action", "registered").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Msg("patient registered successfully")
}

func AlreadyRegistered(lineID string, eventID uint) {
	Log.Warn().
		Str("action", "duplicate_registration").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Msg("patient already registered for this event")
}

func VerificationFailed(lineID string, eventID uint, reason string) {
	Log.Warn().
		Str("action", "verification_failed").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Str("detail", reason).
		Msg("patient not found in hospital records")
}

func MentalHealthNotCompleted(lineID string, eventID uint) {
	Log.Info().
		Str("action", "mental_health_not_completed").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Msg("mental health screening not completed — psychologist station required")
}

// ── Payment & Route ───────────────────────────────────────────────────────────

func PaymentScanReceived(lineID string, eventID uint) {
	Log.Info().
		Str("action", "payment_scan_received").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Msg("payment QR scan received")
}

func TransferForcedFalse(lineID string, eventID uint) {
	Log.Info().
		Str("action", "transfer_forced_false").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Str("detail", "passport holder — not eligible for rights transfer").
		Msg("needs_transfer overridden to false")
}

func RouteGenerated(lineID string, eventID uint, routeType string) {
	Log.Info().
		Str("action", "route_generated").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Str("route_type", routeType).
		Msg("patient route generated")
}

// ── Station Completion ────────────────────────────────────────────────────────

func StationCompleted(lineID string, eventID uint, station string) {
	Log.Info().
		Str("action", "station_completed").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Str("station", station).
		Msg("station marked as completed")
}

func StationNotInRoute(lineID string, eventID uint, station string) {
	Log.Warn().
		Str("action", "station_not_in_route").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Str("station", station).
		Msg("patient attempted to complete a station not in their route")
}

func StationPrerequisiteNotMet(lineID string, eventID uint, station, missing string) {
	Log.Warn().
		Str("action", "prerequisite_not_met").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Str("station", station).
		Str("missing", missing).
		Msg("required previous station not completed")
}

// ── Doctor Queue ──────────────────────────────────────────────────────────────

func QueueJoined(lineID string, eventID uint) {
	Log.Info().
		Str("action", "queue_joined").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Msg("patient joined doctor consultation queue")
}

func QueueDuplicate(lineID string, eventID uint) {
	Log.Warn().
		Str("action", "queue_duplicate").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Msg("patient already in queue — returning existing entry")
}

func QueueAssigned(lineID string, eventID uint, roomName string) {
	Log.Info().
		Str("action", "queue_assigned").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Str("room_name", roomName).
		Msg("patient assigned to doctor room")
}

func ConsultationCompleted(lineID string, eventID uint, roomName string) {
	Log.Info().
		Str("action", "consultation_completed").
		Str("line_id", lineID).
		Uint("event_id", eventID).
		Str("room_name", roomName).
		Msg("doctor consultation completed")
}

func RoomRefillTriggered(roomName, nextLineID string) {
	Log.Info().
		Str("action", "room_refill").
		Str("room_name", roomName).
		Str("next_line_id", nextLineID).
		Msg("room refill: next patient assigned")
}

func PatientSkipped(queueID uint) {
	Log.Info().
		Str("action", "patient_skipped").
		Uint("queue_id", queueID).
		Msg("patient queue entry skipped by staff")
}

// ── Errors ────────────────────────────────────────────────────────────────────

func DBError(action string, err error) {
	Log.Error().
		Err(err).
		Str("action", action).
		Msg("database error")
}

func ExcelError(file string, err error) {
	Log.Error().
		Err(err).
		Str("action", "excel_read_failed").
		Str("file", file).
		Msg("failed to read Excel file at startup")
}

func TxError(action string, err error) {
	Log.Error().
		Err(err).
		Str("action", action).
		Msg("database transaction failed")
}

func UnhandledError(action string, err error) {
	Log.Error().
		Err(err).
		Str("action", action).
		Msg("unhandled error")
}
