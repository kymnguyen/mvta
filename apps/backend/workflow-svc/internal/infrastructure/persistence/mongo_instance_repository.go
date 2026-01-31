package persistence

import (
	"context"
	"time"

	"workflow-svc/internal/domain/repository"
	"workflow-svc/internal/domain/workflow"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoInstanceRepository struct {
	collection *mongo.Collection
}

func NewMongoInstanceRepository(db *mongo.Database) (repository.InstanceRepository, error) {
	collection := db.Collection("workflow_instances")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "workflow_name", Value: 1},
				{Key: "current_state", Value: 1},
			},
		},
		{
			Keys:    bson.D{{Key: "correlation_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "timeout_at", Value: 1}},
			Options: options.Index().SetPartialFilterExpression(
				bson.D{{Key: "timeout_at", Value: bson.D{{Key: "$exists", Value: true}}}},
			),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, err
	}

	return &mongoInstanceRepository{collection: collection}, nil
}

func (r *mongoInstanceRepository) Create(ctx context.Context, instance *workflow.WorkflowInstance) error {
	instance.Version = 1
	instance.CreatedAt = time.Now()
	instance.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, instance)
	return err
}

func (r *mongoInstanceRepository) Update(ctx context.Context, instance *workflow.WorkflowInstance) error {
	filter := bson.M{
		"_id":     instance.ID,
		"version": instance.Version,
	}

	update := bson.M{
		"$set": bson.M{
			"current_state": instance.CurrentState,
			"context":       instance.Context,
			"history":       instance.History,
			"updated_at":    time.Now(),
			"timeout_at":    instance.TimeoutAt,
		},
		"$inc": bson.M{"version": 1},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return workflow.ErrConcurrentModification
	}

	instance.Version++
	instance.UpdatedAt = time.Now()
	return nil
}

func (r *mongoInstanceRepository) FindByID(ctx context.Context, id string) (*workflow.WorkflowInstance, error) {
	var instance workflow.WorkflowInstance
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&instance)
	if err == mongo.ErrNoDocuments {
		return nil, workflow.ErrInstanceNotFound
	}
	return &instance, err
}

func (r *mongoInstanceRepository) FindByCorrelationID(ctx context.Context, correlationID string) (*workflow.WorkflowInstance, error) {
	var instance workflow.WorkflowInstance
	err := r.collection.FindOne(ctx, bson.M{"correlation_id": correlationID}).Decode(&instance)
	if err == mongo.ErrNoDocuments {
		return nil, workflow.ErrInstanceNotFound
	}
	return &instance, err
}

func (r *mongoInstanceRepository) List(ctx context.Context, filter workflow.InstanceFilter) ([]*workflow.WorkflowInstance, error) {
	query := bson.M{}
	if filter.WorkflowName != "" {
		query["workflow_name"] = filter.WorkflowName
	}
	if filter.State != "" {
		query["current_state"] = filter.State
	}
	if filter.CorrelationID != "" {
		query["correlation_id"] = filter.CorrelationID
	}

	cursor, err := r.collection.Find(ctx, query, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var instances []*workflow.WorkflowInstance
	if err := cursor.All(ctx, &instances); err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *mongoInstanceRepository) FindPendingTimeouts(ctx context.Context, limit int) ([]*workflow.WorkflowInstance, error) {
	query := bson.M{
		"timeout_at": bson.M{"$lte": time.Now()},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timeout_at", Value: 1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var instances []*workflow.WorkflowInstance
	if err := cursor.All(ctx, &instances); err != nil {
		return nil, err
	}
	return instances, nil
}
