package dto

type RegisterRequest struct {
	LineID     string `json:"line_id"    binding:"required"`
	FirstName  string `json:"first_name"  binding:"required"`
	LastName   string `json:"last_name"   binding:"required"`
	TelNo      string `json:"tel_no"`
	ID         string `json:"id"          binding:"required"` // national_id or passport_id — backend detects
	StudentID  string `json:"student_id"`
	EmployeeID string `json:"employee_id"`
	EventID    uint   `json:"event_id"    binding:"required"`
}

type RegisterResponse struct {
	LineID     string `json:"line_id"`
	EventID    uint   `json:"event_id"`
	RouteReady bool   `json:"route_ready"` // false until scan-after-payment
	Message    string `json:"message"`
}
