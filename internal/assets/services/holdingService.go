package services

import (
	"thyra/internal/assets/models"       // Replace with the actual path to your models
	"thyra/internal/assets/repositories" // Replace with the actual path to your repositories

	"github.com/google/uuid"
)

type HoldingsService interface {
	GetAccountHoldings(accountId uuid.UUID) ([]models.Holding, error)
	GetAssetDetails(assetId uuid.UUID) (*models.Asset, error)
	GetHoldingsWithAssetDetails(accountId uuid.UUID) ([]HoldingWithAssetDetails, error)
}

type HoldingWithAssetDetails struct {
	Holding models.Holding
	Asset   models.Asset
}

type holdingsService struct {
	repo repositories.HoldingsRepository
}

func NewHoldingsService(repo repositories.HoldingsRepository) HoldingsService {
	return &holdingsService{repo: repo}
}

// GetAccountHoldings retrieves holdings for a specific account
func (s *holdingsService) GetAccountHoldings(accountId uuid.UUID) ([]models.Holding, error) {
	// Here you can add additional business logic before/after fetching the holdings
	return s.repo.GetHoldingsByAccountId(accountId)
}

// GetAssetDetails retrieves detailed information about a specific asset
func (s *holdingsService) GetAssetDetails(assetId uuid.UUID) (*models.Asset, error) {
	// Additional business logic can be added here if necessary
	return s.repo.GetAssetInformation(assetId)
}

func (s *holdingsService) GetHoldingsWithAssetDetails(accountId uuid.UUID) ([]HoldingWithAssetDetails, error) {
	holdings, err := s.repo.GetHoldingsByAccountId(accountId)
	if err != nil {
		return nil, err
	}

	var holdingsWithDetails []HoldingWithAssetDetails
	for _, holding := range holdings {
		asset, err := s.repo.GetAssetInformation(holding.AssetID)
		if err != nil {
			return nil, err
		}
		holdingsWithDetails = append(holdingsWithDetails, HoldingWithAssetDetails{
			Holding: holding,
			Asset:   *asset,
		})
	}

	return holdingsWithDetails, nil
}