package command

type UpdateVehicleFuelLevelCommand struct {
	VehicleID string
	FuelLevel float64
}

func (c *UpdateVehicleFuelLevelCommand) CommandName() string {
	return "UpdateVehicleFuelLevel"
}
