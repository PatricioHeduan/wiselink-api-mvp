package filters

import "time"

type Filter struct {
	Date   time.Time `json:"time,ommitempty"`
	Status bool      `json:"status,ommitempty"`
	Title  string    `json:"title,omitempty"`
}
