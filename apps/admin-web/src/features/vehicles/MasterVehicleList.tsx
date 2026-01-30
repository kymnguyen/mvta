import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { vehicleSvcApi, type VehicleSvc } from '../../shared/api/vehicle-svc';

export const MasterVehicleList: React.FC = () => {
  const queryClient = useQueryClient();
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [editingVehicle, setEditingVehicle] = useState<VehicleSvc | null>(null);

  const { data: vehicles, isLoading, error } = useQuery({
    queryKey: ['vehicleSvc'],
    queryFn: vehicleSvcApi.getVehicles,
  });

  const createMutation = useMutation({
    mutationFn: vehicleSvcApi.createVehicle,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vehicleSvc'] });
      setShowCreateForm(false);
    },
  });

  const updateLocationMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) =>
      vehicleSvcApi.updateLocation(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vehicleSvc'] });
      setEditingVehicle(null);
    },
  });

  const updateStatusMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) =>
      vehicleSvcApi.updateStatus(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vehicleSvc'] });
    },
  });

  const handleCreateSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const data = {
      vin: formData.get('vin') as string,
      vehicleName: formData.get('vehicleName') as string,
      vehicleModel: formData.get('vehicleModel') as string,
      licenseNumber: formData.get('licenseNumber') as string,
      status: formData.get('status') as string,
      latitude: parseFloat(formData.get('latitude') as string),
      longitude: parseFloat(formData.get('longitude') as string),
      altitude: parseFloat(formData.get('altitude') as string) || 0,
      mileage: parseFloat(formData.get('mileage') as string) || 0,
      fuelLevel: parseFloat(formData.get('fuelLevel') as string) || 100,
    };
    createMutation.mutate(data);
  };

  const handleUpdateLocation = (vehicle: VehicleSvc, e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const data = {
      latitude: parseFloat(formData.get('latitude') as string),
      longitude: parseFloat(formData.get('longitude') as string),
      altitude: parseFloat(formData.get('altitude') as string) || 0,
      timestamp: Date.now(),
    };
    updateLocationMutation.mutate({ id: vehicle.id, data });
  };

  const handleStatusChange = (id: string, status: string) => {
    updateStatusMutation.mutate({ id, data: { status } });
  };

  if (isLoading) return <div className="loading">Loading vehicles...</div>;
  if (error) return <div className="error">Error loading vehicles</div>;

  return (
    <div className="vehicle-svc-list">
      <div className="header">
        <h1>Vehicle Service (vehicle-svc)</h1>
        <button
          className="btn btn-primary"
          onClick={() => setShowCreateForm(!showCreateForm)}
        >
          {showCreateForm ? 'Cancel' : '+ Create Vehicle'}
        </button>
      </div>

      {showCreateForm && (
        <div className="create-form-container">
          <h2>Create New Vehicle</h2>
          <form onSubmit={handleCreateSubmit} className="vehicle-form">
            <div className="form-grid">
              <div className="form-group">
                <label>VIN *</label>
                <input type="text" name="vin" required />
              </div>
              <div className="form-group">
                <label>Vehicle Name *</label>
                <input type="text" name="vehicleName" required />
              </div>
              <div className="form-group">
                <label>Model *</label>
                <input type="text" name="vehicleModel" required />
              </div>
              <div className="form-group">
                <label>License Number *</label>
                <input type="text" name="licenseNumber" required />
              </div>
              <div className="form-group">
                <label>Status *</label>
                <select name="status" required>
                  <option value="available">Available</option>
                  <option value="in_use">In Use</option>
                  <option value="maintenance">Maintenance</option>
                  <option value="out_of_service">Out of Service</option>
                </select>
              </div>
              <div className="form-group">
                <label>Latitude *</label>
                <input type="number" step="any" name="latitude" required />
              </div>
              <div className="form-group">
                <label>Longitude *</label>
                <input type="number" step="any" name="longitude" required />
              </div>
              <div className="form-group">
                <label>Altitude</label>
                <input type="number" step="any" name="altitude" defaultValue="0" />
              </div>
              <div className="form-group">
                <label>Mileage (km)</label>
                <input type="number" step="any" name="mileage" defaultValue="0" />
              </div>
              <div className="form-group">
                <label>Fuel Level (%)</label>
                <input
                  type="number"
                  step="any"
                  name="fuelLevel"
                  defaultValue="100"
                  min="0"
                  max="100"
                />
              </div>
            </div>
            <div className="form-actions">
              <button type="submit" className="btn btn-primary" disabled={createMutation.isPending}>
                {createMutation.isPending ? 'Creating...' : 'Create Vehicle'}
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="vehicles-table">
        <table>
          <thead>
            <tr>
              <th>VIN</th>
              <th>Name</th>
              <th>Model</th>
              <th>License</th>
              <th>Status</th>
              <th>Location</th>
              <th>Mileage</th>
              <th>Fuel</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {vehicles?.map((vehicle: VehicleSvc) => (
              <tr key={vehicle.id}>
                <td>{vehicle.vin}</td>
                <td>{vehicle.vehicleName}</td>
                <td>{vehicle.vehicleModel}</td>
                <td>{vehicle.licenseNumber}</td>
                <td>
                  <select
                    value={vehicle.status}
                    onChange={(e) => handleStatusChange(vehicle.id, e.target.value)}
                    className={`status-select status-${vehicle.status}`}
                  >
                    <option value="available">Available</option>
                    <option value="in_use">In Use</option>
                    <option value="maintenance">Maintenance</option>
                    <option value="out_of_service">Out of Service</option>
                  </select>
                </td>
                <td>
                  {vehicle.latitude.toFixed(4)}, {vehicle.longitude.toFixed(4)}
                </td>
                <td>{vehicle.mileage.toFixed(1)} km</td>
                <td>{vehicle.fuelLevel.toFixed(1)}%</td>
                <td>
                  <button
                    className="btn btn-small"
                    onClick={() => setEditingVehicle(vehicle)}
                  >
                    Update Location
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {editingVehicle && (
        <div className="modal-overlay" onClick={() => setEditingVehicle(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h2>Update Location - {editingVehicle.vehicleName}</h2>
            <form onSubmit={(e) => handleUpdateLocation(editingVehicle, e)} className="vehicle-form">
              <div className="form-group">
                <label>Latitude *</label>
                <input
                  type="number"
                  step="any"
                  name="latitude"
                  defaultValue={editingVehicle.latitude}
                  required
                />
              </div>
              <div className="form-group">
                <label>Longitude *</label>
                <input
                  type="number"
                  step="any"
                  name="longitude"
                  defaultValue={editingVehicle.longitude}
                  required
                />
              </div>
              <div className="form-group">
                <label>Altitude</label>
                <input
                  type="number"
                  step="any"
                  name="altitude"
                  defaultValue={editingVehicle.altitude}
                />
              </div>
              <div className="form-actions">
                <button
                  type="button"
                  className="btn btn-secondary"
                  onClick={() => setEditingVehicle(null)}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="btn btn-primary"
                  disabled={updateLocationMutation.isPending}
                >
                  {updateLocationMutation.isPending ? 'Updating...' : 'Update'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
