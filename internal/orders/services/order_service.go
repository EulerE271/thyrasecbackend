package services

import (
	"errors"
	"fmt"
	"log"
	accountutils "thyra/internal/accounts/utils"
	"thyra/internal/orders/models"
	"thyra/internal/orders/repositories"
	"thyra/internal/orders/utils"
	positionsmodel "thyra/internal/positions/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type OrdersService struct {
	db   *sqlx.DB
	repo *repositories.OrdersRepository
}

func NewOrdersService(db *sqlx.DB, repo *repositories.OrdersRepository) *OrdersService {
	return &OrdersService{db: db, repo: repo}
}

/* Checks and reserves cash when buying an instrument */
func (s *OrdersService) CheckAndReserveCash(tx *sqlx.Tx, accountID uuid.UUID, amount decimal.Decimal) error {
	// Check if there is sufficient available cash
	sufficientCash, err := s.repo.CheckAvailableCash(tx, accountID, amount)
	if err != nil {
		return err
	}
	if !sufficientCash {
		return errors.New("insufficient funds")
	}

	// Reserve the cash
	err = s.repo.ReserveCash(tx, accountID, amount)
	if err != nil {
		return err
	}

	return nil
}

/*Creates a reservation and updates the holding table when selling an asset */
func (s *OrdersService) CreateReservationAndDeductHoldings(db *sqlx.DB, order models.Order) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Check available holdings
	availableQuantity, err := s.repo.CheckHoldings(tx, order.AccountID, order.AssetID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if float64(availableQuantity) < order.Quantity {
		tx.Rollback()
		return errors.New("insufficient holdings")
	}

	// Create reservation (assuming reserved duration is provided)
	reservedUntil := time.Now().Add(1000 * time.Hour)
	if err := s.repo.InsertReservation(tx, order, reservedUntil); err != nil {
		tx.Rollback()
		return err
	}

	// Deduct from holdings
	if err := s.repo.DeductHolding(tx, order.AccountID, order.AssetID, order.Quantity); err != nil {
		fmt.Printf("Error in deduct holding: %v", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *OrdersService) ConfirmOrder(db *sqlx.DB, orderID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Retrieve the order to confirm
	order, err := s.repo.GetOrder(tx, orderID)
	if err != nil {
		log.Fatalf("error in order: %v", err)
		tx.Rollback()
		return err
	}

	// Check if the order is in a state that can be confirmed
	if order.Status != models.StatusCreated {
		tx.Rollback()
		return errors.New("order cannot be confirmed in its current state")
	}
	fmt.Printf("ordertype: %v", order.OrderType)
	orderTypeName, err := s.repo.GetOrderType(tx, order.OrderType)
	if err != nil {
		log.Fatalf("error in orderTypename: %v", err)
		tx.Rollback()
		return err
	}

	// Perform actions specific to confirming the order
	if orderTypeName == "order_type_buy" {
		totalAmountDecimal := decimal.NewFromFloat(order.TotalAmount)
		_, err := s.repo.CheckAvailableCash(tx, order.AccountID, totalAmountDecimal)
		if err != nil {
			log.Fatalf("error in checkAvailableCash: %v", err)
			tx.Rollback()
			return err
		}

		/*if !sufficientFunds {
			tx.Rollback()
			return errors.New("insufficient funds to execute buy order")
		} */
	} else if orderTypeName == "order_type_sell" {
		// Logic to reserve assets for a sell order
		availableQuantity, err := s.repo.CheckHoldings(tx, order.AccountID, order.AssetID)
		if err != nil {
			fmt.Printf("error in check holdings: %v", err)
			tx.Rollback()
			return err
		}
		if float64(availableQuantity) < order.Quantity {
			tx.Rollback()
			return errors.New("insufficient holdings for sell order")
		}
	} else {
		tx.Rollback()
		return errors.New("unknown order type")
	}

	// Update order status to confirmed
	if err := s.repo.UpdateOrderStatus(tx, orderID, models.StatusConfirmed); err != nil {
		log.Fatalf("error in updateorderstatus: %v", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *OrdersService) ExecuteOrder(db *sqlx.DB, orderID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Retrieve the order to confirm
	order, err := s.repo.GetOrder(tx, orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if the order is in a state that can be confirmed
	if order.Status != models.StatusConfirmed {
		tx.Rollback()
		return errors.New("order cannot be executed in its current state")
	}

	orderTypeName, err := s.repo.GetOrderType(tx, order.OrderType)
	if err != nil {
		tx.Rollback()
		return err
	}

	if orderTypeName == "order_type_buy" {
		totalAmountDecimal := decimal.NewFromFloat(order.TotalAmount)
		_, err := s.repo.CheckAvailableCash(tx, order.AccountID, totalAmountDecimal)
		if err != nil {
			tx.Rollback()
			return err
		}
		/*if !sufficientFunds {
			tx.Rollback()
			return errors.New("insufficient funds to execute buy order")
		} */
	} else if orderTypeName == "order_type_sell" {
		// Logic to reserve assets for a sell order
		availableQuantity, err := s.repo.CheckHoldings(tx, order.AccountID, order.AssetID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if float64(availableQuantity) < order.Quantity {
			tx.Rollback()
			return errors.New("insufficient holdings for sell order")
		}
	} else {
		tx.Rollback()
		return errors.New("unknown order type")
	}

	// Update order status to confirmed
	if err := s.repo.UpdateOrderStatus(tx, orderID, models.StatusExecuted); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *OrdersService) GetAllOrders() ([]models.OrderWithDetails, error) {
	return s.repo.GetAllOrders()
}

func (s *OrdersService) CreateBuyOrder(newOrder models.Order) (models.Order, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return models.Order{}, err
	}

	// Set unique identifiers and status for the new order
	newOrder.ID = uuid.New()
	newOrder.OrderNumber = utils.GenerateOrderNumber()
	newOrder.Status = models.StatusCreated

	// Insert the order into the database using the repository
	if err := s.repo.InsertOrder(tx, newOrder); err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	// Example business logic for handling cash reservations
	totalAmountDecimal := decimal.NewFromFloat(newOrder.TotalAmount)
	if err := s.CheckAndReserveCash(tx, newOrder.AccountID, totalAmountDecimal); err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	houseAccount, err := accountutils.GetHouseAccount(tx) // Ensure this returns the house account ID
	if err != nil {
		tx.Rollback()
		log.Fatalf("error fetching house account: %v", err)
		return models.Order{}, err
	}
	houseAccountUUID := uuid.MustParse(houseAccount)

	if err := s.CheckAndReserveCash(tx, houseAccountUUID, totalAmountDecimal); err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	// Insert cash reservation for the order
	reservedUntil := time.Now().Add(24 * time.Hour)
	if err := s.repo.InsertCashReservation(tx, newOrder, reservedUntil); err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.Order{}, err
	}

	return newOrder, nil
}

func (s *OrdersService) CreateSellOrder(newOrder models.Order) (models.Order, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return models.Order{}, err
	}

	newOrder.ID = uuid.New()
	newOrder.Status = models.StatusCreated
	newOrder.OrderNumber = utils.GenerateOrderNumber()

	if err := s.repo.ReserveAsset(tx, newOrder.AccountID, newOrder.Quantity, newOrder.AssetID); err != nil {
		log.Println("Error reserving asset:", err)
		tx.Rollback()
		return models.Order{}, err
	}

	if err := s.repo.InsertOrder(tx, newOrder); err != nil {
		log.Println("Error reserving asset:", err)
		tx.Rollback()
		return models.Order{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Order{}, err
	}

	return newOrder, nil

}

func (s *OrdersService) GetOrder(orderID string) (models.Order, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return models.Order{}, err
	}

	newOrder, err := s.repo.GetOrder(tx, orderID)
	if err != nil {
		fmt.Println("Something went wrong fetching orders")
		return models.Order{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Order{}, err
	}

	return *newOrder, nil
}

func (s *OrdersService) GetAssetType(assetID uuid.UUID) (uuid.UUID, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return uuid.Nil, err
	}

	assetTypeID, err := s.repo.GetAssetType(tx, assetID)
	if err != nil {
		fmt.Println("Error fetching Asset Type")
		return uuid.Nil, err
	}
	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return assetTypeID, nil

}

func (s *OrdersService) GetTransactionTypeByOrderTypeID(orderTypeID uuid.UUID) (uuid.UUID, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return uuid.Nil, err
	}

	transactionType, err := s.repo.GetTransactionTypeByOrderTypeID(tx, orderTypeID)
	if err != nil {
		fmt.Println("Something went wrong fetching transaction type")
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return transactionType, err
}

func (s *OrdersService) UpdateOrder(orderID string, settledQuantity float64, settledAmount float64, status string, tradeDate *time.Time, settlementDate *time.Time, comment string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	// Defer a rollback in case something fails. The rollback will be ignored if the commit is successful.
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		}
	}()

	if err = s.repo.UpdateOrder(tx, orderID, settledQuantity, settledAmount, status, tradeDate, settlementDate, comment); err != nil {
		// Consider using a more robust logging mechanism here.
		fmt.Println("something went wrong updating order:", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *OrdersService) ReleaseReservation(orderID string, houseAccount string) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		}
	}()

	if err := s.repo.ReleaseReservation(tx, orderID, houseAccount); err != nil {
		fmt.Println("Something went wrong releasing the reservation:", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *OrdersService) UpdateAccountBalance(accountID uuid.UUID, balanceChange float64) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		}
	}()

	if err := s.repo.UpdateAccountBalance(tx, accountID, balanceChange); err != nil {
		fmt.Println("Something went wrong releasing the reservation:", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *OrdersService) InsertHolding(holding positionsmodel.Holding) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		}
	}()

	if err := s.repo.InsertHolding(tx, holding); err != nil {
		fmt.Println("Something went wrong releasing the reservation:", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *OrdersService) GetOrderTypeByName(name string) (uuid.UUID, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return uuid.Nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	orderTypeID, err := s.repo.GetOrderTypeByName(tx, name)
	if err != nil {
		fmt.Println("Something went wrong fetching order type: ", err)
		return uuid.Nil, err
	}

	if err = tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return orderTypeID, err

}

func (s *OrdersService) GetOrderType(ID uuid.UUID) (string, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return "", err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	orderTypeName, err := s.repo.GetOrderType(tx, ID)
	if err != nil {
		fmt.Println("Something went wrong fetching order type: %v", err)
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return orderTypeName, err

}
