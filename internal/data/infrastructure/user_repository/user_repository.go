package user_repository

import (
	"context"
	"wiselink/pkg/Domain/events"
	"wiselink/pkg/Domain/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	Client *mongo.Client
}
type UserRepositoryI interface {
	GetByEmail(ctx context.Context, email string) (int, user.User)
	FindLastId(ctx context.Context) int
	CreateUser(ctx context.Context, u user.User) int
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (int, user.User) {
	var u user.User
	usersCollection := ur.Client.Database("wsMVP").Collection("users")
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		// error when no matching document is found
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound, u
		} else {
			return user.InternalError, u
		}
	} else {
		return user.Success, u
	}
}
func (ur *UserRepository) FindLastId(ctx context.Context) int {
	var u user.User
	eventsCollection := ur.Client.Database("wsMVP").Collection("users")
	fo := options.FindOne()
	fo.SetSort(bson.D{{"$natural", -1}})
	err := eventsCollection.FindOne(ctx, nil, fo).Decode(&u)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return 0
		}
		return -1
	}
	return u.Id
}
func (er *UserRepository) CreateUser(ctx context.Context, u user.User) int {
	usersCollection := er.Client.Database("wsMVP").Collection("users")
	_, err := usersCollection.InsertOne(ctx, u)
	if err != nil {
		return events.InternalError
	}
	return events.Success

}
