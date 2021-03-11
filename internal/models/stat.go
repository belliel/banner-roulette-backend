package models

import "github.com/google/uuid"

type Stat struct {
	ID        uuid.UUID `json:"id" bson:"_id" db:"id"`
	GuestID   uuid.UUID `json:"guest_id" bson:"guest_id" db:"guest_id"`
	BannerID  uuid.UUID `json:"banner_id" bson:"banner_id" db:"banner_id"`
	ShowCount int       `json:"show_count" bson:"show_count" db:"show_count"`
}
