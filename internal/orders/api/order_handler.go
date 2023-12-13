package handlers

import (
	"fmt"
	"log"
	"net/http"
	accountutils "thyra/internal/accounts/utils"
	"thyra/internal/orders/models"
	"thyra/internal/orders/services" // using alias
	orderutils "thyra/internal/orders/utils"
	positionmodel "thyra/internal/positions/models"
	transactionmodels "thyra/internal/transactions/models"
	transactionservice "thyra/internal/transactions/services"
	authutils "thyra/internal/users/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderHandler struct {
	Service *services.OrdersService
}

func NewOrderHandler(Service *services.OrdersService, transactionservice *transactionservice.TransactionService) *OrderHandler {
	return &OrderHandler{Service: Service}
}

func GetAllOrdersHandler(Service *services.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, authUserRole, isAuthenticated := authutils.GetAuthenticatedUser(c)
		if !isAuthenticated {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated or role not found"})
			return
		}

		if authUserRole != "admin" && authUserRole != "order_manager" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized. Admin or Order Manager access required"})
			return
		}

		orders, err := Service.GetAllOrders()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders", "details": err.Error()})
			return
		}

		if len(orders) == 0 {
			c.JSON(http.StatusNoContent, gin.H{"message": "No orders found"})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

func CreateBuyOrderHandler(Service services.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newOrder models.Order
		if err := c.ShouldBindJSON(&newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		createdOrder, err := Service.CreateBuyOrder(newOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdOrder)
	}
}

func CreateSellOrderHandler(Service services.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newOrder models.Order
		if err := c.ShouldBindJSON(&newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		createdOrder, err := Service.CreateSellOrder(newOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdOrder)
	}
}

func ConfirmOrderHandler(Service services.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("orderId")

		db, _ := c.Get("db")
		sqlxDB, _ := db.(*sqlx.DB)

		// Check if the order exists and is in a state that can be confirmed
		order, err := Service.GetOrder(orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
			return
		}

		if order.Status != "created" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be confirmed in its current state"})
			return
		}

		// Perform specific actions for confirming an order
		if err := Service.ConfirmOrder(sqlxDB, orderID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm order", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order confirmed successfully"})
	}
}

func ExecuteOrderHandler(Service services.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {

		orderID := c.Param("orderId")

		db, _ := c.Get("db")
		sqlxDB, _ := db.(*sqlx.DB)

		// Check if the order exists and is in a state that can be executed
		order, err := Service.GetOrder(orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
			return
		}

		if order.Status != "confirmed" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be executed in its current state"})
			return
		}

		// Perform specific actions for executing an order
		if err := Service.ExecuteOrder(sqlxDB, orderID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute order", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order executed successfully"})
	}
}

type SettlementRequest struct {
	SettledQuantity float64    `json:"quantity"`
	SettledAmount   float64    `json:"amount"`
	Comment         string     `json:"comment"`
	TradeDate       *time.Time `json:"tradeDate"`
	SettlementDate  *time.Time `json:"settlementDate"`
}

func SettlementBuyHandler(Service services.OrdersService, transactionservice transactionservice.TransactionService) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error userIdInterface": "UserID not found in context"})
			return
		}

		userIDStr, ok := userIDInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error userIDStr": "UserID is not a string"})
			return
		}

		// Parse the userID string as UUID (only if needed later in the code)
		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is not a valid UUID", "details": err.Error()})
			return
		}

		orderID := c.Param("orderId")

		db, exists := c.Get("db")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
			return
		}
		sqlxDB, ok := db.(*sqlx.DB)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection"})
			return
		}

		// Get the order from the database
		order, err := Service.GetOrder(orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
			return
		}

		// Ensure that the order is in the correct state
		if order.Status != models.StatusExecuted {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Order cannot be settled in its current state"})
			return
		}

		var settlementRequest SettlementRequest
		if err := c.BindJSON(&settlementRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		assetType, err := Service.GetAssetType(order.AssetID)
		if err != nil {
			log.Fatalf("Error fetching assetType: %v", err)
		}

		/*fmt.Println(orderType)
		/*transactionType, err := repositories.GetOrderType(sqlxDB, orderType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong fetching orderType"})
			return
		} */

		transactionType, err := Service.GetTransactionTypeByOrderTypeID(order.OrderType)
		if err != nil {
			log.Fatalf("Error fetching transactionTyper: %v", err)
		}
		orderNumber := orderutils.GenerateOrderNumber()

		/*transactionRepo := transactionRepository.NewTransactionRepository(sqlxDB)
		transactionService := transactionServices.NewTransactionService(transactionRepo)*/

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
			CreatedById:               userUUID,
			UpdatedById:               userUUID,
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

		houseAccount, err := accountutils.GetHouseAccount(sqlxDB)
		if err != nil {
			log.Fatalf("Error getting house account: %v", err)
			return
		}

		// Parse the houseAccount string into UUID
		houseAccountUUID, err := uuid.Parse(houseAccount)
		if err != nil {
			log.Fatalf("Error parsing house account UUID: %v", err)
			return
		}

		clientInstrumentTransaction.TransactionOwnerAccountId = houseAccountUUID

		_, _, err = transactionservice.CreateInstrumentPurchaseTransaction(c, order.AccountID, userUUID.String(), &clientCashTransaction, &clientInstrumentTransaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction", "details": err.Error()})
			return
		}

		fmt.Println(settlementRequest.SettledQuantity)
		err = Service.UpdateOrder(orderID, settlementRequest.SettledQuantity, settlementRequest.SettledAmount, "settled", settlementRequest.TradeDate, settlementRequest.SettlementDate, settlementRequest.Comment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order", "details": err.Error()})
			return
		}

		// Release the reservation
		err = Service.ReleaseReservation(orderID, houseAccount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to release reservation", "details": err.Error()})
			return
		}

		// Update account balance
		/*fmt.Println("order total amount: %v", order.TotalAmount) */
		fmt.Println("order accountID: %v", order.AccountID)
		err = Service.UpdateAccountBalance(order.AccountID, order.TotalAmount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance", "details": err.Error()})
			return
		}

		// Insert into holding
		holding := positionmodel.Holding{
			ID:        uuid.New(),
			AccountID: order.AccountID,
			AssetID:   order.AssetID,
			Quantity:  order.Quantity,
			// Populate other fields as needed
		}

		err = Service.InsertHolding(holding)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into holding", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order settled successfully"})
	}

}

func GetOrderTypeByName(Service services.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {

		name := c.Query("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order type name is required"})
			return
		}

		orderType, err := Service.GetOrderTypeByName(name)
		if err != nil {
			// Handle errors, such as not finding the order type
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"order_type_id": orderType})

	}
}

func GetOrderTypeByID(Service services.OrdersService) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID  is required"})
		}

		UUID := uuid.MustParse(id)

		orderType, err := Service.GetOrderType(UUID)
		if err != nil {
			// Handle errors, such as not finding the order type
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"order_type_name": orderType})
	}
}
