package user

type Admin struct {
	Id          int    `json:"id" ,bson:"id"`
	Email       string `json:"email" ,bson:"email"`
	AccessToken string `json:"accessToken" ,bson:"accessToken"`
}
