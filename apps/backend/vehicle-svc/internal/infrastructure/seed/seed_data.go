package seed

import (
	"context"
	"time"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/entity"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/valueobject"
	"go.uber.org/zap"
)

type Seeder struct {
	vehicleRepo repository.VehicleRepository
	logger      *zap.Logger
}

func NewSeeder(vehicleRepo repository.VehicleRepository, logger *zap.Logger) *Seeder {
	return &Seeder{
		vehicleRepo: vehicleRepo,
		logger:      logger,
	}
}

func (s *Seeder) SeedVehicles(ctx context.Context) error {
	// Check if vehicles already exist
	vehicles, err := s.vehicleRepo.FindAll(ctx, 1, 0)
	if err != nil {
		s.logger.Error("failed to check existing vehicles", zap.Error(err))
		return err
	}

	if len(vehicles) > 0 {
		s.logger.Info("vehicles already exist, skipping seed", zap.Int("count", len(vehicles)))
		return nil
	}

	s.logger.Info("seeding initial vehicles data...")

	seedVehicles := []struct {
		vin           string
		vehicleName   string
		vehicleModel  string
		licenseNumber string
		status        string
		latitude      float64
		longitude     float64
		altitude      float64
		mileage       float64
		fuelLevel     float64
	}{
		{
			vin:           "1HGBH41JXMN109186",
			vehicleName:   "Fleet Vehicle 001",
			vehicleModel:  "Toyota Camry 2023",
			licenseNumber: "ABC-123",
			status:        "available",
			latitude:      37.7749,
			longitude:     -122.4194,
			altitude:      10.0,
			mileage:       15000.0,
			fuelLevel:     85.0,
		},
		{
			vin:           "2HGBH41JXMN109187",
			vehicleName:   "Fleet Vehicle 002",
			vehicleModel:  "Honda Accord 2023",
			licenseNumber: "DEF-456",
			status:        "in_use",
			latitude:      34.0522,
			longitude:     -118.2437,
			altitude:      15.0,
			mileage:       22000.0,
			fuelLevel:     60.0,
		},
		{
			vin:           "3HGBH41JXMN109188",
			vehicleName:   "Fleet Vehicle 003",
			vehicleModel:  "Ford F-150 2023",
			licenseNumber: "GHI-789",
			status:        "available",
			latitude:      40.7128,
			longitude:     -74.0060,
			altitude:      5.0,
			mileage:       8500.0,
			fuelLevel:     92.0,
		},
		{
			vin:           "4HGBH41JXMN109189",
			vehicleName:   "Fleet Vehicle 004",
			vehicleModel:  "Tesla Model 3 2023",
			licenseNumber: "JKL-012",
			status:        "maintenance",
			latitude:      41.8781,
			longitude:     -87.6298,
			altitude:      12.0,
			mileage:       35000.0,
			fuelLevel:     45.0,
		},
		{
			vin:           "5HGBH41JXMN109190",
			vehicleName:   "Fleet Vehicle 005",
			vehicleModel:  "Chevrolet Silverado 2023",
			licenseNumber: "MNO-345",
			status:        "available",
			latitude:      29.7604,
			longitude:     -95.3698,
			altitude:      8.0,
			mileage:       12000.0,
			fuelLevel:     78.0,
		},
		{
			vin:           "6HGBH41JXMN109191",
			vehicleName:   "Fleet Vehicle 006",
			vehicleModel:  "BMW X5 2023",
			licenseNumber: "PQR-678",
			status:        "in_use",
			latitude:      33.4484,
			longitude:     -112.0740,
			altitude:      20.0,
			mileage:       18000.0,
			fuelLevel:     55.0,
		},
		{
			vin:           "7HGBH41JXMN109192",
			vehicleName:   "Fleet Vehicle 007",
			vehicleModel:  "Mercedes-Benz E-Class 2023",
			licenseNumber: "STU-901",
			status:        "available",
			latitude:      39.7392,
			longitude:     -104.9903,
			altitude:      18.0,
			mileage:       9000.0,
			fuelLevel:     88.0,
		},
		{
			vin:           "8HGBH41JXMN109193",
			vehicleName:   "Fleet Vehicle 008",
			vehicleModel:  "Nissan Altima 2023",
			licenseNumber: "VWX-234",
			status:        "out_of_service",
			latitude:      32.7157,
			longitude:     -117.1611,
			altitude:      7.0,
			mileage:       45000.0,
			fuelLevel:     25.0,
		},
		{
			vin:           "9HGBH41JXMN109194",
			vehicleName:   "Fleet Vehicle 009",
			vehicleModel:  "Hyundai Sonata 2023",
			licenseNumber: "YZA-567",
			status:        "available",
			latitude:      47.6062,
			longitude:     -122.3321,
			altitude:      25.0,
			mileage:       11000.0,
			fuelLevel:     95.0,
		},
		{
			vin:           "1AHGBH41JXMN10919",
			vehicleName:   "Fleet Vehicle 010",
			vehicleModel:  "Volkswagen Passat 2023",
			licenseNumber: "BCD-890",
			status:        "in_use",
			latitude:      42.3601,
			longitude:     -71.0589,
			altitude:      14.0,
			mileage:       28000.0,
			fuelLevel:     68.0,
		},
	}

	for _, v := range seedVehicles {
		vehicleID := valueobject.GenerateVehicleID()

		licenseNumber, err := valueobject.NewLicenseNumber(v.licenseNumber)
		if err != nil {
			s.logger.Error("failed to create license number", zap.Error(err), zap.String("vin", v.vin))
			continue
		}

		status, err := valueobject.NewVehicleStatus(v.status)
		if err != nil {
			s.logger.Error("failed to create vehicle status", zap.Error(err), zap.String("vin", v.vin))
			continue
		}

		location, err := valueobject.NewLocation(v.latitude, v.longitude, v.altitude, time.Now().Unix())
		if err != nil {
			s.logger.Error("failed to create location", zap.Error(err), zap.String("vin", v.vin))
			continue
		}

		mileage, err := valueobject.NewMileage(v.mileage)
		if err != nil {
			s.logger.Error("failed to create mileage", zap.Error(err), zap.String("vin", v.vin))
			continue
		}

		fuelLevel, err := valueobject.NewFuelLevel(v.fuelLevel)
		if err != nil {
			s.logger.Error("failed to create fuel level", zap.Error(err), zap.String("vin", v.vin))
			continue
		}

		vehicle, err := entity.NewVehicle(
			vehicleID,
			v.vin,
			v.vehicleName,
			v.vehicleModel,
			licenseNumber,
			status,
			location,
			mileage,
			fuelLevel,
		)
		if err != nil {
			s.logger.Error("failed to create vehicle entity", zap.Error(err), zap.String("vin", v.vin))
			continue
		}

		if err := s.vehicleRepo.Save(ctx, vehicle); err != nil {
			s.logger.Error("failed to save vehicle", zap.Error(err), zap.String("vin", v.vin))
			continue
		}

		s.logger.Info("seeded vehicle", zap.String("vin", v.vin), zap.String("name", v.vehicleName))
	}

	s.logger.Info("vehicle seeding completed", zap.Int("total", len(seedVehicles)))
	return nil
}
