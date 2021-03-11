package mongodb

import (
	"context"
	"github.com/BellZaph/banner-roulette-backend/internal/models"
	"github.com/BellZaph/banner-roulette-backend/pkg/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type BannerRepo struct {
	db *mongo.Collection
}

func NewBannerRepo(db *mongo.Database) *BannerRepo {
	return &BannerRepo{db: db.Collection(bannerCollection)}
}

func (b *BannerRepo) Create(ctx context.Context, banner models.Banner) error {
	_, err := b.db.InsertOne(ctx, banner)
	return err
}

func (b *BannerRepo) Update(ctx context.Context, banner models.Banner) error {
	updateQuery := bson.M{}

	updateQuery = utils.ToMSI(banner, "bson", []string{"_id"})

	_, err := b.db.UpdateOne(ctx, bson.M{"_id": banner.ID}, bson.M{"$set": updateQuery})
	return err
}

func (b *BannerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := b.db.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (b *BannerRepo) GetById(ctx context.Context, id uuid.UUID) (models.Banner, error) {
	var banner models.Banner
	err := b.db.FindOne(ctx, bson.M{"_id": id}).Decode(&banner)
	return banner, err
}

func (b *BannerRepo) GetRandom(ctx context.Context, hour int) (models.Banner, error) {
	var banner models.Banner

	for hour > 24 {
		hour -= 24
	}

	err := b.db.FindOne(
		ctx,
		bson.M{
			"show_start_date": bson.M{"$gte": time.Now()},
			"show_end_date": bson.M{"$lt": time.Now()},
			"show_hour_start": bson.M{"$gte": hour},
			"show_hour_end": bson.M{"$lt": hour},
			"$expr": bson.M{"$lt": bson.A{"show_count_cap", "show_count"}},
			"$sample": bson.M{"size": 1 },
		},
	).Decode(&banner)
	return banner, err
}

func (b *BannerRepo) GetRandomLimit(ctx context.Context, hour, limit int) ([]models.Banner, error) {
	var banner []models.Banner

	for hour > 24 {
		hour -= 24
	}

	err := b.db.FindOne(
		ctx,
		bson.M{
			"show_start_date": bson.M{"$gte": time.Now()},
			"show_end_date": bson.M{"$lt": time.Now()},
			"show_hour_start": bson.M{"$gte": hour},
			"show_hour_end": bson.M{"$lt": hour},
			"$expr": bson.M{"$lt": bson.A{"show_count_cap", "show_count"}},
			"$sample": bson.M{"size": limit },
		},
	).Decode(&banner)
	return banner, err
}

