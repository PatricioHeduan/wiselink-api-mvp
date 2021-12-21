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
	FindUserLastId(ctx context.Context) int
	CreateUser(ctx context.Context, u user.User) int
	DeleteUser(ctx context.Context, email string) int
	UpdateUser(ctx context.Context, u user.User) int
	GetLastAdminId(ctx context.Context) int
	AddAdmin(ctx context.Context, a user.Admin) int
	DeleteAdmin(ctx context.Context, a user.Admin) int
	GetAdminByEmail(ctx context.Context, email string) (int, user.Admin)
	VerifyAdminExistance(ctx context.Context, accessToken string) int
	ModifyUserEvents(ctx context.Context, u user.User) int
	GetById(ctx context.Context, id int) (int, user.User)
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
func (ur *UserRepository) FindUserLastId(ctx context.Context) int {
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

func (ur *UserRepository) GetLastAdminId(ctx context.Context) int {
	var a user.Admin
	eventsCollection := ur.Client.Database("wsMVP").Collection("Admins")
	fo := options.FindOne()
	fo.SetSort(bson.D{{"$natural", -1}})
	err := eventsCollection.FindOne(ctx, nil, fo).Decode(&a)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return 0
		}
		return -1
	}
	return a.Id
}

func (ur *UserRepository) AddAdmin(ctx context.Context, a user.Admin) int {
	adminsCollection := ur.Client.Database("wsMVP").Collection("admins")
	_, err := adminsCollection.InsertOne(ctx, a)
	if err != nil {
		return events.InternalError
	}
	return events.Success
}

func (ur *UserRepository) DeleteAdmin(ctx context.Context, a user.Admin) int {
	adminsCollection := ur.Client.Database("wsMVP").Collection("admins")
	_, err := adminsCollection.DeleteOne(ctx, bson.M{"accessToken": a.AccessToken})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}

func (ur *UserRepository) GetAdminByEmail(ctx context.Context, email string) (int, user.Admin) {
	var a user.Admin
	adminsCollection := ur.Client.Database("wsMVP").Collection("admins")
	err := adminsCollection.FindOne(ctx, bson.M{"email": email}).Decode(&a)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound, a
		}
		return user.InternalError, a
	}
	return user.Success, a
}

func (ur *UserRepository) VerifyAdminExistance(ctx context.Context, accessToken string) int {
	var a user.Admin
	adminsCollection := ur.Client.Database("wsMVP").Collection("admins")
	err := adminsCollection.FindOne(ctx, bson.M{"access_token": accessToken}).Decode(&a)
	if err != nil {
		if err.Error() == "mongo: no documents in result set" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}

func (ur *UserRepository) ModifyUserEvents(ctx context.Context, u user.User) int {
	usersCollection := ur.Client.Database("wsMVP").Collection("users")
	_, err := usersCollection.UpdateOne(ctx, bson.M{"Id": u.Id}, bson.M{"suscriptedTo": u.SuscriptedTo})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}

func (ur *UserRepository) GetById(ctx context.Context, id int) (int, user.User) {
	var u user.User
	usersCollection := ur.Client.Database("wsMVP").Collection("users")
	err := usersCollection.FindOne(ctx, bson.M{"Id": id}).Decode(&u)
	if err != nil {
		if err.Error() == "mongo: no documents in result set" {
			return user.NotFound, u
		}
		return user.InternalError, u
	}
	return user.Success, u
}
