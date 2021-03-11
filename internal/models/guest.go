package models

import "github.com/google/uuid"

type Guest struct {
	ID          uuid.UUID `json:"id" bson:"_id" db:"id"`
	Fingerprint string    `json:"fingerprint" bson:"fingerprint" db:"fingerprint"`
}
