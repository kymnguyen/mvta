package persistence

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/entity"
	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/valueobject"
)

// MongoVehicleRepository implements the VehicleRepository interface using MongoDB.
type MongoVehicleRepository struct {
	collection *mongo.Collection
}

// NewMongoVehicleRepository creates a new MongoDB vehicle repository.
func NewMongoVehicleRepository(collection *mongo.Collection) *MongoVehicleRepository {
	return &MongoVehicleRepository{collection: collection}
}

// vehicleDocument represents a vehicle document in MongoDB.
type vehicleDocument struct {
	ID            string  `bson:"_id"`
	VIN           string  `bson:"vin"`
	VehicleName   string  `bson:"vehicleName"`
	VehicleModel  string  `bson:"vehicleModel"`
	LicenseNumber string  `bson:"licenseNumber"`
	Status        string  `bson:"status"`
	Latitude      float64 `bson:"latitude"`
	Longitude     float64 `bson:"longitude"`
	Altitude      float64 `bson:"altitude"`
	Mileage       int64   `bson:"mileage"`
	FuelLevel     int     `bson:"fuelLevel"`
	Version       int64   `bson:"version"`
	CreatedAt     int64   `bson:"createdAt"`
	UpdatedAt     int64   `bson:"updatedAt"`
}

// Save persists a vehicle aggregate with optimistic concurrency control.
func (r *MongoVehicleRepository) Save(ctx context.Context, vehicle *entity.Vehicle) error {
	doc := vehicleDocument{
		ID:            vehicle.ID().String(),
		VIN:           vehicle.VIN(),
		VehicleName:   vehicle.VehicleName(),
		VehicleModel:  vehicle.VehicleModel(),
		LicenseNumber: vehicle.LicenseNumber().String(),
		Status:        string(vehicle.Status()),
		Latitude:      vehicle.CurrentLocation().Latitude(),
		Longitude:     vehicle.CurrentLocation().Longitude(),
		Altitude:      vehicle.CurrentLocation().Altitude(),
		Mileage:       vehicle.Mileage().Kilometers(),
		FuelLevel:     vehicle.FuelLevel().Percentage(),
		Version:       vehicle.Version().Value(),
		CreatedAt:     vehicle.CreatedAt().Unix(),
		UpdatedAt:     vehicle.UpdatedAt().Unix(),
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"_id":     vehicle.ID().String(),
		"version": vehicle.Version().Value() - 1, // Optimistic concurrency control
	}
	update := bson.M{
		"$set": doc,
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	// If upserted, no version conflict
	if result.UpsertedID != nil {
		return nil
	}

	// If no document matched, version conflict occurred
	if result.ModifiedCount == 0 && result.MatchedCount == 0 {
		return fmt.Errorf("optimistic concurrency conflict: vehicle version mismatch")
	}

	return nil
}

// FindByID retrieves a vehicle by its ID.
func (r *MongoVehicleRepository) FindByID(ctx context.Context, id valueobject.VehicleID) (*entity.Vehicle, error) {
	var doc vehicleDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("vehicle not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to find vehicle: %w", err)
	}

	status, _ := valueobject.NewVehicleStatus(doc.Status)
	location, _ := valueobject.NewLocation(doc.Latitude, doc.Longitude, doc.Altitude, 0)
	licenseNumber, _ := valueobject.NewLicenseNumber(doc.LicenseNumber)
	mileage, _ := valueobject.NewMileage(doc.Mileage)
	fuelLevel, _ := valueobject.NewFuelLevel(doc.FuelLevel)
	version, _ := valueobject.NewVersion(doc.Version)

	vehicle := entity.LoadFromHistory(
		id,
		doc.VIN,
		doc.VehicleName,
		doc.VehicleModel,
		licenseNumber,
		status,
		location,
		mileage,
		fuelLevel,
		version,
		time.Unix(doc.CreatedAt, 0),
		time.Unix(doc.UpdatedAt, 0),
	)

	return vehicle, nil
}

// FindAll retrieves all vehicles with pagination.
func (r *MongoVehicleRepository) FindAll(ctx context.Context, limit int, offset int) ([]*entity.Vehicle, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find vehicles: %w", err)
	}
	defer cursor.Close(ctx)

	var vehicles []*entity.Vehicle
	if err = cursor.All(ctx, &vehicles); err != nil {
		return nil, fmt.Errorf("failed to decode vehicles: %w", err)
	}

	var results []*entity.Vehicle
	for _, doc := range vehicles {
		var vehDoc vehicleDocument
		bsonBytes, _ := bson.Marshal(doc)
		bson.Unmarshal(bsonBytes, &vehDoc)

		vehicleID, _ := valueobject.NewVehicleID(vehDoc.ID)
		status, _ := valueobject.NewVehicleStatus(vehDoc.Status)
		location, _ := valueobject.NewLocation(vehDoc.Latitude, vehDoc.Longitude, vehDoc.Altitude, 0)
		licenseNumber, _ := valueobject.NewLicenseNumber(vehDoc.LicenseNumber)
		mileage, _ := valueobject.NewMileage(vehDoc.Mileage)
		fuelLevel, _ := valueobject.NewFuelLevel(vehDoc.FuelLevel)
		version, _ := valueobject.NewVersion(vehDoc.Version)

		vehicle := entity.LoadFromHistory(
			vehicleID,
			vehDoc.VIN,
			vehDoc.VehicleName,
			vehDoc.VehicleModel,
			licenseNumber,
			status,
			location,
			mileage,
			fuelLevel,
			version,
			time.Unix(vehDoc.CreatedAt, 0),
			time.Unix(vehDoc.UpdatedAt, 0),
		)
		results = append(results, vehicle)
	}

	return results, nil
}

// Delete removes a vehicle from storage.
func (r *MongoVehicleRepository) Delete(ctx context.Context, id valueobject.VehicleID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id.String()})
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("vehicle not found: %s", id.String())
	}

	return nil
}

// ExistsByVIN checks if a vehicle with the given VIN exists.
func (r *MongoVehicleRepository) ExistsByVIN(ctx context.Context, vin string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"vin": vin})
	if err != nil {
		return false, fmt.Errorf("failed to check vin existence: %w", err)
	}

	return count > 0, nil
}
