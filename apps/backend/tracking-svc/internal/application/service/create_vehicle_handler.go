package service

import (
	"context"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/entity"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/valueobject"
)

type CreateVehicleCommandHandler struct {
	vehicleRepo repository.VehicleRepository
}

func NewCreateVehicleCommandHandler(vehicleRepo repository.VehicleRepository) *CreateVehicleCommandHandler {
	return &CreateVehicleCommandHandler{
		vehicleRepo: vehicleRepo,
	}
}

func (h *CreateVehicleCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	command := cmd.(*command.CreateVehicleCommand)

	exists, err := h.vehicleRepo.ExistsByVIN(ctx, command.VIN)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	vehicleID := valueobject.GenerateVehicleID()

	licenseNumber, err := valueobject.NewLicenseNumber(command.LicenseNumber)
	if err != nil {
		return err
	}

	status, err := valueobject.NewVehicleStatus(command.Status)
	if err != nil {
		return err
	}

	location, err := valueobject.NewLocation(command.Latitude, command.Longitude, command.Altitude, 0)
	if err != nil {
		return err
	}

	mileage, err := valueobject.NewMileage(command.Mileage)
	if err != nil {
		return err
	}

	fuelLevel, err := valueobject.NewFuelLevel(command.FuelLevel)
	if err != nil {
		return err
	}

	vehicle, err := entity.NewVehicle(
		vehicleID,
		command.VIN,
		command.VehicleName,
		command.VehicleModel,
		licenseNumber,
		status,
		location,
		mileage,
		fuelLevel,
	)
	if err != nil {
		return err
	}

	return h.vehicleRepo.Save(ctx, vehicle)
}
