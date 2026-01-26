package command

type UpdateVehicleMileageCommand struct {
	VehicleID string
	Mileage   float64
}

func (c *UpdateVehicleMileageCommand) CommandName() string {
	return "UpdateVehicleMileage"
}
