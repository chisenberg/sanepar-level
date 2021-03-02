package entity

import "time"

// Report -
type Report struct {
	UpdatedAt time.Time `json:"updatedAt"`
	Dams      []Dam     `json:"dams"`
}
