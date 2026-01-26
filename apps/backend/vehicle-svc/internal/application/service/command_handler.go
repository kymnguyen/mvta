package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/entity"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/valueobject"
)

// CreateVehicleCommandHandler handles vehicle creation commands.
type CreateVehicleCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

// NewCreateVehicleCommandHandler creates a new command handler.
func NewCreateVehicleCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *CreateVehicleCommandHandler {
	return &CreateVehicleCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

// Handle processes the create vehicle command.
func (h *CreateVehicleCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	createCmd, ok := cmd.(*command.CreateVehicleCommand)
	if !ok {
		return fmt.Errorf("invalid command type for CreateVehicleCommandHandler")
	}

	// Validate input
	status, err := valueobject.NewVehicleStatus(createCmd.Status)
	if err != nil {
		return fmt.Errorf("invalid status: %w", err)
	}

	location, err := valueobject.NewLocation(
		createCmd.Latitude,
		createCmd.Longitude,
		createCmd.Altitude,
		0,
	)
	if err != nil {
		return fmt.Errorf("invalid location: %w", err)
	}

	licenseNumber, err := valueobject.NewLicenseNumber(createCmd.LicenseNumber)
	if err != nil {
		return fmt.Errorf("invalid license number: %w", err)
	}

	mileage, err := valueobject.NewMileage(createCmd.Mileage)
	if err != nil {
		return fmt.Errorf("invalid mileage: %w", err)
	}

	fuelLevel, err := valueobject.NewFuelLevel(createCmd.FuelLevel)
	if err != nil {
		return fmt.Errorf("invalid fuel level: %w", err)
	}

	// Check if VIN already exists
	exists, err := h.vehicleRepo.ExistsByVIN(ctx, createCmd.VIN)
	if err != nil {
		return fmt.Errorf("failed to check vin existence: %w", err)
	}
	if exists {
		return fmt.Errorf("vehicle with vin %s already exists", createCmd.VIN)
	}

	// Create new vehicle aggregate
	vehicleID := valueobject.GenerateVehicleID()
	vehicle, err := entity.NewVehicle(
		vehicleID,
		createCmd.VIN,
		createCmd.VehicleName,
		createCmd.VehicleModel,
		licenseNumber,
		status,
		location,
		mileage,
		fuelLevel,
	)
	if err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}

	// Save vehicle and publish events via outbox pattern
	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	// Publish uncommitted domain events to outbox for asynchronous processing
	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}

// UpdateVehicleLocationCommandHandler handles vehicle location update commands.
type UpdateVehicleLocationCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

// NewUpdateVehicleLocationCommandHandler creates a new command handler.
func NewUpdateVehicleLocationCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *UpdateVehicleLocationCommandHandler {
	return &UpdateVehicleLocationCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

// Handle processes the update location command.
func (h *UpdateVehicleLocationCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	updateCmd, ok := cmd.(*command.UpdateVehicleLocationCommand)
	if !ok {
		return fmt.Errorf("invalid command type for UpdateVehicleLocationCommandHandler")
	}

	// Parse vehicle ID
	vehicleID, err := valueobject.NewVehicleID(updateCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	// Retrieve existing vehicle
	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	// Create new location
	location, err := valueobject.NewLocation(
		updateCmd.Latitude,
		updateCmd.Longitude,
		updateCmd.Altitude,
		updateCmd.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("invalid location: %w", err)
	}

	// Update location on aggregate
	if err := vehicle.UpdateLocation(location); err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	// Save updated vehicle
	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	// Publish events via outbox
	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}

// UpdateVehicleMileageCommandHandler handles vehicle mileage update commands.
type UpdateVehicleMileageCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

// NewUpdateVehicleMileageCommandHandler creates a new command handler.
func NewUpdateVehicleMileageCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *UpdateVehicleMileageCommandHandler {
	return &UpdateVehicleMileageCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

// Handle processes the update mileage command.
func (h *UpdateVehicleMileageCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	mileageCmd, ok := cmd.(*command.UpdateVehicleMileageCommand)
	if !ok {
		return fmt.Errorf("invalid command type for UpdateVehicleMileageCommandHandler")
	}

	// Parse vehicle ID
	vehicleID, err := valueobject.NewVehicleID(mileageCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	// Create new mileage
	newMileage, err := valueobject.NewMileage(mileageCmd.Mileage)
	if err != nil {
		return fmt.Errorf("invalid mileage: %w", err)
	}

	// Retrieve existing vehicle
	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	// Update mileage on aggregate
	if err := vehicle.UpdateMileage(newMileage); err != nil {
		return fmt.Errorf("failed to update mileage: %w", err)
	}

	// Save updated vehicle
	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	// Publish events via outbox
	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}

// UpdateVehicleFuelLevelCommandHandler handles vehicle fuel level update commands.
type UpdateVehicleFuelLevelCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

// NewUpdateVehicleFuelLevelCommandHandler creates a new command handler.
func NewUpdateVehicleFuelLevelCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *UpdateVehicleFuelLevelCommandHandler {
	return &UpdateVehicleFuelLevelCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

// Handle processes the update fuel level command.
func (h *UpdateVehicleFuelLevelCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	fuelCmd, ok := cmd.(*command.UpdateVehicleFuelLevelCommand)
	if !ok {
		return fmt.Errorf("invalid command type for UpdateVehicleFuelLevelCommandHandler")
	}

	// Parse vehicle ID
	vehicleID, err := valueobject.NewVehicleID(fuelCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	// Create new fuel level
	newFuelLevel, err := valueobject.NewFuelLevel(fuelCmd.FuelLevel)
	if err != nil {
		return fmt.Errorf("invalid fuel level: %w", err)
	}

	// Retrieve existing vehicle
	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	// Update fuel level on aggregate
	if err := vehicle.UpdateFuelLevel(newFuelLevel); err != nil {
		return fmt.Errorf("failed to update fuel level: %w", err)
	}

	// Save updated vehicle
	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	// Publish events via outbox
	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}

// ChangeVehicleStatusCommandHandler handles vehicle status change commands.
type ChangeVehicleStatusCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

// NewChangeVehicleStatusCommandHandler creates a new command handler.
func NewChangeVehicleStatusCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *ChangeVehicleStatusCommandHandler {
	return &ChangeVehicleStatusCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

// Handle processes the change status command.
func (h *ChangeVehicleStatusCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	statusCmd, ok := cmd.(*command.ChangeVehicleStatusCommand)
	if !ok {
		return fmt.Errorf("invalid command type for ChangeVehicleStatusCommandHandler")
	}

	// Parse vehicle ID
	vehicleID, err := valueobject.NewVehicleID(statusCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	// Validate new status
	newStatus, err := valueobject.NewVehicleStatus(statusCmd.NewStatus)
	if err != nil {
		return fmt.Errorf("invalid status: %w", err)
	}

	// Retrieve existing vehicle
	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	// Change status on aggregate
	if err := vehicle.ChangeStatus(newStatus); err != nil {
		return fmt.Errorf("failed to change status: %w", err)
	}

	// Save updated vehicle
	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	// Publish events via outbox
	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}
