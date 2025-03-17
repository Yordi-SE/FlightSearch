package use_case

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Yordi-SE/FlightSearch/models"
	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
)

type SabreClient struct {
	ClientID     string
	ClientSecret string
	Token        string
	PCC          string
}

func NewSabreClient(clientID, clientSecret string, PCC string) *SabreClient {
	return &SabreClient{ClientID: clientID, ClientSecret: clientSecret, PCC: PCC}
}

// GetToken fetches an OAuth token from Sabre's authentication API

// SearchFlights calls the Bargain Finder Max API
func (c *SabreClient) SearchFlights(req *DTO.FlightSearchRequest) ([]models.Flight, error) {
	if c.Token == "" {
		if err := c.GetToken(); err != nil {
			return nil, err
		}
	}

	url := "https://api.cert.platform.sabre.com/v5/offers/shop"
	sabreReq := map[string]interface{}{
		"OTA_AirLowFareSearchRQ": map[string]interface{}{
			"Version": "5",
			"POS": map[string]interface{}{
				"Source": []map[string]interface{}{
					{
						"PseudoCityCode": c.PCC,
						"RequestorID": map[string]interface{}{
							"CompanyName": map[string]string{
								"Code": "TN",
							},
							"ID":   "1",
							"Type": "1",
						},
					},
				},
			},
			"OriginDestinationInformation": []map[string]interface{}{
				{
					"OriginLocation": map[string]interface{}{
						"LocationCode": "ADD",
					},
					"DestinationLocation": map[string]interface{}{
						"LocationCode": "NBO",
					},
					"DepartureDateTime": "2025-03-16T20:00:00",
				},
			},
			"TravelerInfoSummary": map[string]interface{}{
				"AirTravelerAvail": []map[string]interface{}{
					{
						"PassengerTypeQuantity": []map[string]interface{}{
							{
								"Code":     "ADT",
								"Quantity": 1,
							},
						},
					},
				},
			},
			"TPA_Extensions": map[string]interface{}{
				"IntelliSellTransaction": map[string]interface{}{
					"RequestType": map[string]string{
						"Name": "50ITINS",
					},
				},
			},
		},
	}
	c.buildSabreRequest(req)
	payload, err := json.Marshal(sabreReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
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

	// Check if no itineraries were returned
	if sabreResp.GroupedItineraryResponse.Statistics.ItineraryCount == 0 {
		for _, msg := range sabreResp.GroupedItineraryResponse.Messages {
			if msg.Type == "SCHEDULES" && msg.Text == "NO FLIGHT SCHEDULES FOR QUALIFIERS USED" {
				return nil, fmt.Errorf("no flights found matching your search criteria (e.g., dates, route, or preferences)")
			}
		}
		return nil, fmt.Errorf("no flights available for your search; try adjusting your dates or preferences")
	}
	return parseSabreResponse(sabreResp, req)
}

// buildSabreRequest constructs the Sabre API request payload
func (c *SabreClient) buildSabreRequest(req *DTO.FlightSearchRequest) models.OTA_AirLowFareSearchRQ {
	passengers := []models.PassengerTypeQuantity{}
	for _, p := range req.Passengers {
		passengers = append(passengers, models.PassengerTypeQuantity{
			Code:     p.Type,
			Quantity: p.Count,
		})
	}

	originDest := []models.OriginDest{
		{
			OriginLocation: models.Location{
				LocationCode: req.Origin,
				LocationType: "A",
			},
			DestinationLocation: models.Location{
				LocationCode: req.Destination,
				LocationType: "A",
			},
			DepartureDateTime: req.DepartureDateTime,
		},
	}
	if req.TripType == "round_trip" {
		originDest = append(originDest, models.OriginDest{
			OriginLocation:      models.Location{LocationCode: req.Destination, LocationType: "A"},
			DestinationLocation: models.Location{LocationCode: req.Origin, LocationType: "A"},
			DepartureDateTime:   req.ReturnDateTime,
		})
	}

	return models.OTA_AirLowFareSearchRQ{
		Version: "5",
		POS: models.POS{
			Source: []models.Source{
				{
					PseudoCityCode: c.PCC,
					RequestorID: models.RequestorID{
						CompanyName: models.CompanyName{
							Code: "TN",
						},
						ID:   "1",
						Type: "1",
					},
				},
			},
		},
		OriginDestinationInformation: originDest,
		TravelerInfoSummary: models.TravelerInfoSummary{
			AirTravelerAvail: []models.AirTravelerAvail{
				{
					PassengerTypeQuantity: passengers,
				},
			},
		},
		TPA_Extensions: models.TPAExtensions{
			IntelliSellTransaction: models.IntelliSellTransaction{
				RequestType: models.RequestType{
					Name: "50ITINS",
				},
			},
		},
	}
}

// parseSabreResponse converts Sabre's response to our Flight model
func parseSabreResponse(resp DTO.SabreResponse, req *DTO.FlightSearchRequest) ([]models.Flight, error) {
	var flights []models.Flight

	// Map baggage allowances by ID
	baggageMap := make(map[int]string)
	for _, allowance := range resp.GroupedItineraryResponse.BaggageAllowanceDescs {
		baggageMap[allowance.ID] = fmt.Sprintf("%dPC", allowance.PieceCount)
	}

	// Map schedules by ID
	scheduleMap := make(map[int]struct {
		FlightNumber  string
		Origin        string
		Destination   string
		DepartureTime string
		ArrivalTime   string
	})
	for _, sched := range resp.GroupedItineraryResponse.ScheduleDescs {
		scheduleMap[sched.ID] = struct {
			FlightNumber  string
			Origin        string
			Destination   string
			DepartureTime string
			ArrivalTime   string
		}{
			FlightNumber:  fmt.Sprintf("%s%d", sched.Carrier.Marketing, sched.Carrier.MarketingFlightNumber),
			Origin:        sched.Departure.Airport,
			Destination:   sched.Arrival.Airport,
			DepartureTime: sched.Departure.Time,
			ArrivalTime:   sched.Arrival.Time,
		}
	}

	// Map legs by ID
	legMap := make(map[int][]int) // Maps leg ID to schedule refs
	for _, leg := range resp.GroupedItineraryResponse.LegDescs {
		var schedRefs []int
		for _, sched := range leg.Schedules {
			schedRefs = append(schedRefs, sched.Ref)
		}
		legMap[leg.ID] = schedRefs
	}

	// Process itineraries
	for _, group := range resp.GroupedItineraryResponse.ItineraryGroups {
		for _, itin := range group.Itineraries {
			for _, pricing := range itin.PricingInformation {
				if len(pricing.Fare.PassengerInfoList) == 0 {
					continue
				}
				passengerInfo := pricing.Fare.PassengerInfoList[0].PassengerInfo
				totalPrice := pricing.Fare.TotalFare.TotalPrice
				numLegs := len(itin.Legs)

				// Build baggage info per segment
				baggageInfo := make(map[int]string)
				for _, bag := range passengerInfo.BaggageInformation {
					allowance, ok := baggageMap[bag.Allowance.Ref]
					if ok {
						for _, seg := range bag.Segments {
							baggageInfo[seg.ID] = allowance
						}
					}
				}

				// Process each leg
				for i, legRef := range itin.Legs {
					legID := legRef.Ref
					schedRefs, ok := legMap[legID]
					if !ok {
						continue
					}

					for idx, schedRef := range schedRefs {
						if flightData, ok := scheduleMap[schedRef]; ok {
							baggage := []models.BaggageInfo{}
							for _, p := range req.Passengers {
								if allowance, exists := baggageInfo[idx]; exists {
									baggage = append(baggage, models.BaggageInfo{
										PassengerType: p.Type,
										Allowance:     allowance,
									})
								}
							}

							departureDate := group.GroupDescription.LegDescriptions[i].DepartureDate
							departureTime := fmt.Sprintf("%sT%s", departureDate, flightData.DepartureTime[:8])

							flights = append(flights, models.Flight{
								FlightNumber:  flightData.FlightNumber,
								Origin:        flightData.Origin,
								Destination:   flightData.Destination,
								DepartureTime: departureTime,
								ArrivalTime:   fmt.Sprintf("%sT%s", departureDate, flightData.ArrivalTime[:8]),
								Price:         totalPrice / float64(numLegs),
								Baggage:       baggage,
							})
						}
					}
				}
			}
		}
	}

	fmt.Println("Parsed Flights:", len(flights))
	return flights, nil
}
