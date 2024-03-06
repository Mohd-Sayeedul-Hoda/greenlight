package data

import(
	"time"
)

type Movie struct{
	ID int64
	CreatedAt time.Time
	Title string
	Year int32
	Runtime int32 // movie lenght
	Genres []string
	Version int32
}
