package services

import (
	"errors"
	"fmt"
	"log"
	"thyra/internal/orders/models"
	"thyra/internal/orders/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

/* Checks and reserves cash when buying an instrument */
func CheckAndReserveCash(tx *sqlx.Tx, accountID uuid.UUID, amount decimal.Decimal) error {
	// Check if there is sufficient available cash
	sufficientCash, err := repositories.CheckAvailableCash(tx, accountID, amount)
	if err != nil {
		return err
	}
	if !sufficientCash {
		return errors.New("insufficient funds")
	}

	// Reserve the cash
	err = repositories.ReserveCash(tx, accountID, amount)
	if err != nil {
		return err
	}

	return nil
}

/*Creates a reservation and updates the holding table when selling an asset */
func CreateReservationAndDeductHoldings(db *sqlx.DB, order models.Order) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Check available holdings
	availableQuantity, err := repositories.CheckHoldings(tx, order.AccountID, order.AssetID)
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
	if err := repositories.InsertReservation(tx, order, reservedUntil); err != nil {
		tx.Rollback()
		return err
	}

	// Deduct from holdings
	if err := repositories.DeductHoldings(tx, order.AccountID, order.AssetID, int(order.Quantity)); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func ConfirmOrder(db *sqlx.DB, orderID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Retrieve the order to confirm
	order, err := repositories.GetOrder(db, orderID)
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
	orderTypeName, err := repositories.GetOrderType(tx, order.OrderType)
	if err != nil {
		log.Fatalf("error in orderTypename: %v", err)
		tx.Rollback()
		return err
	}

	// Perform actions specific to confirming the order
	if orderTypeName == "order_type_buy" {
		totalAmountDecimal := decimal.NewFromFloat(order.TotalAmount)
		_, err := repositories.CheckAvailableCash(tx, order.AccountID, totalAmountDecimal)
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
		availableQuantity, err := repositories.CheckHoldings(tx, order.AccountID, order.AssetID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if float64(availableQuantity) < order.Quantity {
			tx.Rollback()
			return errors.New("insufficient holdings for sell order")
		}
		if err := CreateReservationAndDeductHoldings(db, *order); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		tx.Rollback()
		return errors.New("unknown order type")
	}

	// Update order status to confirmed
	if err := repositories.UpdateOrderStatus(tx, orderID, models.StatusConfirmed); err != nil {
		log.Fatalf("error in updateorderstatus: %v", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func ExecuteOrder(db *sqlx.DB, orderID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Retrieve the order to confirm
	order, err := repositories.GetOrder(db, orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if the order is in a state that can be confirmed
	if order.Status != models.StatusConfirmed {
		tx.Rollback()
		return errors.New("order cannot be executed in its current state")
	}

	orderTypeName, err := repositories.GetOrderType(tx, order.OrderType)
	if err != nil {
		tx.Rollback()
		return err
	}

	if orderTypeName == "order_type_buy" {
		totalAmountDecimal := decimal.NewFromFloat(order.TotalAmount)
		_, err := repositories.CheckAvailableCash(tx, order.AccountID, totalAmountDecimal)
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
		availableQuantity, err := repositories.CheckHoldings(tx, order.AccountID, order.AssetID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if float64(availableQuantity) < order.Quantity {
			tx.Rollback()
			return errors.New("insufficient holdings for sell order")
		}
		if err := CreateReservationAndDeductHoldings(db, *order); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		tx.Rollback()
		return errors.New("unknown order type")
	}

	// Update order status to confirmed
	if err := repositories.UpdateOrderStatus(tx, orderID, models.StatusExecuted); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
