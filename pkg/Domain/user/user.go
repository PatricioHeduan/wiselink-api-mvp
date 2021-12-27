package user

type User struct {
	Id                int    `json:"id" ,bson:"id"`
	Name              string `json:"name" ,bson:"name"`
	Email             string `json:"email" ,bson:"email"`
	TemporaryPassword string `json:"password" ,bson:"password"` //temporaryPassword por el motivo de que solo llega cuando tiene que encodearse o decodearse, nunca se guarda en la base de datos
	AccessToken       string `json:"accessToken" ,bson:"accessToken"`
	SuscriptedTo      []int  `json:"suscriptedTo" ,bson:"suscriptedTo"` //slices de indentificadores de eventos a los que el usuario está inscripto

}
