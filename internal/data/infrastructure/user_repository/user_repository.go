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
	DeleteUser(ctx context.Context, email string) int
	UpdateUser(ctx context.Context, u user.User) int
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
func (ur *UserRepository) CreateUser(ctx context.Context, u user.User) int {
	usersCollection := ur.Client.Database("wsMVP").Collection("users")
	_, err := usersCollection.InsertOne(ctx, u)
	if err != nil {
		return events.InternalError
	}
	return events.Success

}
func (ur *UserRepository) DeleteUser(ctx context.Context, email string) int {
	usersCollection := ur.Client.Database("wsMVP").Collection("users")
	_, err := usersCollection.DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}
func (ur *UserRepository) UpdateUser(ctx context.Context, u user.User) int {
	usersCollection := ur.Client.Database("wsMVP").Collection("users")
	_, err := usersCollection.UpdateOne(ctx, bson.M{"Id": u.Id}, bson.M{"$set": u})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}
