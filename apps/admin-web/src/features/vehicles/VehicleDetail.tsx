import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { useParams, Link } from 'react-router-dom';
import { vehicleApi } from '../../shared/api/client';
import { format } from 'date-fns';

export const VehicleDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  const { data: vehicle, isLoading: vehicleLoading } = useQuery({
    queryKey: ['vehicle', id],
    queryFn: () => vehicleApi.getVehicle(id!),
    enabled: !!id,
  });

  const { data: history, isLoading: historyLoading } = useQuery({
    queryKey: ['vehicleHistory', id],
    queryFn: () => vehicleApi.getVehicleHistory(id!),
    enabled: !!id,
  });

  if (vehicleLoading || historyLoading) {
    return <div className="loading">Loading vehicle details...</div>;
  }

  if (!vehicle) {
    return <div className="error">Vehicle not found</div>;
  }

  const getChangeTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
      created: '🆕 Created',
      location_updated: '📍 Location Updated',
      status_changed: '🔄 Status Changed',
      mileage_updated: '🚗 Mileage Updated',
      fuel_updated: '⛽ Fuel Updated',
    };
    return labels[type] || type;
  };

  const formatValue = (value: any): string => {
    if (typeof value === 'object' && value !== null) {
      return JSON.stringify(value, null, 2);
    }
    return String(value);
  };

  return (
    <div className="vehicle-detail">
      <Link to="/" className="back-link">
        ← Back to Fleet
      </Link>

      <div className="vehicle-info-card">
        <div className="vehicle-header">
          <h1>{vehicle.vehicleName}</h1>
          <span className={`status status-${vehicle.status}`}>
            {vehicle.status}
          </span>
        </div>

        <div className="info-grid">
          <div className="info-item">
            <label>VIN</label>
            <span>{vehicle.vin}</span>
          </div>
          <div className="info-item">
            <label>Model</label>
            <span>{vehicle.vehicleModel}</span>
          </div>
          <div className="info-item">
            <label>License Number</label>
            <span>{vehicle.licenseNumber}</span>
          </div>
          <div className="info-item">
            <label>Mileage</label>
            <span>{vehicle.mileage.toFixed(1)} km</span>
          </div>
          <div className="info-item">
            <label>Fuel Level</label>
            <span>{vehicle.fuelLevel.toFixed(1)}%</span>
          </div>
          <div className="info-item">
            <label>Location</label>
            <span>
              {vehicle.latitude.toFixed(6)}, {vehicle.longitude.toFixed(6)}
            </span>
          </div>
        </div>
      </div>

      <div className="history-section">
        <h2>Change History</h2>
        {history && history.changes.length > 0 ? (
          <div className="history-timeline">
            {history.changes.map((change, index) => (
              <div key={index} className="history-item">
                <div className="history-marker" />
                <div className="history-content">
                  <div className="history-header">
                    <span className="change-type">
                      {getChangeTypeLabel(change.changeType)}
                    </span>
                    <span className="change-time">
                      {format(new Date(change.changedAt), 'PPpp')}
                    </span>
                  </div>

                  <div className="change-details">
                    {Object.keys(change.newValue).length > 0 && (
                      <div className="change-values">
                        <h4>Changes:</h4>
                        {Object.entries(change.newValue).map(([key, value]) => (
                          <div key={key} className="change-value-item">
                            <strong>{key}:</strong>
                            <div className="value-comparison">
                              {change.oldValue[key] !== undefined && (
                                <span className="old-value">
                                  {formatValue(change.oldValue[key])}
                                </span>
                              )}
                              <span className="arrow">→</span>
                              <span className="new-value">
                                {formatValue(value)}
                              </span>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                    <div className="version-info">Version: {change.version}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="no-history">No change history available</p>
        )}
      </div>
    </div>
  );
};
