package services

import (
	"errors"
	"fmt"
	"thyra/internal/transactions/models"
	"thyra/internal/transactions/repositories"
	orderno "thyra/internal/transactions/services/ordernumber"
	"thyra/internal/transactions/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func InsertTransaction(database *sqlx.DB, transaction *models.Transaction, query string) (uuid.UUID, error) {
	var transactionID uuid.UUID
	tx, err := database.Beginx()
	if err != nil {
		return transactionID, err
	}

	namedStmt, err := tx.PrepareNamed(query)
	if err != nil {
		tx.Rollback()
		return transactionID, err
	}

	err = namedStmt.QueryRowx(transaction.ToMap()).Scan(&transactionID)
	if err != nil {
		tx.Rollback()
		return transactionID, err
	}

	err = tx.Commit()
	if err != nil {
		return transactionID, err
	}

	return transactionID, nil
}

func InsertParentTransaction(database *sqlx.DB, transaction *models.Transaction) error {
	query := `
        INSERT INTO thyrasec.transactions (
            id, type, asset1_id, asset2_id, amount_asset1, amount_asset2, created_by_id,
            updated_by_id, created_at, updated_at, corrected, canceled, status_transaction,
            comment, transaction_owner_id, account_owner_id, account_asset1_id,
            account_asset2_id, trade_date, settlement_date, order_no
        )
        VALUES (
            :id, :type, :asset1_id, :asset2_id, :amount_asset1, :amount_asset2, :created_by_id,
            :updated_by_id, :created_at, :updated_at, :corrected, :canceled, :status_transaction,
            :comment, :transaction_owner_id, :account_owner_id, :account_asset1_id,
            :account_asset2_id, :trade_date, :settlement_date, :order_no
        )
    ` // Removed the RETURNING id part

	_, err := database.NamedExec(query, transaction)
	if err != nil {
		return err
	}

	return nil
}

type TransactionService struct {
	transactionRepo repositories.TransactionRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}

func (s *TransactionService) CreateDeposit(c *gin.Context, userID string, transactionData *models.Transaction) (uuid.UUID, uuid.UUID, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid user ID")
	}

	OrderNumber := orderno.GenerateOrderNumber()
	houseAccountID, err := utils.GetHouseAccount(c) // Assuming GetHouseAccount does not require context
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	houseAccountUUID, err := uuid.Parse(houseAccountID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	debitTransactionID := uuid.New()
	creditTransactionID := uuid.New()

	// Setup debit and credit transactions
	creditTransaction := *transactionData
	creditTransaction.Id = creditTransactionID
	creditTransaction.CreatedAt = time.Now()
	creditTransaction.UpdatedAt = time.Now()
	creditTransaction.CreatedById = userUUID
	creditTransaction.UpdatedById = userUUID
	creditTransaction.OrderNumber = OrderNumber

	debitTransaction := models.Transaction{
		Id:                 debitTransactionID,
		Type:               creditTransaction.Type,
		CreatedById:        userUUID,
		UpdatedById:        userUUID,
		Asset1Id:           creditTransaction.Asset1Id,
		AmountAsset1:       creditTransaction.AmountAsset1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		StatusTransaction:  creditTransaction.StatusTransaction,
		Comment:            creditTransaction.Comment,
		TransactionOwnerId: creditTransaction.TransactionOwnerId,
		AccountOwnerId:     houseAccountUUID,
		AccountAsset1Id:    creditTransaction.AccountAsset2Id,
		AccountAsset2Id:    uuid.Nil,
		OrderNumber:        OrderNumber,
		Trade_date:         creditTransaction.Trade_date,
		Settlement_date:    creditTransaction.Settlement_date,
	}

	creditTransaction.AmountAsset1 = 0
	creditTransaction.AccountAsset1Id = uuid.Nil
	creditTransaction.Asset1Id = uuid.Nil

	/* CHECKS THE BALANCE OF THE CUSTOMER ACCOUNT */
	currentBalance, err := s.transactionRepo.GetAccountBalance(transactionData.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	currentAvailableBalance, err := s.transactionRepo.GetAccountAvailableBalance(transactionData.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	/*BALANCE CHEKC FOR HOUSE */
	currentBalanceHouse, err := s.transactionRepo.GetAccountBalance(debitTransaction.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	currentAvailableBalanceHouse, err := s.transactionRepo.GetAccountAvailableBalance(debitTransaction.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	/* CALCULATES BALANCE FOR CUSTOMER */
	newBalance := currentBalance + transactionData.AmountAsset1
	availableBalance := currentAvailableBalance + transactionData.AmountAsset1

	/* CALCULATES BALANCE FOR HOUSE */
	newBalanceHouse := currentBalanceHouse + transactionData.AmountAsset1
	availableBalanceHouse := currentAvailableBalanceHouse + transactionData.AmountAsset1

	err = s.transactionRepo.UpdateAccountBalance(creditTransaction.AccountAsset1Id, newBalance, availableBalance)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.UpdateAccountBalance(debitTransaction.AccountAsset1Id, newBalanceHouse, availableBalanceHouse)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// Insert into database using repository
	err = s.transactionRepo.InsertTransaction(&debitTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.InsertTransaction(&creditTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return debitTransactionID, creditTransactionID, nil
}

func (s *TransactionService) CreateInstrumentPurchaseTransaction(c *gin.Context, userID string, transactionData *models.Transaction, transactionInstrumentData *models.Transaction) (uuid.UUID, uuid.UUID, error) {

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid user ID")
	}

	OrderNumber := orderno.GenerateOrderNumber()
	houseAccountID, err := utils.GetHouseAccount(c) // Assuming GetHouseAccount does not require context
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	houseAccountUUID, err := uuid.Parse(houseAccountID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	debitTransactionID := uuid.New()
	creditTransactionID := uuid.New()
	debitInstrumentTransactionID := uuid.New()
	creditInstrumentTransactionID := uuid.New()

	// Debit transaction - already bound from the request
	debitTransaction := transactionData
	debitTransaction.Id = debitTransactionID
	debitTransaction.CreatedAt = time.Now()
	debitTransaction.UpdatedAt = time.Now()
	debitTransaction.CreatedById = userUUID
	debitTransaction.UpdatedById = userUUID
	debitTransaction.OrderNumber = OrderNumber

	creditTransaction := &models.Transaction{
		Id:                 creditTransactionID,
		Type:               debitTransaction.Type,
		CreatedById:        userUUID,
		UpdatedById:        userUUID,
		Asset1Id:           debitTransaction.Asset1Id,
		AmountAsset1:       debitTransaction.AmountAsset1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		StatusTransaction:  debitTransaction.StatusTransaction,
		Comment:            debitTransaction.Comment,
		TransactionOwnerId: debitTransaction.TransactionOwnerId,
		AccountOwnerId:     houseAccountUUID,
		AccountAsset1Id:    debitTransaction.AccountAsset2Id,
		AccountAsset2Id:    uuid.Nil,
		OrderNumber:        OrderNumber,
		Trade_date:         debitTransaction.Trade_date,
		Settlement_date:    debitTransaction.Settlement_date,
	}

	/*Configures the instrument transaction */
	debitInstrumentTransaction := transactionInstrumentData
	debitInstrumentTransaction.Id = debitInstrumentTransactionID
	debitInstrumentTransaction.CreatedAt = time.Now()
	debitInstrumentTransaction.UpdatedAt = time.Now()
	debitInstrumentTransaction.CreatedById = userUUID
	debitInstrumentTransaction.UpdatedById = userUUID
	debitInstrumentTransaction.OrderNumber = OrderNumber

	creditInstrumentTransaction := &models.Transaction{
		Id:                 creditInstrumentTransactionID,
		Type:               debitInstrumentTransaction.Type,
		CreatedById:        userUUID,
		UpdatedById:        userUUID,
		Asset1Id:           debitInstrumentTransaction.Asset1Id,
		AmountAsset1:       debitInstrumentTransaction.AmountAsset1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		StatusTransaction:  debitInstrumentTransaction.StatusTransaction,
		Comment:            debitInstrumentTransaction.Comment,
		TransactionOwnerId: debitInstrumentTransaction.TransactionOwnerId,
		AccountOwnerId:     houseAccountUUID,
		AccountAsset1Id:    debitInstrumentTransaction.AccountAsset2Id,
		AccountAsset2Id:    uuid.Nil,
		OrderNumber:        OrderNumber,
		Trade_date:         debitInstrumentTransaction.Trade_date,
		Settlement_date:    debitInstrumentTransaction.Settlement_date,
	}

	/* CHECKS THE BALANCE OF THE CUSTOMER ACCOUNT */
	currentBalance, err := s.transactionRepo.GetAccountBalance(transactionData.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	currentAvailableBalance, err := s.transactionRepo.GetAccountAvailableBalance(transactionData.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	/*BALANCE CHEKC FOR HOUSE */
	currentBalanceHouse, err := s.transactionRepo.GetAccountBalance(creditTransaction.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	currentAvailableBalanceHouse, err := s.transactionRepo.GetAccountAvailableBalance(creditTransaction.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	/* CALCULATES BALANCE FOR CUSTOMER */
	newBalance := currentBalance - transactionData.AmountAsset1
	availableBalance := currentAvailableBalance - transactionData.AmountAsset1

	/* CALCULATES BALANCE FOR HOUSE */
	newBalanceHouse := currentBalanceHouse - transactionData.AmountAsset1
	availableBalanceHouse := currentAvailableBalanceHouse - transactionData.AmountAsset1

	transactionAmount := transactionData.AmountAsset1

	if availableBalance < transactionAmount {
		return uuid.Nil, uuid.Nil, fmt.Errorf("insufficient funds in customer's account")
	}

	if availableBalanceHouse < transactionAmount {
		return uuid.Nil, uuid.Nil, fmt.Errorf("insufficient funds in house account")
	}

	err = s.transactionRepo.UpdateAccountBalance(debitTransaction.AccountAsset1Id, newBalance, availableBalance)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.UpdateAccountBalance(creditTransaction.AccountAsset1Id, newBalanceHouse, availableBalanceHouse)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// Insert into database using repository
	err = s.transactionRepo.InsertTransaction(debitTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.InsertTransaction(creditTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	//Inserts instrument transactions
	err = s.transactionRepo.InsertTransaction(debitInstrumentTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.InsertTransaction(creditInstrumentTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	err = s.transactionRepo.InsertTransaction(creditInstrumentTransaction)

	return debitTransactionID, creditTransactionID, nil

}

func (s *TransactionService) CreateWithdrawal(c *gin.Context, userID string, transactionData *models.Transaction) (uuid.UUID, uuid.UUID, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid user ID")
	}

	OrderNumber := orderno.GenerateOrderNumber()
	houseAccountID, err := utils.GetHouseAccount(c) // Assuming GetHouseAccount does not require context
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	houseAccountUUID, err := uuid.Parse(houseAccountID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	debitTransactionID := uuid.New()
	creditTransactionID := uuid.New()

	// Debit transaction - already bound from the request
	debitTransaction := transactionData
	debitTransaction.Id = debitTransactionID
	debitTransaction.CreatedAt = time.Now()
	debitTransaction.UpdatedAt = time.Now()
	debitTransaction.CreatedById = userUUID
	debitTransaction.UpdatedById = userUUID
	debitTransaction.OrderNumber = OrderNumber

	creditTransaction := &models.Transaction{
		Id:                 creditTransactionID, // Assign the generated UUID
		Type:               debitTransaction.Type,
		CreatedById:        userUUID,
		UpdatedById:        userUUID,
		Asset1Id:           debitTransaction.Asset1Id,
		AmountAsset1:       debitTransaction.AmountAsset1,
		CreatedAt:          time.Now(), // Use time.Now() to set to current time
		UpdatedAt:          time.Now(),
		StatusTransaction:  debitTransaction.StatusTransaction,
		Comment:            debitTransaction.Comment,
		TransactionOwnerId: debitTransaction.TransactionOwnerId,
		AccountOwnerId:     houseAccountUUID,
		AccountAsset1Id:    debitTransaction.AccountAsset2Id,
		AccountAsset2Id:    uuid.Nil,
		OrderNumber:        OrderNumber,
		Trade_date:         debitTransaction.Trade_date,      // Set to zero time if not used
		Settlement_date:    debitTransaction.Settlement_date, // Set to zero time if not used
	}

	/* CHECKS THE BALANCE OF THE CUSTOMER ACCOUNT */
	currentBalance, err := s.transactionRepo.GetAccountBalance(transactionData.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	currentAvailableBalance, err := s.transactionRepo.GetAccountAvailableBalance(transactionData.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	/*BALANCE CHEKC FOR HOUSE */
	currentBalanceHouse, err := s.transactionRepo.GetAccountBalance(creditTransaction.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	currentAvailableBalanceHouse, err := s.transactionRepo.GetAccountAvailableBalance(creditTransaction.AccountAsset1Id)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	/* CALCULATES BALANCE FOR CUSTOMER */
	newBalance := currentBalance - transactionData.AmountAsset1
	availableBalance := currentAvailableBalance - transactionData.AmountAsset1

	/* CALCULATES BALANCE FOR HOUSE */
	newBalanceHouse := currentBalanceHouse - transactionData.AmountAsset1
	availableBalanceHouse := currentAvailableBalanceHouse - transactionData.AmountAsset1

	transactionAmount := transactionData.AmountAsset1

	if availableBalance < transactionAmount {
		return uuid.Nil, uuid.Nil, fmt.Errorf("insufficient funds in customer's account")
	}

	if availableBalanceHouse < transactionAmount {
		return uuid.Nil, uuid.Nil, fmt.Errorf("insufficient funds in house account")
	}

	err = s.transactionRepo.UpdateAccountBalance(debitTransaction.AccountAsset1Id, newBalance, availableBalance)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.UpdateAccountBalance(creditTransaction.AccountAsset1Id, newBalanceHouse, availableBalanceHouse)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// Insert into database using repository
	err = s.transactionRepo.InsertTransaction(debitTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.InsertTransaction(creditTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return debitTransactionID, creditTransactionID, nil
}
