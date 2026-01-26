package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
)

type MongoOutboxRepository struct {
	collection *mongo.Collection
}

func NewMongoOutboxRepository(collection *mongo.Collection) *MongoOutboxRepository {
	return &MongoOutboxRepository{collection: collection}
}

type outboxDocument struct {
	ID          string `bson:"_id"`
	AggregateID string `bson:"aggregateId"`
	EventType   string `bson:"eventType"`
	EventData   string `bson:"eventData"`
	CreatedAt   int64  `bson:"createdAt"`
	PublishedAt *int64 `bson:"publishedAt,omitempty"`
}

// SaveOutboxEvent saves a domain event to the outbox for asynchronous publication.
func (r *MongoOutboxRepository) SaveOutboxEvent(ctx context.Context, aggregateID string, event interface{}) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	doc := outboxDocument{
		ID:          uuid.New().String(),
		AggregateID: aggregateID,
		EventType:   fmt.Sprintf("%T", event),
		EventData:   string(eventData),
		CreatedAt:   time.Now().UTC().Unix(),
	}

	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to save outbox event: %w", err)
	}

	return nil
}

func (r *MongoOutboxRepository) GetPendingEvents(ctx context.Context, limit int) ([]repository.OutboxEvent, error) {
	filter := bson.M{"publishedAt": nil}
	opts := options.Find().SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending events: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []outboxDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode outbox events: %w", err)
	}

	var events []repository.OutboxEvent
	for _, doc := range docs {
		events = append(events, repository.OutboxEvent{
			ID:          doc.ID,
			AggregateID: doc.AggregateID,
			EventType:   doc.EventType,
			EventData:   []byte(doc.EventData),
			CreatedAt:   doc.CreatedAt,
			PublishedAt: doc.PublishedAt,
		})
	}

	return events, nil
}

func (r *MongoOutboxRepository) MarkEventAsPublished(ctx context.Context, eventID string) error {
	now := time.Now().UTC().Unix()
	filter := bson.M{"_id": eventID}
	update := bson.M{
		"$set": bson.M{"publishedAt": now},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark event as published: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("outbox event not found: %s", eventID)
	}

	return nil
}
