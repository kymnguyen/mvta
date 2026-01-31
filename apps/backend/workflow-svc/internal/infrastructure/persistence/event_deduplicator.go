package persistence

import (
	"context"
	"time"

	"workflow-svc/internal/domain/workflow"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventDeduplication struct {
	EventID     string    `bson:"_id"`
	InstanceID  string    `bson:"instance_id"`
	ProcessedAt time.Time `bson:"processed_at"`
}

type EventDeduplicator struct {
	collection *mongo.Collection
}

func NewEventDeduplicator(db *mongo.Database) (*EventDeduplicator, error) {
	collection := db.Collection("processed_events")

	// Create TTL index to auto-delete old events after 7 days
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "processed_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(7 * 24 * 60 * 60), // 7 days
	}

	_, err := collection.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		return nil, err
	}

	return &EventDeduplicator{collection: collection}, nil
}

func (d *EventDeduplicator) MarkProcessed(ctx context.Context, eventID, instanceID string) error {
	doc := EventDeduplication{
		EventID:     eventID,
		InstanceID:  instanceID,
		ProcessedAt: time.Now(),
	}

	_, err := d.collection.InsertOne(ctx, doc)
	if mongo.IsDuplicateKeyError(err) {
		return workflow.ErrDuplicateEvent
	}
	return err
}

func (d *EventDeduplicator) IsProcessed(ctx context.Context, eventID string) (bool, error) {
	count, err := d.collection.CountDocuments(ctx, bson.M{"_id": eventID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
