package dto

type ScanAfterPaymentRequest struct {
	EventID       uint `json:"event_id"       binding:"required"`
	NeedsTransfer bool `json:"needs_transfer"`
}

type StationDTO struct {
	Order       int    `json:"order"`
	StationCode string `json:"station_code"`
	StationName string `json:"station_name"`
	Completed   bool   `json:"completed"`
	Required    bool   `json:"required"`
}

type RouteResponse struct {
	RouteType         string       `json:"route_type"`
	NeedsPsychologist bool         `json:"needs_psychologist"`
	NeedsTransfer     bool         `json:"needs_transfer"`
	CurrentStation    string       `json:"current_station"`
	NextStation       string       `json:"next_station,omitempty"`
	Steps             []StationDTO `json:"steps"`
	Message           string       `json:"message"`
}

type CompleteStationRequest struct {
	EventID uint `json:"event_id" binding:"required"`
}

type StationCompleteResponse struct {
	Message        string       `json:"message"`
	CurrentStation string       `json:"current_station"`
	NextStation    string       `json:"next_station,omitempty"`
	Steps          []StationDTO `json:"steps"`
}
