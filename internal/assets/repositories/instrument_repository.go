package repositories

import (
	"context"
	"thyra/internal/assets/models"

	"github.com/jmoiron/sqlx"
)

type AssetsRepository struct {
	db *sqlx.DB
}

func NewAssetsRepository(db *sqlx.DB) *AssetsRepository {
	return &AssetsRepository{db: db}
}

func (r *AssetsRepository) CreateInstrument(ctx context.Context, instrument models.Instrument) (models.Instrument, error) {
	query := `INSERT INTO assets (instrument_name, isin, ticker, exchange, currency, instrument_type, 
                                 current_price, volume, country, sector, asset_type_id) 
              VALUES (:instrument_name, :isin, :ticker, :exchange, :currency, :instrument_type, 
                      :current_price, :volume, :country, :sector, :asset_type_id)
              RETURNING id, instrument_name, isin, ticker, exchange, currency, 
                        instrument_type, current_price, volume, country, sector, asset_type_id`

	row, err := r.db.NamedQueryContext(ctx, query, instrument)
	if err != nil {
		return models.Instrument{}, err
	}
	if row.Next() {
		err = row.StructScan(&instrument)
	}
	return instrument, err
}

func (r *AssetsRepository) GetAllInstruments(ctx context.Context) ([]models.Instrument, error) {
	var instruments []models.Instrument
	query := "SELECT * FROM assets"

	err := r.db.SelectContext(ctx, &instruments, query)
	return instruments, err
}

func (r *AssetsRepository) GetAllAssetTypes(ctx context.Context) ([]models.Asset, error) {
	var assets []models.Asset
	query := "SELECT * FROM asset_types"

	err := r.db.SelectContext(ctx, &assets, query)
	return assets, err
}
