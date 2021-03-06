package mongodb

import (
	"context"
	"github.com/BellZaph/banner-roulette-backend/internal/models"
	"github.com/BellZaph/banner-roulette-backend/pkg/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type BannerRepo struct {
	db *mongo.Collection
}

func NewBannerRepo(db *mongo.Database) *BannerRepo {
	return &BannerRepo{db: db.Collection(bannerCollection)}
}

func (b *BannerRepo) Create(ctx context.Context, banner models.Banner) error {
	if banner.ID == uuid.Nil {
		banner.ID = uuid.New()
	}
	_, err := b.db.InsertOne(ctx, banner)
	return err
}

func (b *BannerRepo) Update(ctx context.Context, banner models.Banner) error {
	updateQuery := bson.M{}

	updateQuery = utils.ToMSI(banner, "bson", []string{"_id"})

	_, err := b.db.UpdateOne(ctx, bson.M{"_id": banner.ID}, bson.M{"$set": updateQuery}, options.Update().SetUpsert(true))

	return err
}

func (b *BannerRepo) IncrementCount(ctx context.Context, id uuid.UUID) error {
	_, err := b.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"show_count": 1}})
	return err
}

func (b *BannerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := b.db.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (b *BannerRepo) GetById(ctx context.Context, id uuid.UUID) (models.Banner, error) {
	var banner = models.Banner{}
	err := b.db.FindOne(ctx, bson.M{"_id": id}).Decode(&banner)
	return banner, err
}

func (b *BannerRepo) GetRandom(ctx context.Context, hour int) (models.Banner, error) {
	var banner models.Banner

	for hour > 24 {
		hour -= 24
	}

	cursor, err := b.db.Aggregate(
		ctx,
		mongo.Pipeline{
			{{"$match", bson.M{
				"show_start_date": bson.M{"$lte": time.Now()},
				"show_end_date": bson.M{"$gte": time.Now()},
				"show_hour_start": bson.M{"$lte": hour},
				"show_hour_end": bson.M{"$gte": hour},
				"visible": true,
				"$or": bson.A{
					bson.M{"show_count_cap": -1},
					bson.M{"$expr": bson.D{{"$gt", bson.A{"$show_count_cap", "$show_count"}}}},
				},
			}}},
			{{"$sample", bson.D{{"size", 1}}}},
		},
	)

	if err != nil {
		return models.Banner{}, err
	}

	if cursor.RemainingBatchLength() == 0 {
		return models.Banner{}, mongo.ErrNoDocuments
	}

	if cursor.Next(ctx) {
		_ = cursor.Decode(&banner)
	}


	return banner, err
}

func (b *BannerRepo) GetRandomLimit(ctx context.Context, hour, limit int) ([]models.Banner, error) {
	var banner []models.Banner

	for hour > 24 {
		hour -= 24
	}

	cursor, err := b.db.Aggregate(
		ctx,
		mongo.Pipeline{
			{{"$match", bson.M{
				"show_start_date": bson.M{"$lte": time.Now()},
				"show_end_date": bson.M{"$gte": time.Now()},
				"show_hour_start": bson.M{"$lte": hour},
				"show_hour_end": bson.M{"$gte": hour},
				"visible": true,
				"$or": bson.A{
					bson.M{"show_count_cap": -1},
					bson.M{"$expr": bson.D{{"$gt", bson.A{"$show_count_cap", "$show_count"}}}},
				},
			}}},
			{{"$sample", bson.D{{"size", limit}}}},
		},
	)

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &banner); err != nil {
		return nil, err
	}

	_ = cursor.Close(context.Background())
	return banner, err
}

func (b *BannerRepo) GetByPage(ctx context.Context, page int) ([]models.Banner, error) {
	var banner = make([]models.Banner, 0)

	const limit = 15
	if page <= 0 {
		page = 1
	}

	optionsFind := options.Find().SetSkip(int64((page-1) * limit)).SetLimit(limit)

	cursor, err := b.db.Find(
		ctx,
		bson.D{},
		optionsFind,
	)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &banner)
	if err != nil {
		_ = cursor.Close(context.Background())
		return nil, err
	}

	_ = cursor.Close(context.Background())
	return banner, err
}



