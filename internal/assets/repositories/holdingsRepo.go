package repositories

import (
	"thyra/internal/assets/models" // Replace with the actual path to your models

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type HoldingsRepository interface {
	GetHoldingsByAccountId(accountId uuid.UUID) ([]models.Holding, error)
	GetAssetInformation(assetId uuid.UUID) (*models.Asset, error)
	GetCurrencyID(currencyName string) (string, error)
}
type holdingsRepository struct {
	db *sqlx.DB
}

func NewHoldingsRepository(db *sqlx.DB) HoldingsRepository {
	return &holdingsRepository{db: db}
}

// GetHoldingsByAccountId retrieves holdings for a given account ID
func (r *holdingsRepository) GetHoldingsByAccountId(accountId uuid.UUID) ([]models.Holding, error) {
	var holdings []models.Holding
	query := `SELECT id, account_id, asset_id, quantity FROM thyrasec.holdings WHERE account_id = $1`
	err := r.db.Select(&holdings, query, accountId)
	if err != nil {
		return nil, err
	}
	return holdings, nil
}

func (r *holdingsRepository) GetAssetInformation(assetId uuid.UUID) (*models.Asset, error) {
	var asset models.Asset
	query := `SELECT id, instrument_name, isin, ticker, exchange, currency, instrument_type, current_price, volume, country, sector, asset_type_id, created_at, updated_at FROM thyrasec.assets WHERE id = $1`
	err := r.db.Get(&asset, query, assetId)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *holdingsRepository) GetCurrencyID(currencyName string) (string, error) {

	var currencyID string

	query := `SELECT id FROM thyrasec.currencies WHERE name = $1`

	err := r.db.QueryRow(query, currencyName).Scan(&currencyID)
	if err != nil {
		return "", err
	}

	return currencyID, nil

}
