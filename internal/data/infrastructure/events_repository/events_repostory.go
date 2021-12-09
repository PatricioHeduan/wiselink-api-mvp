package events_repository

import (
	"context"
	"wiselink/pkg/Domain/events"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventsRepository struct {
	Client *mongo.Client
}
type EventsRepositoryI interface {
	FindLastId(ctx context.Context) int
	CreateEvent(ctx context.Context, e events.Event) int
	UpdateEvent(ctx context.Context, e events.Event) int
	DeleteEvent(ctx context.Context, id int) int
}

func (er *EventsRepository) FindLastId(ctx context.Context) int {
	var e events.Event
	eventsCollection := er.Client.Database("wsMVP").Collection("events")
	fo := options.FindOne()
	fo.SetSort(bson.D{{"$natural", -1}})
	err := eventsCollection.FindOne(ctx, nil, fo).Decode(&e)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return 0
		}
		return -1
	}
	return e.Id
}
func (er *EventsRepository) CreateEvent(ctx context.Context, e events.Event) int {
	eventsCollection := er.Client.Database("wsMVP").Collection("events")
	_, err := eventsCollection.InsertOne(ctx, e)
	if err != nil {
		return events.InternalError
	}
	return events.Success

}

func (er *EventsRepository) UpdateEvent(ctx context.Context, e events.Event) int {
	eventsCollection := er.Client.Database("wsMVP").Collection("events")
	_, err := eventsCollection.UpdateOne(ctx, bson.M{"id": e.Id}, bson.M{"$set": e})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return events.NotFound
		} else {
			return events.InternalError
		}
	}
	return events.Success
}

func (er *EventsRepository) DeleteEvent(ctx context.Context, id int) int {
	eventsCollection := er.Client.Database("wsMVP").Collection("events")
	_, err := eventsCollection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return events.NotFound
		} else {
			return events.InternalError
		}
	}
	return events.Success
}
