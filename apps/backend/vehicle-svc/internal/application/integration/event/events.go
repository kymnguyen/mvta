package event

type UserAuthorizedEvent struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	Timestamp int64  `json:"timestamp"`
}

type TrackingCorrectionAppliedEvent struct {
	VehicleID string `json:"vehicle_id"`
	Field     string `json:"field"` // e.g., "mileage", "fuel_level"
	OldValue  string `json:"old_value"`
	NewValue  string `json:"new_value"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}

type TrackingAlertEvent struct {
	VehicleID string `json:"vehicle_id"`
	AlertType string `json:"alert_type"` // e.g., "low_fuel", "high_mileage"
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
