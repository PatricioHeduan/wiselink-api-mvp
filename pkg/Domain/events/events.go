package events

type Event struct {
	Id        int    `json:"id" ,bson:"id"`
	Title     string `json:"title" ,bson:"title"`
	ShortD    string `json:"shortD" ,bson:"shortD"`
	LongD     string `json:"longD" ,bson:"longD"`
	Date      string `json:"date" ,bson:"date"`
	Hour      string `json:"hour" ,bson:"hour"`
	Organizer string `json:"organizer" ,bson:"organizer"`
	Place     string `json:"place" ,bson:"place"`
	Status    bool   `json:"status" ,bson:"status"`
}
