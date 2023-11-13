package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHouseAccount(c *gin.Context) (string, error) {
	type Result struct {
		ID              string `json:"id"`
		AccountTypeName string `json:"account_type_name"`
	}

	token, exists := c.Get("token") // Assuming you have set the token in the context with key "token"
	if !exists {
		return "", fmt.Errorf("token not found in context")
	}

	tokenString, ok := token.(string)
	if !ok {
		return "", fmt.Errorf("token found in context is not a string")
	}

	cookie := http.Cookie{
		Name:  "token",
		Value: tokenString,
		// Set other cookie attributes here, if needed.
	}

	req, err := http.NewRequest("GET", "http://localhost:8083/v1/fetch/account/house", nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Add the cookie to the request
	req.AddCookie(&cookie)

	// ...

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to account service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-ok response status: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Log the body for debugging
	log.Printf("Response body: %s", body)

	// Define a struct to match the expected JSON response structure
	var response struct {
		HouseAccountID string `json:"house_account_id"`
	}

	// Unmarshal the body into the struct
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	// Check if the house account ID was found
	if response.HouseAccountID == "" {
		return "", fmt.Errorf("no House account was found")
	}

	// Return the house account ID
	return response.HouseAccountID, nil
}
