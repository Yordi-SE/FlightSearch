package use_case

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GetToken retrieves an authentication token from Sabre's API
// It updates the SabreClient's Token field with the new access token
// Returns:
//
//	error - Any error encountered during the token retrieval process
func (c *SabreClient) GetToken() error {
	// Log the attempt to get a token (client ID and secret partially for debugging)
	fmt.Println("Getting token", c.ClientID, c.ClientSecret)
	url := c.SABREAUTHURL // Sabre's certification token endpoint

	// Encode client ID and secret separately using base64
	encodedID := base64.StdEncoding.EncodeToString([]byte(c.ClientID))
	encodedSecret := base64.StdEncoding.EncodeToString([]byte(c.ClientSecret))

	// Combine encoded credentials with a colon separator
	creds := encodedID + ":" + encodedSecret
	fmt.Println("Base64 Encoded Credentials (ID:Secret):", creds) // Log for debugging

	// Encode the combined credentials again for the Basic Auth header
	Token := base64.StdEncoding.EncodeToString([]byte(creds))
	fmt.Println("Final Base64 Token for Header:", Token) // Log for debugging

	// Prepare the payload for client credentials grant type
	payload := strings.NewReader("grant_type=client_credentials")

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("failed to create token request: %v", err)
	}

	// Set necessary headers for authentication
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Required for form data
	req.Header.Set("Authorization", "Basic "+Token)                     // Basic Auth with encoded credentials

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %v", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed after use

	// Check if the response status is successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token request returned status %d: %s", resp.StatusCode, string(body))
	}

	// Define a struct to parse the token response
	var tokenResp struct {
		AccessToken string `json:"access_token"` // Field to capture the access token
	}

	// Decode the JSON response into the struct
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %v", err)
	}

	// Store the token in the client instance
	c.Token = tokenResp.AccessToken
	fmt.Println("Successfully retrieved Token:", c.Token) // Log the retrieved token

	return nil // Success
}
