package command

import "context"

// Command is the base interface for all commands.
type Command interface {
	CommandName() string
}

// CommandHandler is the base interface for command handlers.
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// CommandBus defines the interface for dispatching commands.
type CommandBus interface {
	// Dispatch sends a command for processing.
	Dispatch(ctx context.Context, cmd Command) error

	// Register registers a command handler for a specific command type.
	Register(commandName string, handler CommandHandler)
}

// CreateVehicleCommand represents the command to create a new vehicle.
type CreateVehicleCommand struct {
	VIN           string
	VehicleName   string
	VehicleModel  string
	LicenseNumber string
	Status        string
	Latitude      float64
	Longitude     float64
	Altitude      float64
	Mileage       int64
	FuelLevel     int
}

// CommandName returns the command name.
func (c *CreateVehicleCommand) CommandName() string {
	return "CreateVehicle"
}

// UpdateVehicleLocationCommand represents the command to update vehicle location.
type UpdateVehicleLocationCommand struct {
	VehicleID string
	Latitude  float64
	Longitude float64
	Altitude  float64
	Timestamp int64
}

// CommandName returns the command name.
func (c *UpdateVehicleLocationCommand) CommandName() string {
	return "UpdateVehicleLocation"
}

// ChangeVehicleStatusCommand represents the command to change vehicle status.
type ChangeVehicleStatusCommand struct {
	VehicleID string
	NewStatus string
}

// CommandName returns the command name.
func (c *ChangeVehicleStatusCommand) CommandName() string {
	return "ChangeVehicleStatus"
}

// UpdateVehicleMileageCommand represents the command to update vehicle mileage.
type UpdateVehicleMileageCommand struct {
	VehicleID string
	Mileage   int64
}

// CommandName returns the command name.
func (c *UpdateVehicleMileageCommand) CommandName() string {
	return "UpdateVehicleMileage"
}

// UpdateVehicleFuelLevelCommand represents the command to update vehicle fuel level.
type UpdateVehicleFuelLevelCommand struct {
	VehicleID string
	FuelLevel int
}

// CommandName returns the command name.
func (c *UpdateVehicleFuelLevelCommand) CommandName() string {
	return "UpdateVehicleFuelLevel"
}
