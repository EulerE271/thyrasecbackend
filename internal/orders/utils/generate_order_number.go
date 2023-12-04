package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderNumber() string {
	rand.Seed(time.Now().UnixNano())

	// Generate a random number with the desired number of digits.
	// In this case, we're generating a 6-digit number to follow the initial '4'.
	randomNumbers := rand.Intn(99999-10000) + 10000
	orderNumber := fmt.Sprintf("4%d", randomNumbers)

	return orderNumber
}
