package vehicle

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"context"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/dto"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/query"
)

// --- Mocks ---

type MockCommandBus struct{ mock.Mock }

func (m *MockCommandBus) Dispatch(ctx context.Context, cmd command.Command) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

func (m *MockCommandBus) Register(commandName string, handler command.CommandHandler) {
	// no-op for test
}

type MockQueryBus struct{ mock.Mock }

func (m *MockQueryBus) Dispatch(ctx context.Context, q query.Query) (query.QueryResult, error) {
	args := m.Called(ctx, q)
	return args.Get(0), args.Error(1)
}

func (m *MockQueryBus) Register(queryName string, handler query.QueryHandler) {
	// no-op for test
}

// --- Tests ---
func TestUpdateLocation_Success(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	reqBody := dto.UpdateVehicleLocationRequest{
		Latitude: 1.1, Longitude: 2.2, Altitude: 3.3, Timestamp: 123456,
	}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/location", bytes.NewReader(b))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	cmdBus.On("Dispatch", mock.Anything, mock.AnythingOfType("*command.UpdateVehicleLocationCommand")).Return(nil)
	h.UpdateLocation(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateLocation_BadRequest(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/location", bytes.NewReader([]byte("bad json")))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	h.UpdateLocation(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateLocation_CommandError(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	reqBody := dto.UpdateVehicleLocationRequest{Latitude: 1.1}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/location", bytes.NewReader(b))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	cmdBus.On("Dispatch", mock.Anything, mock.AnythingOfType("*command.UpdateVehicleLocationCommand")).Return(errors.New("fail"))
	h.UpdateLocation(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUpdateLocation_BadID(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles//location", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()
	h.UpdateLocation(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateMileage_Success(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	reqBody := dto.UpdateVehicleMileageRequest{Mileage: 1234.5}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/mileage", bytes.NewReader(b))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	cmdBus.On("Dispatch", mock.Anything, mock.AnythingOfType("*command.UpdateVehicleMileageCommand")).Return(nil)
	h.UpdateMileage(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateMileage_BadRequest(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/mileage", bytes.NewReader([]byte("bad json")))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	h.UpdateMileage(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateMileage_CommandError(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	reqBody := dto.UpdateVehicleMileageRequest{Mileage: 1234.5}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/mileage", bytes.NewReader(b))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	cmdBus.On("Dispatch", mock.Anything, mock.AnythingOfType("*command.UpdateVehicleMileageCommand")).Return(errors.New("fail"))
	h.UpdateMileage(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUpdateMileage_BadID(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles//mileage", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()
	h.UpdateMileage(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateFuelLevel_Success(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	reqBody := dto.UpdateVehicleFuelLevelRequest{FuelLevel: 55.5}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/fuel", bytes.NewReader(b))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	cmdBus.On("Dispatch", mock.Anything, mock.AnythingOfType("*command.UpdateVehicleFuelLevelCommand")).Return(nil)
	h.UpdateFuelLevel(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateFuelLevel_BadRequest(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/fuel", bytes.NewReader([]byte("bad json")))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	h.UpdateFuelLevel(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateFuelLevel_CommandError(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	reqBody := dto.UpdateVehicleFuelLevelRequest{FuelLevel: 55.5}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/fuel", bytes.NewReader(b))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	cmdBus.On("Dispatch", mock.Anything, mock.AnythingOfType("*command.UpdateVehicleFuelLevelCommand")).Return(errors.New("fail"))
	h.UpdateFuelLevel(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUpdateFuelLevel_BadID(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles//fuel", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()
	h.UpdateFuelLevel(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestChangeStatus_Success(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	reqBody := dto.ChangeVehicleStatusRequest{Status: "inactive"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/status", bytes.NewReader(b))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	cmdBus.On("Dispatch", mock.Anything, mock.AnythingOfType("*command.ChangeVehicleStatusCommand")).Return(nil)
	h.ChangeStatus(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestChangeStatus_BadRequest(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/status", bytes.NewReader([]byte("bad json")))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	h.ChangeStatus(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestChangeStatus_CommandError(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	reqBody := dto.ChangeVehicleStatusRequest{Status: "inactive"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/VIN123/status", bytes.NewReader(b))
	req.SetPathValue("id", "VIN123")
	w := httptest.NewRecorder()
	cmdBus.On("Dispatch", mock.Anything, mock.AnythingOfType("*command.ChangeVehicleStatusCommand")).Return(errors.New("fail"))
	h.ChangeStatus(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestChangeStatus_BadID(t *testing.T) {
	cmdBus := new(MockCommandBus)
	qryBus := new(MockQueryBus)
	h := &VehicleHandler{commandBus: cmdBus, queryBus: qryBus, logger: zap.NewNop()}
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles//status", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()
	h.ChangeStatus(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
