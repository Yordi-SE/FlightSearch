package use_case

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
	"github.com/Yordi-SE/FlightSearch/utils"
)

type SabreClient struct {
	ClientID     string
	ClientSecret string
	Token        string
	URL          string
	PCC          string
}

func NewSabreClient(clientID, clientSecret string, PCC string, url string) *SabreClient {
	return &SabreClient{ClientID: clientID, ClientSecret: clientSecret, PCC: PCC, URL: url}
}

// SearchFlights calls the Bargain Finder Max API
func (c *SabreClient) SearchFlights(req *DTO.FlightSearchRequest) (*DTO.FlightSearchResponse, error) {
	if c.Token == "" {
		if err := c.GetToken(); err != nil {
			return nil, err
		}
	}
	
	sabreReq := utils.BuildSabreRequest(req, c.PCC)
	fmt.Printf("Request: %+v\n", sabreReq)
	payload, err := json.Marshal(sabreReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create flight request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("flight request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response Body:", string(body))
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("flight request returned status %d: %s", resp.StatusCode, string(body))
	}

	var sabreResp DTO.SabreResponse
	if err := json.Unmarshal(body, &sabreResp); err != nil {
		// Log the raw response for debugging and return a generic error
		fmt.Println("Failed to unmarshal response:", err)
		return nil, fmt.Errorf("invalid response format from Sabre API: %v (raw response: %s)", err, string(body))
	}

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

	if sabreResp.GroupedItineraryResponse.Statistics.ItineraryCount == 0 {
		for _, msg := range sabreResp.GroupedItineraryResponse.Messages {
			if msg.Type == "SCHEDULES" && msg.Text == "NO FLIGHT SCHEDULES FOR QUALIFIERS USED" {
				return nil, fmt.Errorf("no flights found matching your search criteria (e.g., dates, route, or preferences)")
			}
		}
		return nil, fmt.Errorf("no flights available for your search; try adjusting your dates or preferences")
	}
	return utils.ParseSabreResponse(sabreResp, req)
}
