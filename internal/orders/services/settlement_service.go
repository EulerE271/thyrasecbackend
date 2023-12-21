package services

import (
	"errors"
	ordermodels "thyra/internal/orders/models"
	"thyra/internal/orders/repositories"
	orderutils "thyra/internal/orders/utils"
	"time"

	accountutils "thyra/internal/accounts/utils"
	transactionmodels "thyra/internal/transactions/models"
	transactionrepo "thyra/internal/transactions/repositories"
	transactionservice "thyra/internal/transactions/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SettlementService struct {
	db   *sqlx.DB
	repo *repositories.OrdersRepository
}

func NewSettlementService(db *sqlx.DB, repo *repositories.OrdersRepository) *OrdersService {
	return &OrdersService{db: db, repo: repo}
}

func (s *SettlementService) SellOrder(c *gin.Context, orderID, userIDStr string, settlementRequest ordermodels.SettlementRequest) error {

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	order, err := s.repo.GetOrder(tx, orderID)
	if err != nil {
		return err
	}

	if order.Status != ordermodels.StatusExecuted {
		return errors.New("order cannot be settle in it's current state")
	}

	assetType, err := s.repo.GetAssetType(tx, order.AssetID)
	if err != nil {
		return err
	}

	transactionType, err := s.repo.GetTransactionTypeByOrderTypeID(tx, order.OrderType)
	if err != nil {
		return err
	}

	orderNumber := orderutils.GenerateOrderNumber()
	transactionRepo := transactionrepo.NewTransactionRepository(tx)
	transactionService := transactionservice.NewTransactionService(transactionRepo, &sqlx.DB{})

	clientCashTransaction := transactionmodels.Transaction{

		Id:                        uuid.New(),
		Type:                      transactionType,
		AssetId:                   order.AssetID,
		CashAmount:                &settlementRequest.SettledAmount,
		AssetQuantity:             &settlementRequest.SettledQuantity,
		CashAccountId:             order.AccountID,
		AssetAccountId:            order.AccountID,
		AssetType:                 assetType,
		TransactionCurrency:       order.Currency,
		AssetPrice:                &order.PricePerUnit,
		CreatedById:               userID,
		UpdatedById:               userID,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
		Corrected:                 false,
		Canceled:                  false,
		Comment:                   order.Comment,
		TransactionOwnerId:        order.OwnerID,
		TransactionOwnerAccountId: order.AccountID,
		TradeDate:                 *settlementRequest.TradeDate,
		SettlementDate:            *settlementRequest.SettlementDate,
		OrderNumber:               orderNumber,
	}

	clientInstrumentTransaction := clientCashTransaction

	houseAccount, err := accountutils.GetHouseAccount(tx)
	if err != nil {
		return err
	}

	houseAccountUUID, err := uuid.Parse(houseAccount)

	clientInstrumentTransaction.TransactionOwnerAccountId = houseAccountUUID

	_, _, err = transactionService.CreateInstrumentSellTransaction(c, order.AccountID, userID.String(), &clientCashTransaction, &clientInstrumentTransaction)

	err = s.repo.UpdateOrder(tx, orderID, settlementRequest.SettledQuantity, settlementRequest.SettledAmount, "settled", settlementRequest.TradeDate, settlementRequest.SettlementDate, settlementRequest.Comment)
	if err != nil {
		return err
	}

	err = s.repo.DeductHolding(tx, order.AccountID, order.AssetID, *order.SettledQuantity)
	if err != nil {
		return err
	}

	return tx.Commit()
}
