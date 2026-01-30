import axios from 'axios';

// API clients for different backend services
const trackingApi = axios.create({
  baseURL: '/api/tracking',
  headers: {
    'Content-Type': 'application/json',
  },
});


export interface Vehicle {
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
  createdAt: string;
  updatedAt: string;
}

export interface VehicleChangeRecord {
  vehicleId: string;
  vin: string;
  changeType: string;
  oldValue: Record<string, any>;
  newValue: Record<string, any>;
  changedAt: string;
  version: number;
}

export interface VehicleChangeHistoryResponse {
  vehicleId: string;
  changes: VehicleChangeRecord[];
  total: number;
}

export const trackingVehicleApi = {
  getVehicles: async (): Promise<Vehicle[]> => {
    const response = await trackingApi.get('/vehicles');
    return response.data.vehicles;
  },

  getVehicle: async (id: string): Promise<Vehicle> => {
    const response = await trackingApi.get(`/vehicles/${id}`);
    return response.data;
  },

  getVehicleHistory: async (
    id: string,
    limit = 50,
    offset = 0
  ): Promise<VehicleChangeHistoryResponse> => {
    const response = await trackingApi.get(`/vehicles/${id}/history`, {
      params: { limit, offset },
    });
    return response.data;
  },
};

export default trackingVehicleApi;
