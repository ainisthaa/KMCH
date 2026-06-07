package dto

type StaffQueueItem struct {
	QueueID        uint   `json:"queue_id"`
	FirstName      string `json:"first_name"`
	MaskedLastName string `json:"last_name"`
	Status         string `json:"status"`
	RoomName       string `json:"room_name,omitempty"`
}

type StaffRoomStatus struct {
	RoomID        string           `json:"room_id"`
	RoomName      string           `json:"room_name"`
	ActiveCount   int64            `json:"active_count"`
	AvailableSlots int64           `json:"available_slots"`
	Patients      []DisplayPatient `json:"patients"`
}

type StaffDashboardResponse struct {
	WaitingCount  int64             `json:"waiting_count"`
	AssignedCount int64             `json:"assigned_count"`
	Rooms         []StaffRoomStatus `json:"rooms"`
}
