package command

type ChangeVehicleStatusCommand struct {
	VehicleID string
	NewStatus string
}

func (c *ChangeVehicleStatusCommand) CommandName() string {
	return "ChangeVehicleStatus"
}
