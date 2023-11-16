package handlers

import (
	"net/http"
	"thyra/internal/assets/models"
	"thyra/internal/assets/repositories"
	"thyra/internal/assets/services"
	"thyra/internal/common/db"
	"thyra/internal/transactions/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

func GetAllOrdersHandler(c *gin.Context) {
	// Extract the authenticated user's ID and role from context.
	_, authUserRole, isAuthenticated := utils.GetAuthenticatedUser(c)
	if !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated or role not found"})
		return
	}

	// Ensure the authenticated user is authorized to fetch orders.
	if authUserRole != "admin" && authUserRole != "order_manager" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized. Admin or Order Manager access required"})
		return
	}

	// Database connection
	db := db.GetConnection(c)
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection issue"})
		return
	}

	// Fetching all orders with related account and asset details from the database.
	query := `
	SELECT o.*, a.account_number, asst.instrument_name, asst.instrument_type 
	FROM thyrasec.orders o
	JOIN thyrasec.accounts a ON o.account_id = a.id
	JOIN thyrasec.assets asst ON o.asset_id = asst.id
	`
	var orders []models.OrderWithDetails // Define a new struct that includes the additional fields
	if err := db.Select(&orders, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders", "details": err.Error()})
		return
	}

	// Handling empty orders list.
	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No orders found"})
		return
	}

	// Sending the response with the orders and their details.
	c.JSON(http.StatusOK, orders)
}

/* Handler for creating buy order*/
func CreateBuyOrderHandler(c *gin.Context) {
	var newOrder models.Order

	// Parse the request body into the newOrder struct.
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	totalAmountDecimal := decimal.NewFromFloat(newOrder.TotalAmount)

	// Generate a new UUID for the order.
	newOrder.ID = uuid.New()
	newOrder.Status = models.StatusCreated

	// Get a database connection.
	dbConn := db.GetConnection(c)
	if dbConn == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection issue"})
		return
	}

	// Start a database transaction.
	tx, err := dbConn.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction", "details": err.Error()})
		return
	}

	// Insert the order into the database.
	stmt := `INSERT INTO thyrasec.orders (id, account_id, asset_id, order_type, quantity, price_per_unit, total_amount, status, created_at, updated_at) 
              VALUES (:id, :account_id, :asset_id, :order_type, :quantity, :price_per_unit, :total_amount, :status, NOW(), NOW())`
	_, err = tx.NamedExec(stmt, newOrder)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}

	if err := services.CheckAndReserveCash(tx, newOrder.AccountID, totalAmountDecimal); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reserve cash", "details": err.Error()})
		return
	}

	// Insert cash reservation for the order
	reservedUntil := time.Now().Add(24 * time.Hour)
	if err := repositories.InsertCashReservation(tx, newOrder, reservedUntil); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cash reservation", "details": err.Error()})
		return
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction", "details": err.Error()})
		return
	}
	// Respond with the created order.
	c.JSON(http.StatusCreated, newOrder)
}

/* Creates a sell order */
func CreateSellOrderHandler(c *gin.Context) {
	var newOrder models.Order

	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	newOrder.ID = uuid.New()
	newOrder.Status = models.StatusCreated

	dbConn := db.GetConnection(c)
	if dbConn == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection issue"})
		return
	}

	stmt := `INSERT INTO thyrasec.orders (id, account_id, asset_id, order_type, quantity, price_per_unit, total_amount, status, created_at, updated_at) 
              VALUES (:id, :account_id, :asset_id, :order_type, :quantity, :price_per_unit, :total_amount, :status, NOW(), NOW())`

	_, err := dbConn.NamedExec(stmt, newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}

	if err := services.CreateReservationAndDeductHoldings(dbConn, newOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process sell order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newOrder)
}

func ConfirmOrderHandler(c *gin.Context) {
	orderID := c.Param("orderId")

	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	// Check if the order exists and is in a state that can be confirmed
	order, err := repositories.GetOrder(sqlxDB, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
		return
	}

	if order.Status != "created" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be confirmed in its current state"})
		return
	}

	// Perform specific actions for confirming an order
	if err := services.ConfirmOrder(sqlxDB, orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order confirmed successfully"})
}

func ExecuteOrderHandler(c *gin.Context) {
	orderID := c.Param("orderId")

	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	// Check if the order exists and is in a state that can be executed
	order, err := repositories.GetOrder(sqlxDB, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
		return
	}

	if order.Status != "confirmed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be executed in its current state"})
		return
	}

	// Perform specific actions for executing an order
	if err := services.ExecuteOrder(sqlxDB, orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order executed successfully"})
}

func SettlementHandler(c *gin.Context) {
	orderID := c.Param("orderId")

	db, _ := c.Get("db")
	sqlxDB, _ := db.(*sqlx.DB)

	// Get the order from the database
	order, err := repositories.GetOrder(sqlxDB, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
		return
	}

	// Ensure that the order is in the correct state
	if order.Status != models.StatusExecuted {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Order cannot be settled in its current state", "error": err.Error()})
		return
	}

	// Release the reservation
	if err := repositories.ReleaseReservation(sqlxDB, orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to release reservation", "details": err.Error()})
		return
	}

	// Update account balance - modify as per your application logic
	if err := repositories.UpdateAccountBalance(sqlxDB, order.AccountID, order.TotalAmount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance", "details": err.Error()})
		return
	}

	// Insert into holding - create a holding object from order details
	holding := models.Holding{
		ID:        uuid.New(),
		AccountID: order.AccountID,
		AssetID:   order.AssetID,
		Quantity:  order.Quantity,
	}

	if err := repositories.InsertHolding(sqlxDB, holding); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into holding", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order settled successfully"})
}
