package use_case

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Yordi-SE/FlightSearch/config"
	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
	"github.com/Yordi-SE/FlightSearch/utils"
)

// SabreClient represents a client for interacting with the Sabre API
type SabreClient struct {
	ClientID     string // API client ID for authentication
	ClientSecret string // API client secret for authentication
	Token        string // Current access token (refreshed as needed)
	URL          string // Sabre API endpoint URL
	PCC          string // Pseudo City Code for agency identification
	SABREAUTHURL string // Sabre authentication endpoint URL
}

// NewSabreClient creates and initializes a new SabreClient instance
// Args:
//
//	clientID - The API client ID
//	clientSecret - The API client secret
//	PCC - Pseudo City Code
//	url - The Sabre API endpoint
//
// Returns:
//
//	Pointer to a new SabreClient instance
func NewSabreClient(Config *config.Config) *SabreClient {
	return &SabreClient{
		ClientID:     Config.ClientID,
		ClientSecret: Config.ClientSecret,
		PCC:          Config.PCC,
		URL:          Config.URL,
		SABREAUTHURL: Config.SABREAUTHURL,
	}
}

// SearchFlights executes a flight search using Sabre's Bargain Finder Max API
// Args:
//
//	req - The flight search request containing search parameters
//
// Returns:
//
//	Pointer to FlightSearchResponse with search results or an error if the request fails
func (c *SabreClient) SearchFlights(req *DTO.FlightSearchRequest) (*DTO.FlightSearchResponse, error) {
	// Ensure we have a valid token; fetch one if not present
	if c.Token == "" {
		if err := c.GetToken(); err != nil {
			return nil, fmt.Errorf("failed to obtain authentication token: %v", err)
		}
	}

	// Build the Sabre-specific request format from our internal request
	sabreReq := utils.BuildSabreRequest(req, c.PCC)

	// Marshal the request into JSON
	payload, err := json.Marshal(sabreReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create flight request: %v", err)
	}

	// Set required headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.Token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("flight request failed: %v", err)
	}
	defer resp.Body.Close() // Ensure body is closed after we're done

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Log response details for debugging

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flight request returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the Sabre response into our structure
	var sabreResp DTO.SabreResponse
	if err := json.Unmarshal(body, &sabreResp); err != nil {
		// Log raw response for debugging and return a detailed error
		fmt.Println("Failed to unmarshal response:", err)
		return nil, fmt.Errorf("invalid response format from Sabre API: %v (raw response: %s)", err, string(body))
	}

	// Handle specific error messages from Sabre
	for _, msg := range sabreResp.GroupedItineraryResponse.Messages {
		if msg.Severity == "Error" {
			switch msg.Text {
			case "No complete journey can be built in IF2/ADVJR1.":
				return nil, fmt.Errorf("no flights available for the specified route and dates")
			case "Error during Processing":
				return nil, fmt.Errorf("an error occurred while searching for flights; please try again later")
			default:
				return nil, fmt.Errorf("sabre processing error: %s (%s)", msg.Text, msg.Code)
			}
		}
	}

	// Check if any itineraries were found
	if sabreResp.GroupedItineraryResponse.Statistics.ItineraryCount == 0 {
		for _, msg := range sabreResp.GroupedItineraryResponse.Messages {
			if msg.Type == "SCHEDULES" && msg.Text == "NO FLIGHT SCHEDULES FOR QUALIFIERS USED" {
				return nil, fmt.Errorf("no flights found matching your search criteria (e.g., dates, route, or preferences)")
			}
		}
		return nil, fmt.Errorf("no flights available for your search; try adjusting your dates or preferences")
	}

	// Parse the response into our flight model and return
	return utils.ParseSabreResponse(sabreResp, req)
}
