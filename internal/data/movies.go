package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	RunTime   Runtime     `json:"runtime,omitzero"` // string directive would convert value to string
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
