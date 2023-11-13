package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"thyra/internal/common/db"
	"thyra/internal/transactions/models"
	"thyra/internal/transactions/services"
	orderno "thyra/internal/transactions/services/ordernumber"
	"thyra/internal/transactions/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateDeposit(c *gin.Context) {
	database := db.GetConnection(c)
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		// Handle the missing userID case
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID not found in context"})
		return
	}

	// Assert the type of userID to be string
	userIDStr, ok := userIDInterface.(string)
	if !ok {
		// Handle the case where userID is not a string
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID is not a string"})
		return
	}

	// Parse the userID string as UUID
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		// Handle the error case where userID is not a valid UUID
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is not a valid UUID", "details": err.Error()})
		return
	}

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// handle error
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	newTransaction, err := utils.BindTransaction(c)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input in main transaction data", "error": err.Error()})
		return
	}

	OrderNumber := orderno.GenerateOrderNumber()
	houseAccountID, err := utils.GetHouseAccount(c)
	if err != nil {
		// Handle error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get House account", "details": err.Error()})
		return
	}

	houseccountUUID, err := uuid.Parse(houseAccountID)
	if err != nil {
		// Handle error
		return
	}

	debitTransactionID := uuid.New()
	creditTransactionID := uuid.New()

	// Debit transaction - already bound from the request
	creditTransaction := newTransaction
	creditTransaction.Id = creditTransactionID
	creditTransaction.CreatedAt = time.Now()
	creditTransaction.UpdatedAt = time.Now()
	creditTransaction.CreatedById = userUUID
	creditTransaction.UpdatedById = userUUID
	creditTransaction.OrderNumber = OrderNumber

	// Credit transaction - create a new instance
	debitTransaction := &models.Transaction{
		Id:                 debitTransactionID, // Assign the generated UUID
		Type:               creditTransaction.Type,
		CreatedById:        userUUID,
		UpdatedById:        userUUID,
		Asset1Id:           creditTransaction.Asset1Id,
		AmountAsset1:       creditTransaction.AmountAsset1,
		CreatedAt:          time.Now(), // Use time.Now() to set to current time
		UpdatedAt:          time.Now(),
		StatusTransaction:  creditTransaction.StatusTransaction,
		Comment:            creditTransaction.Comment,
		TransactionOwnerId: creditTransaction.TransactionOwnerId,
		AccountOwnerId:     houseccountUUID,
		AccountAsset1Id:    creditTransaction.AccountAsset2Id,
		AccountAsset2Id:    uuid.Nil,
		OrderNumber:        OrderNumber,
		Trade_date:         creditTransaction.Trade_date,
		Settlement_date:    creditTransaction.Settlement_date,
	}

	//Update to remove fields from incorrect mappings
	creditTransaction.AccountAsset1Id = uuid.Nil
	creditTransaction.AmountAsset1 = nil
	creditTransaction.Asset1Id = uuid.Nil

	/* END OF PARENT */

	var currentBalance float64
	err = database.Get(&currentBalance, "SELECT account_balance FROM thyrasec.accounts WHERE account_holder_id = $1", userUUID)
	if err != nil {
		// Handle error, e.g., account not found or DB query failed
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account balance", "details": err.Error()})
		return
	}

	if *creditTransaction.AmountAsset1 > 0 {
		currentBalance += float64(*creditTransaction.AmountAsset1)
	} else {
		// Handle invalid amount
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deposit amount"})
		return
	}

	_, err = database.Exec("UPDATE thyrasec.accounts SET account_balance = $1, updated_at = $2 WHERE account_holder_id = $3", currentBalance, time.Now(), userUUID)
	if err != nil {
		// Handle error, e.g., DB update failed
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance", "details": err.Error()})
		return
	}

	/*PERFORM BALANCE UPDATING LOGIC*/

	/*INSERTS INTO DATABASE */
	err = services.InsertParentTransaction(database, debitTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create debit parent transaction", "error": err.Error()})
		return
	}
	err = services.InsertParentTransaction(database, creditTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create credit parent transaction", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully", "Parent debitTransactionID": debitTransactionID, "Parent CreditTransactionID": creditTransactionID})

}

func CreateWithdrawal(c *gin.Context) {

	database := db.GetConnection(c)
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		// Handle the missing userID case
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID not found in context"})
		return
	}

	// Assert the type of userID to be string
	userIDStr, ok := userIDInterface.(string)
	if !ok {
		// Handle the case where userID is not a string
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID is not a string"})
		return
	}

	// Parse the userID string as UUID
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		// Handle the error case where userID is not a valid UUID
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is not a valid UUID", "details": err.Error()})
		return
	}

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// handle error
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	newTransaction, err := utils.BindTransaction(c)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input in main transaction data", "error": err.Error()})
		return
	}

	OrderNumber := orderno.GenerateOrderNumber()
	houseAccountID, err := utils.GetHouseAccount(c)
	if err != nil {
		// Handle error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get House account", "details": err.Error()})
		return
	}

	houseccountUUID, err := uuid.Parse(houseAccountID)
	if err != nil {
		// Handle error
		return
	}

	debitTransactionID := uuid.New()
	creditTransactionID := uuid.New()

	// Debit transaction - already bound from the request
	debitTransaction := newTransaction
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
		AccountOwnerId:     houseccountUUID,
		AccountAsset1Id:    debitTransaction.AccountAsset2Id,
		AccountAsset2Id:    uuid.Nil,
		OrderNumber:        OrderNumber,
		Trade_date:         debitTransaction.Trade_date,      // Set to zero time if not used
		Settlement_date:    debitTransaction.Settlement_date, // Set to zero time if not used
	}

	//Update to remove fields from incorrect mappings
	debitTransaction.AccountAsset1Id = uuid.Nil
	debitTransaction.AmountAsset1 = nil
	debitTransaction.Asset1Id = uuid.Nil

	/* END OF PARENT */

	/*gRCP calls for updating the account service */

	if creditTransaction.AmountAsset1 != nil && *creditTransaction.AmountAsset1 > 0 {
		*creditTransaction.AmountAsset1 = -*creditTransaction.AmountAsset1
	}
	if creditTransaction.AmountAsset2 != nil && *creditTransaction.AmountAsset2 > 0 {
		*creditTransaction.AmountAsset2 = -*creditTransaction.AmountAsset2
	}

	// Adjust the sign of amounts for the debitTransaction (addition to house account)
	// If they are negative, make them positive
	if debitTransaction.AmountAsset1 != nil && *debitTransaction.AmountAsset1 > 0 {
		*debitTransaction.AmountAsset1 = -*debitTransaction.AmountAsset1
	}
	if debitTransaction.AmountAsset2 != nil && *debitTransaction.AmountAsset2 > 0 {
		*debitTransaction.AmountAsset2 = -*debitTransaction.AmountAsset2
	}

	withdrawalAmount := 0.0
	if creditTransaction.AmountAsset1 != nil {
		withdrawalAmount = float64(*creditTransaction.AmountAsset1)
	}
	fmt.Println(withdrawalAmount)

	/*INSERTS INTO DATABASE */
	//Inserts two parent transactions
	err = services.InsertParentTransaction(database, debitTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create debit parent transaction", "error": err.Error()})
		return
	}
	err = services.InsertParentTransaction(database, creditTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create credit parent transaction", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully", "Parent debitTransactionID": debitTransactionID, "Parent CreditTransactionID": creditTransactionID})

}
