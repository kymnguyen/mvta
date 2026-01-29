package persistence

import (
"context"

"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"

"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/entity"
)

type MongoVehicleChangeHistoryRepository struct {
	collection *mongo.Collection
}

func NewMongoVehicleChangeHistoryRepository(collection *mongo.Collection) *MongoVehicleChangeHistoryRepository {
	return &MongoVehicleChangeHistoryRepository{collection: collection}
}

func (r *MongoVehicleChangeHistoryRepository) Save(ctx context.Context, history *entity.VehicleChangeHistory) error {
	_, err := r.collection.InsertOne(ctx, history)
	return err
}

func (r *MongoVehicleChangeHistoryRepository) FindByVehicleID(ctx context.Context, vehicleID string, limit int, offset int) ([]*entity.VehicleChangeHistory, error) {
	filter := bson.M{"vehicleId": vehicleID}
	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"changedAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var histories []*entity.VehicleChangeHistory
	if err := cursor.All(ctx, &histories); err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *MongoVehicleChangeHistoryRepository) FindByChangeType(ctx context.Context, changeType string, limit int, offset int) ([]*entity.VehicleChangeHistory, error) {
	filter := bson.M{"changeType": changeType}
	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"changedAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var histories []*entity.VehicleChangeHistory
	if err := cursor.All(ctx, &histories); err != nil {
		return nil, err
	}

	return histories, nil
}
