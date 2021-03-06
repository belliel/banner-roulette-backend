package models

import (
	"github.com/google/uuid"
	"time"
)

type Banner struct {
	ID                uuid.UUID `json:"id" bson:"_id" db:"id"`
	Name              string    `json:"name" bson:"name" db:"name"`
	RawHTML           string    `json:"raw_html" bson:"raw_html" db:"raw_html"`
	ImageURI          string    `json:"image_uri" bson:"image_uri" db:"image_uri"`
	Size              string    `json:"size" bson:"size" db:"size"`
	URI               string    `json:"uri" bson:"uri" db:"uri"`
	Alt               string    `json:"alt" bson:"alt" db:"alt"`
	ShowStartDate     time.Time `json:"show_start_date" bson:"show_start_date" db:"show_start_date"`
	ShowEndDate       time.Time `json:"show_end_date" bson:"show_end_date" db:"show_end_date"`
	ShowCountCap      int       `json:"show_count_cap" bson:"show_count_cap" db:"show_count_cap"`
	ShowCountPerGuest int       `json:"show_count_per_guest" bson:"show_count_per_guest" db:"show_count_per_guest"`
	ShowHourStart     int       `json:"show_hour_start" bson:"show_hour_start" db:"show_hour_start"`
	ShowHourEnd       int       `json:"show_hour_end" bson:"show_hour_end" db:"show_hour_end"`
	ShowCount         int       `json:"show_count" bson:"show_count" db:"show_count"`
	Visible           bool      `json:"visible" bson:"visible" db:"visible"`
}
