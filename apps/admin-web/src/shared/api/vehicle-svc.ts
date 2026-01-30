import axios from 'axios';

const vehicleSvcClient = axios.create({
  baseURL: '/api/vehicle',
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface VehicleSvc {
  id: string;
  vin: string;
  vehicleName: string;
  vehicleModel: string;
  licenseNumber: string;
  status: string;
  latitude: number;
  longitude: number;
  altitude: number;
  mileage: number;
  fuelLevel: number;
  version: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateVehicleRequest {
  vin: string;
  vehicleName: string;
  vehicleModel: string;
  licenseNumber: string;
  status: string;
  latitude: number;
  longitude: number;
  altitude?: number;
  mileage?: number;
  fuelLevel?: number;
}

export interface UpdateVehicleLocationRequest {
  latitude: number;
  longitude: number;
  altitude?: number;
  timestamp: number;
}

export interface UpdateVehicleStatusRequest {
  status: string;
}

export interface UpdateVehicleMileageRequest {
  mileage: number;
}

export interface UpdateVehicleFuelRequest {
  fuelLevel: number;
}

export const vehicleSvcApi = {
  // Get all vehicles
  getVehicles: async (): Promise<VehicleSvc[]> => {
    const response = await vehicleSvcClient.get('/vehicles');
    return response.data.vehicles || [];
  },

  // Get vehicle by ID
  getVehicle: async (id: string): Promise<VehicleSvc> => {
    const response = await vehicleSvcClient.get(`/vehicles/${id}`);
    return response.data;
  },

  // Create vehicle
  createVehicle: async (data: CreateVehicleRequest): Promise<VehicleSvc> => {
    const response = await vehicleSvcClient.post('/vehicles', data);
    return response.data;
  },

  // Update vehicle location
  updateLocation: async (
    id: string,
    data: UpdateVehicleLocationRequest
  ): Promise<VehicleSvc> => {
    const response = await vehicleSvcClient.patch(`/vehicles/${id}/location`, data);
    return response.data;
  },

  // Update vehicle status
  updateStatus: async (
    id: string,
    data: UpdateVehicleStatusRequest
  ): Promise<VehicleSvc> => {
    const response = await vehicleSvcClient.patch(`/vehicles/${id}/status`, data);
    return response.data;
  },

  // Update vehicle mileage
  updateMileage: async (
    id: string,
    data: UpdateVehicleMileageRequest
  ): Promise<VehicleSvc> => {
    const response = await vehicleSvcClient.patch(`/vehicles/${id}/mileage`, data);
    return response.data;
  },

  // Update vehicle fuel level
  updateFuel: async (
    id: string,
    data: UpdateVehicleFuelRequest
  ): Promise<VehicleSvc> => {
    const response = await vehicleSvcClient.patch(`/vehicles/${id}/fuel`, data);
    return response.data;
  },
};
