package accountno

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenerateAccountNumber(uuidStr string) string {
	if len(uuidStr) < 4 { // Ensure we have at least 4 characters
		fmt.Println("UUID is too short")
		return ""
	}

	// Extract the first 4 characters of the UUID
	firstFour := uuidStr[:4]

	// Convert these characters to an integer using base 16 (hexadecimal)
	firstDigit, err := strconv.ParseInt(firstFour, 16, 64)
	if err != nil {
		fmt.Println("Failed to parse UUID substring:", err)
		return ""
	}

	// Use modulo to ensure the number is between 0 and 9
	firstDigit = firstDigit % 10

	rand.Seed(time.Now().UnixNano())

	randomNumbers := rand.Intn(1000000)
	accountNumber := fmt.Sprintf("%d%06d", firstDigit, randomNumbers)

	return accountNumber
}
