package service

import (
	"context"
	"github.com/BellZaph/banner-roulette-backend/internal/models"
	"github.com/BellZaph/banner-roulette-backend/internal/repository"
	"github.com/BellZaph/banner-roulette-backend/pkg/hash"
	"github.com/google/uuid"
)

type BannersService struct {
	repo repository.Banners

	hasher hash.Hasher
}

func NewBannersService(repo repository.Banners, hasher hash.Hasher) *BannersService {
	return &BannersService{repo: repo, hasher: hasher}
}

func (b *BannersService) Create(ctx context.Context, input BannerCreateInput) error {

	banner := models.Banner{
		Name:              input.Name,
		RawHTML:           input.RawHTML,
		ImageURI:          input.ImageURI,
		Size:              input.Size,
		URI:               input.URI,
		Alt:               input.Alt,
		ShowStartDate:     input.ShowStartDate,
		ShowEndDate:       input.ShowEndDate,
		ShowCountCap:      input.ShowCountCap,
		ShowHourStart:     input.ShowHourStart,
		ShowHourEnd:       input.ShowHourEnd,
		ShowCount:         0,
		Visible:           input.Visible,
	}

	if err := b.repo.Create(ctx, banner); err != nil {
		return err
	}

	return nil
}

func (b *BannersService) Update(ctx context.Context, input BannerUpdateInput) error {

	banner := models.Banner{
		Name:              input.Name,
		RawHTML:           input.RawHTML,
		ImageURI:          input.ImageURI,
		Size:              input.Size,
		URI:               input.URI,
		Alt:               input.Alt,
		ShowStartDate:     input.ShowStartDate,
		ShowEndDate:       input.ShowEndDate,
		ShowCountCap:      input.ShowCountCap,
		ShowHourStart:     input.ShowHourStart,
		ShowHourEnd:       input.ShowHourEnd,
		ShowCount:         0,
		Visible:           input.Visible,
	}

	return b.repo.Update(ctx, banner)
}

func (b *BannersService) IncrementCount(ctx context.Context, id uuid.UUID) error {
	return b.repo.IncrementCount(ctx, id)
}

func (b *BannersService) Delete(ctx context.Context, id uuid.UUID) error {
	return b.repo.Delete(ctx, id)
}

func (b *BannersService) GetById(ctx context.Context, id uuid.UUID) (models.Banner, error) {
	return b.repo.GetById(ctx, id)
}

func (b *BannersService) GetRandom(ctx context.Context, hour int) (models.Banner, error) {
	return b.repo.GetRandom(ctx, hour)
}

func (b *BannersService) GetRandoms(ctx context.Context, hour int, limit int) ([]models.Banner, error) {
	return b.repo.GetRandomLimit(ctx, hour, limit)
}

