package dto

type DisplayPatient struct {
	FirstName      string `json:"first_name"`
	MaskedLastName string `json:"last_name"`
}

type RoomDisplay struct {
	RoomID   string           `json:"room_id"`
	RoomName string           `json:"room_name"`
	Patients []DisplayPatient `json:"patients"`
}

type DisplayResponse struct {
	Rooms []RoomDisplay `json:"rooms"`
}
