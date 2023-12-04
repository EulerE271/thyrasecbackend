package services

import (
	"errors"
	"fmt"
	"strconv"
	accountutils "thyra/internal/accounts/utils"
	orderutils "thyra/internal/orders/utils"
	"thyra/internal/transactions/models"
	"thyra/internal/transactions/repositories"
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
            updated_by_id, created_at, updated_at, corrected, canceled,
            comment, transaction_owner_id, account_owner_id, account_asset1_id,
            account_asset2_id, trade_date, settlement_date, order_no
        )
        VALUES (
            :id, :type, :asset1_id, :asset2_id, :amount_asset1, :amount_asset2, :created_by_id,
            :updated_by_id, :created_at, :updated_at, :corrected, :canceled,
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
	//Gets UUID for the current authenticated user
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid user ID")
	}

	//Generates ordernumber
	OrderNumber := orderutils.GenerateOrderNumber()

	//Fetches house account
	houseAccountID, err := accountutils.GetHouseAccount(c)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	houseAccountUUID, err := uuid.Parse(houseAccountID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	clientTransactionID := uuid.New()
	houseTransactionID := uuid.New()

	clientTransaction := *transactionData
	clientTransaction.Id = clientTransactionID
	clientTransaction.CreatedById = userUUID
	clientTransaction.UpdatedById = userUUID
	clientTransaction.CreatedAt = time.Now()
	clientTransaction.UpdatedAt = time.Now()
	clientTransaction.Corrected = false
	clientTransaction.Corrected = false
	clientTransaction.OrderNumber = OrderNumber

	houseTransaction := clientTransaction
	houseTransaction.Id = houseTransactionID
	houseTransaction.TransactionOwnerAccountId = houseAccountUUID

	/* CHECKS THE BALANCE OF THE CUSTOMER ACCOUNT */
	currentBalance, err := s.transactionRepo.GetAccountBalance(clientTransaction.TransactionOwnerAccountId)
	if err != nil {
		fmt.Printf("error in currentBalance: %v")
		return uuid.Nil, uuid.Nil, err
	}

	currentAvailableBalance, err := s.transactionRepo.GetAccountAvailableBalance(clientTransaction.TransactionOwnerAccountId)
	if err != nil {
		fmt.Printf("error in currentAvailableBalance: %v")
		return uuid.Nil, uuid.Nil, err
	}

	/* BALANCE CHECK FOR HOUSE */

	currentBalanceHouse, err := s.transactionRepo.GetAccountBalance(houseAccountUUID)
	if err != nil {
		fmt.Printf("error in currentBalanceHouse: %v")
		return uuid.Nil, uuid.Nil, err
	}

	currentAvailableBalanceHouse, err := s.transactionRepo.GetAccountAvailableBalance(houseAccountUUID)
	if err != nil {
		fmt.Printf("error in currentAvailableBalanceHouse: %v")
		return uuid.Nil, uuid.Nil, err
	}

	/* CALCULATES BALANCE FOR CUSTOMER AND HOUSE */
	var newBalance, availableBalance, newBalanceHouse, availableBalanceHouse float64
	if clientTransaction.CashAmount != nil {
		newBalance = currentBalance + *clientTransaction.CashAmount
		availableBalance = currentAvailableBalance + *clientTransaction.CashAmount
		newBalanceHouse = currentBalanceHouse + *clientTransaction.CashAmount // Assuming house loses this amount
		availableBalanceHouse = currentAvailableBalanceHouse + *clientTransaction.CashAmount
	} else {
		newBalance = currentBalance
		availableBalance = currentAvailableBalance
		newBalanceHouse = currentBalanceHouse
		availableBalanceHouse = currentAvailableBalanceHouse
	}

	err = s.transactionRepo.UpdateAccountBalance(clientTransaction.CashAccountId, newBalance, availableBalance)
	if err != nil {
		fmt.Printf("error in UpdateAccountBalance: %v", err)
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.UpdateAccountBalance(houseAccountUUID, newBalanceHouse, availableBalanceHouse)
	if err != nil {
		fmt.Printf("error in UpdateAccountBalance: %v", err)
		return uuid.Nil, uuid.Nil, err
	}

	// Insert into database using repository
	err = s.transactionRepo.InsertTransaction(&clientTransaction)
	if err != nil {
		fmt.Printf("error in ClientTransaction: %v", err)
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.InsertTransaction(&houseTransaction)
	if err != nil {
		fmt.Printf("error in InsertCreditTransaction: %v", err)
		return uuid.Nil, uuid.Nil, err
	}

	return clientTransactionID, houseAccountUUID, nil
}

func (s *TransactionService) CreateInstrumentPurchaseTransaction(c *gin.Context, accountID uuid.UUID, userID string, transactionData *models.Transaction, transactionInstrumentData *models.Transaction) (uuid.UUID, uuid.UUID, error) {

	/*	userUUID, err := uuid.Parse(userID)
		if err != nil {
			return uuid.Nil, uuid.Nil, errors.New("invalid user ID")
		} */

	//Fetches house account
	houseAccountID, err := accountutils.GetHouseAccount(c)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	fmt.Printf("House account: %v", houseAccountID)
	houseAccountUUID, err := uuid.Parse(houseAccountID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	fmt.Println(transactionData)

	clientCashTransaction := *transactionData
	clientCashTransaction.AssetQuantity = nil
	clientCashTransaction.Id = uuid.New()

	houseCashTransaction := *transactionData
	houseCashTransaction.TransactionOwnerAccountId = houseAccountUUID
	houseCashTransaction.Id = uuid.New()
	clientCashTransaction.AssetQuantity = nil

	clientInstrumentTransaction := *transactionInstrumentData
	clientInstrumentTransaction.CashAmount = nil
	clientInstrumentTransaction.Id = uuid.New()

	houseInstrumentTransaction := *transactionInstrumentData
	houseInstrumentTransaction.TransactionOwnerAccountId = houseAccountUUID
	houseInstrumentTransaction.Id = uuid.New()
	clientInstrumentTransaction.CashAmount = nil
	clientInstrumentTransaction.Id = uuid.New()

	/* CHECKS THE BALANCE OF THE CUSTOMER ACCOUNT */
	currentBalance, err := s.transactionRepo.GetAccountBalance(clientCashTransaction.TransactionOwnerAccountId)
	if err != nil {
		fmt.Printf("error in currentBalance: %v")
		return uuid.Nil, uuid.Nil, err
	}

	currentAvailableBalance, err := s.transactionRepo.GetAccountAvailableBalance(clientCashTransaction.TransactionOwnerAccountId)
	if err != nil {
		fmt.Printf("error in currentAvailableBalance: %v")
		return uuid.Nil, uuid.Nil, err
	}

	/* BALANCE CHECK FOR HOUSE */

	currentBalanceHouse, err := s.transactionRepo.GetAccountBalance(houseAccountUUID)
	if err != nil {
		fmt.Printf("error in currentBalanceHouse: %v")
		return uuid.Nil, uuid.Nil, err
	}

	currentAvailableBalanceHouse, err := s.transactionRepo.GetAccountAvailableBalance(houseAccountUUID)
	if err != nil {
		fmt.Printf("error in currentAvailableBalanceHouse: %v")
		return uuid.Nil, uuid.Nil, err
	}

	/* CALCULATES BALANCE FOR CUSTOMER AND HOUSE */
	var newBalance, availableBalance, newBalanceHouse, availableBalanceHouse float64
	fmt.Printf("cash transaction: %v", clientCashTransaction.CashAmount)
	if clientCashTransaction.CashAmount != nil {
		newBalance = currentBalance + *clientCashTransaction.CashAmount
		availableBalance = currentAvailableBalance + *clientCashTransaction.CashAmount
		newBalanceHouse = currentBalanceHouse + *clientCashTransaction.CashAmount // Assuming house loses this amount
		availableBalanceHouse = currentAvailableBalanceHouse + *clientCashTransaction.CashAmount
	} else {
		newBalance = currentBalance
		availableBalance = currentAvailableBalance
		newBalanceHouse = currentBalanceHouse
		availableBalanceHouse = currentAvailableBalanceHouse
	}

	fmt.Printf("New balance: %v", newBalance)
	fmt.Printf("available balance: %v", availableBalance)
	fmt.Printf("new balance house: %v", newBalanceHouse)
	fmt.Printf("available balance house: %v", availableBalanceHouse)

	err = s.transactionRepo.UpdateAccountBalance(clientCashTransaction.CashAccountId, availableBalance, availableBalance)
	if err != nil {
		fmt.Printf("error in UpdateAccountBalance: %v", err)
		return uuid.Nil, uuid.Nil, err
	}

	err = s.transactionRepo.UpdateAccountBalance(houseAccountUUID, availableBalanceHouse, availableBalanceHouse)
	if err != nil {
		fmt.Printf("error in UpdateAccountBalance: %v", err)
		return uuid.Nil, uuid.Nil, err
	}

	// Insert into database using repository
	err = s.transactionRepo.InsertTransaction(&clientCashTransaction)
	if err != nil {
		fmt.Printf("error in insert client cash: %v", err)
		return uuid.Nil, uuid.Nil, err
	}

	err = s.transactionRepo.InsertTransaction(&houseCashTransaction)
	if err != nil {
		fmt.Printf("error in insert house cash: %v", err)
		return uuid.Nil, uuid.Nil, err
	}
	//Inserts instrument transactions
	err = s.transactionRepo.InsertTransaction(&clientInstrumentTransaction)
	if err != nil {
		fmt.Printf("error in insert client instrument: %v", err)
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.InsertTransaction(&houseInstrumentTransaction)
	if err != nil {
		fmt.Printf("error in insert house instrument: %v", err)
		return uuid.Nil, uuid.Nil, err
	}

	//err = s.transactionRepo.InsertTransaction(creditInstrumentTransaction)

	return clientCashTransaction.Id, houseCashTransaction.Id, nil

}

func (s *TransactionService) CreateWithdrawal(c *gin.Context, userID string, transactionData *models.Transaction) (uuid.UUID, uuid.UUID, error) {
	//Gets UUID for the current authenticated user
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid user ID")
	}

	//Generates ordernumber
	OrderNumber := orderutils.GenerateOrderNumber()

	//Fetches house account
	houseAccountID, err := accountutils.GetHouseAccount(c)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	houseAccountUUID, err := uuid.Parse(houseAccountID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	clientTransactionID := uuid.New()
	houseTransactionID := uuid.New()

	clientTransaction := *transactionData
	clientTransaction.Id = clientTransactionID
	clientTransaction.CreatedById = userUUID
	clientTransaction.UpdatedById = userUUID
	clientTransaction.CreatedAt = time.Now()
	clientTransaction.UpdatedAt = time.Now()
	clientTransaction.Corrected = false
	clientTransaction.Corrected = false
	clientTransaction.OrderNumber = OrderNumber

	houseTransaction := clientTransaction
	houseTransaction.Id = houseTransactionID
	houseTransaction.TransactionOwnerAccountId = houseAccountUUID

	/* CHECKS THE BALANCE OF THE CUSTOMER ACCOUNT */
	currentBalance, err := s.transactionRepo.GetAccountBalance(clientTransaction.TransactionOwnerAccountId)
	if err != nil {
		fmt.Printf("error in currentBalance: %v")
		return uuid.Nil, uuid.Nil, err
	}

	currentAvailableBalance, err := s.transactionRepo.GetAccountAvailableBalance(clientTransaction.TransactionOwnerAccountId)
	if err != nil {
		fmt.Printf("error in currentAvailableBalance: %v")
		return uuid.Nil, uuid.Nil, err
	}

	/* BALANCE CHECK FOR HOUSE */

	currentBalanceHouse, err := s.transactionRepo.GetAccountBalance(houseAccountUUID)
	if err != nil {
		fmt.Printf("error in currentBalanceHouse: %v")
		return uuid.Nil, uuid.Nil, err
	}

	currentAvailableBalanceHouse, err := s.transactionRepo.GetAccountAvailableBalance(houseAccountUUID)
	if err != nil {
		fmt.Printf("error in currentAvailableBalanceHouse: %v")
		return uuid.Nil, uuid.Nil, err
	}
	fmt.Printf("currentAvailableBalanceHouse before calculation: ", currentAvailableBalance)
	/* CALCULATES BALANCE FOR CUSTOMER AND HOUSE */
	var newBalance, availableBalance, newBalanceHouse, availableBalanceHouse float64
	if clientTransaction.CashAmount != nil {
		newBalance = currentBalance - *clientTransaction.CashAmount
		availableBalance = currentAvailableBalance - *clientTransaction.CashAmount
		newBalanceHouse = currentBalanceHouse - *clientTransaction.CashAmount // Assuming house loses this amount
		availableBalanceHouse = currentAvailableBalanceHouse - *clientTransaction.CashAmount
	} else {
		newBalance = currentBalance
		availableBalance = currentAvailableBalance
		newBalanceHouse = currentBalanceHouse
		availableBalanceHouse = currentAvailableBalanceHouse
	}

	var cashAmountStr string
	if clientTransaction.CashAmount != nil {
		cashAmountStr = strconv.FormatFloat(*houseTransaction.CashAmount, 'f', 2, 64)
	} else {
		cashAmountStr = "default_value" // Or any other appropriate handling
	}

	fmt.Println(cashAmountStr)
	fmt.Println(availableBalance)
	if currentAvailableBalance < *clientTransaction.CashAmount {
		return uuid.Nil, uuid.Nil, fmt.Errorf("insufficient funds in customer's account")
	}

	fmt.Println(cashAmountStr)
	fmt.Println(availableBalanceHouse)
	if currentAvailableBalanceHouse < *houseTransaction.CashAmount {
		return uuid.Nil, uuid.Nil, fmt.Errorf("insufficient funds in house account")
	}

	err = s.transactionRepo.UpdateAccountBalance(clientTransaction.CashAccountId, newBalance, availableBalance)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.UpdateAccountBalance(houseAccountUUID, newBalanceHouse, availableBalanceHouse)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// Insert into database using repository
	err = s.transactionRepo.InsertTransaction(&clientTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	err = s.transactionRepo.InsertTransaction(&houseTransaction)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return clientTransactionID, houseTransactionID, nil
}
