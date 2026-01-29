package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/entity"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/valueobject"
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
