package dto

type ScanDoctorQueueRequest struct {
	EventID uint `json:"event_id" binding:"required"`
}

type QueueStatusResponse struct {
	Status       string `json:"status"`
	TotalWaiting int64  `json:"total_waiting,omitempty"`
	Ahead        int64  `json:"ahead"`
	RoomName     string `json:"room_name,omitempty"`
	NextStation  string `json:"next_station,omitempty"`
	Message      string `json:"message"`
}

type CompleteConsultationRequest struct {
	EventID uint `json:"event_id" binding:"required"`
}
