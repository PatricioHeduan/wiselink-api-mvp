package user_repository

import "go.mongodb.org/mongo-driver/mongo"

type UserRepository struct {
	Client *mongo.Client
}
type UserRepositoryI interface {
}
