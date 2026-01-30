import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { trackingVehicleApi, type Vehicle } from '../../shared/api/tracking-svc';
import { Link } from 'react-router-dom';

export const TrackingVehicleList: React.FC = () => {
  const { data: vehicles, isLoading, error } = useQuery({
    queryKey: ['vehicles'],
    queryFn: trackingVehicleApi.getVehicles,
  });

  if (isLoading) return <div className="loading">Loading vehicles...</div>;
  if (error) return <div className="error">Error loading vehicles</div>;

  return (
    <div className="vehicle-list">
      <h1>Vehicle Fleet</h1>
      <div className="vehicle-grid">
        {vehicles?.map((vehicle: Vehicle) => (
          <Link
            key={vehicle.id}
            to={`/vehicles/${vehicle.id}`}
            className="vehicle-card"
          >
            <div className="vehicle-header">
              <h3>{vehicle.vehicleName}</h3>
              <span className={`status status-${vehicle.status}`}>
                {vehicle.status}
              </span>
            </div>
            <div className="vehicle-details">
              <p>
                <strong>VIN:</strong> {vehicle.vin}
              </p>
              <p>
                <strong>Model:</strong> {vehicle.vehicleModel}
              </p>
              <p>
                <strong>License:</strong> {vehicle.licenseNumber}
              </p>
              <p>
                <strong>Mileage:</strong> {vehicle.mileage.toFixed(1)} km
              </p>
              <p>
                <strong>Fuel:</strong> {vehicle.fuelLevel.toFixed(1)}%
              </p>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
};
