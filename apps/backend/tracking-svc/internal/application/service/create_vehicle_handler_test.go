package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/entity"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/valueobject"
)

// --- Mocks ---

type MockVehicleRepo struct{ mock.Mock }

func (m *MockVehicleRepo) Save(ctx context.Context, v *entity.Vehicle) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}
func (m *MockVehicleRepo) FindByID(ctx context.Context, id valueobject.VehicleID) (*entity.Vehicle, error) {
	return nil, nil
}
func (m *MockVehicleRepo) FindAll(ctx context.Context, limit int, offset int) ([]*entity.Vehicle, error) {
	return nil, nil
}
func (m *MockVehicleRepo) Delete(ctx context.Context, id valueobject.VehicleID) error {
	return nil
}
func (m *MockVehicleRepo) ExistsByVIN(ctx context.Context, vin string) (bool, error) {
	args := m.Called(ctx, vin)
	return args.Bool(0), args.Error(1)
}

type MockOutboxRepo struct{ mock.Mock }

func (m *MockOutboxRepo) SaveOutboxEvent(ctx context.Context, aggregateID string, event interface{}) error {
	args := m.Called(ctx, aggregateID, event)
	return args.Error(0)
}
func (m *MockOutboxRepo) GetPendingEvents(ctx context.Context, limit int) ([]repository.OutboxEvent, error) {
	return nil, nil
}
func (m *MockOutboxRepo) MarkEventAsPublished(ctx context.Context, eventID string) error {
	return nil
}

// --- Tests ---
func TestCreateVehicleCommandHandler_Success(t *testing.T) {
	vehicleRepo := new(MockVehicleRepo)
	outboxRepo := new(MockOutboxRepo)
	h := NewCreateVehicleCommandHandler(vehicleRepo, outboxRepo)

	cmd := &command.CreateVehicleCommand{
		VIN:           "VIN123",
		VehicleName:   "TestCar",
		VehicleModel:  "ModelX",
		LicenseNumber: "ABC123",
		Status:        "active",
		Latitude:      1.23,
		Longitude:     4.56,
		Altitude:      7.89,
		Mileage:       1000,
		FuelLevel:     80,
	}

	vehicleRepo.On("ExistsByVIN", mock.Anything, "VIN123").Return(false, nil)
	vehicleRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	outboxRepo.On("SaveOutboxEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	vehicleRepo.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestCreateVehicleCommandHandler_AlreadyExists(t *testing.T) {
	vehicleRepo := new(MockVehicleRepo)
	outboxRepo := new(MockOutboxRepo)
	h := NewCreateVehicleCommandHandler(vehicleRepo, outboxRepo)

	cmd := &command.CreateVehicleCommand{
		VIN:           "VIN123",
		VehicleName:   "TestCar",
		VehicleModel:  "ModelX",
		LicenseNumber: "ABC123",
		Status:        "active",
		Latitude:      1.23,
		Longitude:     4.56,
		Altitude:      7.89,
		Mileage:       1000,
		FuelLevel:     80,
	}

	vehicleRepo.On("ExistsByVIN", mock.Anything, "VIN123").Return(false, nil)
	vehicleRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))

	err := h.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save vehicle")
	vehicleRepo.AssertExpectations(t)
	outboxRepo.AssertNotCalled(t, "SaveOutboxEvent", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateVehicleCommandHandler_InvalidStatus(t *testing.T) {
	vehicleRepo := new(MockVehicleRepo)
	outboxRepo := new(MockOutboxRepo)
	h := NewCreateVehicleCommandHandler(vehicleRepo, outboxRepo)

	cmd := &command.CreateVehicleCommand{
		Status: "badstatus",
	}
	vehicleRepo.On("ExistsByVIN", mock.Anything, "").Return(false, nil)

	err := h.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
}

func TestCreateVehicleCommandHandler_SaveError(t *testing.T) {
	vehicleRepo := new(MockVehicleRepo)
	outboxRepo := new(MockOutboxRepo)
	h := NewCreateVehicleCommandHandler(vehicleRepo, outboxRepo)

	cmd := &command.CreateVehicleCommand{
		VIN:           "VIN123",
		VehicleName:   "TestCar",
		VehicleModel:  "ModelX",
		LicenseNumber: "ABC123",
		Status:        "active",
		Latitude:      1.23,
		Longitude:     4.56,
		Altitude:      7.89,
		Mileage:       1000,
		FuelLevel:     80,
	}

	vehicleRepo.On("ExistsByVIN", mock.Anything, "VIN123").Return(false, nil)
	vehicleRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))

	err := h.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save vehicle")
	vehicleRepo.AssertExpectations(t)
}
