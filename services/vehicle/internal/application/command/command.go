package command

import "context"

type Command interface {
	CommandName() string
}

type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

type CommandBus interface {
	Dispatch(ctx context.Context, cmd Command) error

	Register(commandName string, handler CommandHandler)
}

type CreateVehicleCommand struct {
	VIN           string
	VehicleName   string
	VehicleModel  string
	LicenseNumber string
	Status        string
	Latitude      float64
	Longitude     float64
	Altitude      float64
	Mileage       float64
	FuelLevel     float64
}

func (c *CreateVehicleCommand) CommandName() string {
	return "CreateVehicle"
}

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

type ChangeVehicleStatusCommand struct {
	VehicleID string
	NewStatus string
}

func (c *ChangeVehicleStatusCommand) CommandName() string {
	return "ChangeVehicleStatus"
}

type UpdateVehicleMileageCommand struct {
	VehicleID string
	Mileage   float64
}

func (c *UpdateVehicleMileageCommand) CommandName() string {
	return "UpdateVehicleMileage"
}

type UpdateVehicleFuelLevelCommand struct {
	VehicleID string
	FuelLevel float64
}

func (c *UpdateVehicleFuelLevelCommand) CommandName() string {
	return "UpdateVehicleFuelLevel"
}
