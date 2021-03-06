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
	DeleteUser(ctx context.Context, id int) int
	UpdateUser(ctx context.Context, u user.User, token string) int
	GetLastAdminId(ctx context.Context) int
	AddAdmin(ctx context.Context, a user.Admin) int
	DeleteAdmin(ctx context.Context, a user.Admin) int
	GetAdminByEmail(ctx context.Context, email string) (int, user.Admin)
	VerifyAdminExistance(ctx context.Context, accessToken string) int
	ModifyUserEvents(ctx context.Context, u user.User) int
	GetUserById(ctx context.Context, id int) (int, user.User)
}

//Method to get a user from id an email from the database
func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (int, user.User) {
	var u user.User
	usersCollection := ur.Client.Database("wlMVP").Collection("users")
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound, u
		} else {
			return user.InternalError, u
		}
	} else {
		return user.Success, u
	}
}

//Method to get the last user id from the database
func (ur *UserRepository) FindUserLastId(ctx context.Context) int {
	var u user.User
	usersCollection := ur.Client.Database("wlMVP").Collection("users")
	fo := options.FindOne()
	fo.SetSort(bson.D{{"$natural", -1}})
	result := usersCollection.FindOne(ctx, bson.D{}, fo)
	err := result.Decode(&u)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return 0
		}
		return -1
	}
	return u.Id
}

//Method to add a single user to the database
func (ur *UserRepository) CreateUser(ctx context.Context, u user.User) int {
	usersCollection := ur.Client.Database("wlMVP").Collection("users")
	_, err := usersCollection.InsertOne(ctx, u)
	if err != nil {
		return events.InternalError
	}
	return events.Success
}

//Method to delete a single user from the database
func (ur *UserRepository) DeleteUser(ctx context.Context, id int) int {
	usersCollection := ur.Client.Database("wlMVP").Collection("users")
	_, err := usersCollection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}

//Method to modify a single user from the database
func (ur *UserRepository) UpdateUser(ctx context.Context, u user.User, token string) int {
	usersCollection := ur.Client.Database("wlMVP").Collection("users")
	_, err := usersCollection.UpdateOne(ctx, bson.M{"id": u.Id}, bson.M{"$set": bson.M{"name": u.Name, "accesstoken": token}})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}

//Method to get the last admin id from the database
func (ur *UserRepository) GetLastAdminId(ctx context.Context) int {
	var a user.Admin
	eventsCollection := ur.Client.Database("wlMVP").Collection("admins")
	fo := options.FindOne()
	fo.SetSort(bson.D{{"$natural", -1}})
	err := eventsCollection.FindOne(ctx, bson.D{}, fo).Decode(&a)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return 0
		}
		return -1
	}
	return a.Id
}

//Method to add a single admin to the database
func (ur *UserRepository) AddAdmin(ctx context.Context, a user.Admin) int {
	adminsCollection := ur.Client.Database("wlMVP").Collection("admins")
	_, err := adminsCollection.InsertOne(ctx, a)
	if err != nil {
		return events.InternalError
	}
	return events.Success
}

//Method to delete a single admin from the database
func (ur *UserRepository) DeleteAdmin(ctx context.Context, a user.Admin) int {
	adminsCollection := ur.Client.Database("wlMVP").Collection("admins")
	_, err := adminsCollection.DeleteOne(ctx, bson.M{"accesstoken": a.AccessToken})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}

//Method to get a single admin from an email from the database
func (ur *UserRepository) GetAdminByEmail(ctx context.Context, email string) (int, user.Admin) {
	var a user.Admin
	adminsCollection := ur.Client.Database("wlMVP").Collection("admins")
	err := adminsCollection.FindOne(ctx, bson.M{"email": email}).Decode(&a)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound, a
		}
		return user.InternalError, a
	}
	return user.Success, a
}

//Method to get a single admin from an accessToken from the database
func (ur *UserRepository) VerifyAdminExistance(ctx context.Context, accessToken string) int {
	var a user.Admin
	adminsCollection := ur.Client.Database("wlMVP").Collection("admins")
	err := adminsCollection.FindOne(ctx, bson.M{"accesstoken": accessToken}).Decode(&a)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}

//Method to modify the users inscriptions in the database
func (ur *UserRepository) ModifyUserEvents(ctx context.Context, u user.User) int {
	usersCollection := ur.Client.Database("wlMVP").Collection("users")
	_, err := usersCollection.UpdateOne(ctx, bson.M{"id": u.Id}, bson.M{"$set": bson.M{"suscriptedto": u.SuscriptedTo}})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound
		}
		return user.InternalError
	}
	return user.Success
}

//Method to get a single user from an id from the database
func (ur *UserRepository) GetUserById(ctx context.Context, id int) (int, user.User) {
	var u user.User
	usersCollection := ur.Client.Database("wlMVP").Collection("users")
	err := usersCollection.FindOne(ctx, bson.M{"id": id}).Decode(&u)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return user.NotFound, u
		}
		return user.InternalError, u
	}
	return user.Success, u
}
