package handler

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/services/vehicle/internal/application/command"
	"github.com/kymnguyen/mvta/services/vehicle/internal/application/dto"
	"github.com/kymnguyen/mvta/services/vehicle/internal/application/query"
)

type VehicleHandler struct {
	commandBus command.CommandBus
	queryBus   query.QueryBus
	logger     *zap.Logger
}

func NewVehicleHandler(
	commandBus command.CommandBus,
	queryBus query.QueryBus,
	logger *zap.Logger,
) *VehicleHandler {
	return &VehicleHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
		logger:     logger,
	}
}

func (h *VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateVehicleRequest
	if err := decodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode create vehicle request", zap.Error(err))
		respondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.CreateVehicleCommand{
		VIN:           req.VIN,
		VehicleName:   req.VehicleName,
		VehicleModel:  req.VehicleModel,
		LicenseNumber: req.LicenseNumber,
		Status:        req.Status,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Altitude:      req.Altitude,
		Mileage:       req.Mileage,
		FuelLevel:     req.FuelLevel,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to create vehicle", zap.Error(err))
		respondError(w, http.StatusInternalServerError, "ERR_CREATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle created successfully", zap.String("vin", req.VIN))
	respondSuccess(w, http.StatusCreated, map[string]string{
		"message": "vehicle created successfully",
	})
}

func (h *VehicleHandler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		respondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	q := &query.GetVehicleQuery{
		VehicleID: vehicleID,
	}

	result, err := h.queryBus.Dispatch(ctx, q)
	if err != nil {
		h.logger.Error("failed to get vehicle", zap.String("vehicleId", vehicleID), zap.Error(err))
		respondError(w, http.StatusNotFound, "ERR_NOT_FOUND", "Vehicle not found")
		return
	}

	respondSuccess(w, http.StatusOK, result)
}

func (h *VehicleHandler) GetAllVehicles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if _, err := scanInt(l, &limit); err != nil {
			respondError(w, http.StatusBadRequest, "ERR_INVALID_LIMIT", "Invalid limit parameter")
			return
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if _, err := scanInt(o, &offset); err != nil {
			respondError(w, http.StatusBadRequest, "ERR_INVALID_OFFSET", "Invalid offset parameter")
			return
		}
	}

	q := &query.GetAllVehiclesQuery{
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.queryBus.Dispatch(ctx, q)
	if err != nil {
		h.logger.Error("failed to get all vehicles", zap.Error(err))
		respondError(w, http.StatusInternalServerError, "ERR_QUERY_FAILED", err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, map[string]interface{}{
		"vehicles": result,
		"limit":    limit,
		"offset":   offset,
	})
}

// UpdateLocation handles PATCH /vehicles/{id}/location - updates vehicle location.
func (h *VehicleHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		respondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	var req dto.UpdateVehicleLocationRequest
	if err := decodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode update location request", zap.Error(err))
		respondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.UpdateVehicleLocationCommand{
		VehicleID: vehicleID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Altitude:  req.Altitude,
		Timestamp: req.Timestamp,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to update vehicle location",
			zap.String("vehicleId", vehicleID),
			zap.Error(err))
		respondError(w, http.StatusInternalServerError, "ERR_UPDATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle location updated", zap.String("vehicleId", vehicleID))
	respondSuccess(w, http.StatusOK, map[string]string{
		"message": "location updated successfully",
	})
}

func (h *VehicleHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		respondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	var req dto.ChangeVehicleStatusRequest
	if err := decodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode change status request", zap.Error(err))
		respondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.ChangeVehicleStatusCommand{
		VehicleID: vehicleID,
		NewStatus: req.Status,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to change vehicle status",
			zap.String("vehicleId", vehicleID),
			zap.Error(err))
		respondError(w, http.StatusInternalServerError, "ERR_UPDATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle status changed", zap.String("vehicleId", vehicleID))
	respondSuccess(w, http.StatusOK, map[string]string{
		"message": "status changed successfully",
	})
}

func (h *VehicleHandler) UpdateMileage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		respondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	var req dto.UpdateVehicleMileageRequest
	if err := decodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode update mileage request", zap.Error(err))
		respondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.UpdateVehicleMileageCommand{
		VehicleID: vehicleID,
		Mileage:   req.Mileage,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to update vehicle mileage",
			zap.String("vehicleId", vehicleID),
			zap.Error(err))
		respondError(w, http.StatusInternalServerError, "ERR_UPDATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle mileage updated", zap.String("vehicleId", vehicleID), zap.Float64("mileage", req.Mileage))
	respondSuccess(w, http.StatusOK, map[string]string{
		"message": "mileage updated successfully",
	})
}

func (h *VehicleHandler) UpdateFuelLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		respondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	var req dto.UpdateVehicleFuelLevelRequest
	if err := decodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode update fuel level request", zap.Error(err))
		respondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.UpdateVehicleFuelLevelCommand{
		VehicleID: vehicleID,
		FuelLevel: req.FuelLevel,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to update vehicle fuel level",
			zap.String("vehicleId", vehicleID),
			zap.Error(err))
		respondError(w, http.StatusInternalServerError, "ERR_UPDATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle fuel level updated", zap.String("vehicleId", vehicleID), zap.Float64("fuelLevel", req.FuelLevel))
	respondSuccess(w, http.StatusOK, map[string]string{
		"message": "fuel level updated successfully",
	})
}
