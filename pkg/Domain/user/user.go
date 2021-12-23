package user

type User struct {
	Id                int    `json:"id" ,bson:"id"`
	Name              string `json:"name" ,bson:"name"`
	Email             string `json:"email" ,bson:"email"`
	TemporaryPassword string `json:"password" ,bson:"password"`
	AccessToken       string `json:"accessToken" ,bson:"accessToken"`
	SuscriptedTo      []int  `json:"suscriptedTo" ,bson:"suscriptedTo"`
}
