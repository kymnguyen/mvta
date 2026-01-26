package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/entity"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/valueobject"
)

type CreateVehicleCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

func NewCreateVehicleCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *CreateVehicleCommandHandler {
	return &CreateVehicleCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

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

	exists, err := h.vehicleRepo.ExistsByVIN(ctx, createCmd.VIN)
	if err != nil {
		return fmt.Errorf("failed to check vin existence: %w", err)
	}
	if exists {
		return fmt.Errorf("vehicle with vin %s already exists", createCmd.VIN)
	}

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

	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}

type UpdateVehicleLocationCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

func NewUpdateVehicleLocationCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *UpdateVehicleLocationCommandHandler {
	return &UpdateVehicleLocationCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

func (h *UpdateVehicleLocationCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	updateCmd, ok := cmd.(*command.UpdateVehicleLocationCommand)
	if !ok {
		return fmt.Errorf("invalid command type for UpdateVehicleLocationCommandHandler")
	}

	vehicleID, err := valueobject.NewVehicleID(updateCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	location, err := valueobject.NewLocation(
		updateCmd.Latitude,
		updateCmd.Longitude,
		updateCmd.Altitude,
		updateCmd.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("invalid location: %w", err)
	}

	if err := vehicle.UpdateLocation(location); err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}

type UpdateVehicleMileageCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

func NewUpdateVehicleMileageCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *UpdateVehicleMileageCommandHandler {
	return &UpdateVehicleMileageCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

func (h *UpdateVehicleMileageCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	mileageCmd, ok := cmd.(*command.UpdateVehicleMileageCommand)
	if !ok {
		return fmt.Errorf("invalid command type for UpdateVehicleMileageCommandHandler")
	}

	vehicleID, err := valueobject.NewVehicleID(mileageCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	newMileage, err := valueobject.NewMileage(mileageCmd.Mileage)
	if err != nil {
		return fmt.Errorf("invalid mileage: %w", err)
	}

	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	if err := vehicle.UpdateMileage(newMileage); err != nil {
		return fmt.Errorf("failed to update mileage: %w", err)
	}

	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}

type UpdateVehicleFuelLevelCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

func NewUpdateVehicleFuelLevelCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *UpdateVehicleFuelLevelCommandHandler {
	return &UpdateVehicleFuelLevelCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

func (h *UpdateVehicleFuelLevelCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	fuelCmd, ok := cmd.(*command.UpdateVehicleFuelLevelCommand)
	if !ok {
		return fmt.Errorf("invalid command type for UpdateVehicleFuelLevelCommandHandler")
	}

	vehicleID, err := valueobject.NewVehicleID(fuelCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	newFuelLevel, err := valueobject.NewFuelLevel(fuelCmd.FuelLevel)
	if err != nil {
		return fmt.Errorf("invalid fuel level: %w", err)
	}

	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	if err := vehicle.UpdateFuelLevel(newFuelLevel); err != nil {
		return fmt.Errorf("failed to update fuel level: %w", err)
	}

	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}

type ChangeVehicleStatusCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

func NewChangeVehicleStatusCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *ChangeVehicleStatusCommandHandler {
	return &ChangeVehicleStatusCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

func (h *ChangeVehicleStatusCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	statusCmd, ok := cmd.(*command.ChangeVehicleStatusCommand)
	if !ok {
		return fmt.Errorf("invalid command type for ChangeVehicleStatusCommandHandler")
	}

	vehicleID, err := valueobject.NewVehicleID(statusCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	newStatus, err := valueobject.NewVehicleStatus(statusCmd.NewStatus)
	if err != nil {
		return fmt.Errorf("invalid status: %w", err)
	}

	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	if err := vehicle.ChangeStatus(newStatus); err != nil {
		return fmt.Errorf("failed to change status: %w", err)
	}

	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}
