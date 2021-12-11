package user

import "wiselink/pkg/Domain/events"

type User struct {
	Id                int            `json:"id" ,bson:"id"`
	Name              int            `json:"name" ,bson:"name"`
	Email             string         `json:"email" ,bson:"email"`
	TemporaryPassword string         `json:"password" ,bson:"password"`
	AccessToken       string         `json:"accessToken" ,bson:"accessToken"`
	SuscriptedTo      []events.Event `json:"suscriptedTo" ,bson:"suscriptedTo"`
}
