package repository

import (
	"context"
	"github.com/BellZaph/banner-roulette-backend/internal/models"
	"github.com/BellZaph/banner-roulette-backend/internal/repository/mongodb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Banner interface {
	Create(ctx context.Context, banner models.Banner) error
	Update(ctx context.Context, banner models.Banner) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (models.Banner, error)
	GetRandom(ctx context.Context, hour int) (models.Banner, error)
	GetRandomLimit(ctx context.Context, hour, limit int) ([]models.Banner, error)
}

type Repository struct {
	Banner Banner
}

func NewRepository(database interface{}) (*Repository, error) {

	r := &Repository{}

	if db, ok := database.(*mongo.Database); ok {

		r.Banner = mongodb.NewBannerRepo(db)

		return r, nil
	}


	return nil, ErrNotImplemented
}
