package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	RunTime   int32     `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitempty"`
	Vesion    int32     `json:"version"`
}
