package use_case

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *SabreClient) GetToken() error {
	fmt.Println("Getting token", c.ClientID, c.ClientSecret)
	url := "https://api.cert.sabre.com/v2/auth/token"

	encodedID := base64.StdEncoding.EncodeToString([]byte(c.ClientID))
	encodedSecret := base64.StdEncoding.EncodeToString([]byte(c.ClientSecret))

	// Concatenate with a colon
	fmt.Println(encodedID, encodedSecret)
	creds := encodedID + ":" + encodedSecret
	fmt.Println(creds)
	Token := base64.StdEncoding.EncodeToString([]byte(creds))
	// Payload remains the same
	payload := strings.NewReader("grant_type=client_credentials")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("failed to create token request: %v", err)
	}
	fmt.Println(Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token request returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %v", err)
	}

	c.Token = tokenResp.AccessToken
	fmt.Println("Token:", c.Token)
	return nil
}
