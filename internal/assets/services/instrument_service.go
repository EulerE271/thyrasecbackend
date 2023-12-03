package services

import (
	"context"
	"errors"
	"thyra/internal/assets/models"
	"thyra/internal/assets/repositories"
)

type AssetsService struct {
	repo *repositories.AssetsRepository
}

func NewAssetsService(repo *repositories.AssetsRepository) *AssetsService {
	return &AssetsService{repo: repo}
}

func (s *AssetsService) CreateInstrument(ctx context.Context, instrument models.Instrument, userRole string) (models.Instrument, error) {
	if userRole != "admin" {
		return models.Instrument{}, errors.New("only admins can create instruments")
	}

	return s.repo.CreateInstrument(ctx, instrument)
}

func (s *AssetsService) GetAllInstruments(ctx context.Context, userRole string) ([]models.Instrument, error) {
	if userRole != "admin" {
		return nil, errors.New("only admins can fetch all instruments")
	}

	return s.repo.GetAllInstruments(ctx)
}

func (s *AssetsService) GetAllAssetTypes(ctx context.Context, userRole string) ([]models.Asset, error) {
	if userRole != "admin" {
		return nil, errors.New("only admins can fetch all asset types")
	}

	return s.repo.GetAllAssetTypes(ctx)
}
