package command

type RecordVehicleChangeCommand struct {
	VehicleID  string
	VIN        string
	ChangeType string                       // created, location_updated, status_changed, mileage_updated, fuel_updated
	OldValue   map[string]interface{}
	NewValue   map[string]interface{}
	Version    int64
}

func (c *RecordVehicleChangeCommand) CommandName() string {
	return "RecordVehicleChange"
}
