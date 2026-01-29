package command

type UpdateVehicleLocationCommand struct {
	VehicleID string
	Latitude  float64
	Longitude float64
	Altitude  float64
	Timestamp int64
}

func (c *UpdateVehicleLocationCommand) CommandName() string {
	return "UpdateVehicleLocation"
}
