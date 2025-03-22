package utils

import (
	"fmt"

	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
)

// ParseSabreResponse converts Sabre's API response into our internal Flight model
// Args:
//
//	resp - The raw response from Sabre API
//	req - The original flight search request
//
// Returns:
//
//	Pointer to FlightSearchResponse containing parsed flights and any error encountered
func ParseSabreResponse(resp DTO.SabreResponse, req *DTO.FlightSearchRequest) (*DTO.FlightSearchResponse, error) {
	var flights DTO.FlightSearchResponse

	// Create a mapping of baggage allowance descriptions using ID as key
	// Format: "2PC" indicates 2 pieces allowed
	baggageMap := make(map[int]string)
	for _, allowance := range resp.GroupedItineraryResponse.BaggageAllowanceDescs {
		baggageMap[allowance.ID] = fmt.Sprintf("%dPC", allowance.PieceCount)
	}

	// Create a mapping of flight schedules using schedule ID as key
	// Stores essential flight details for quick lookup
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

	// Create a mapping of leg descriptions to their schedule references
	// Maps leg ID to array of schedule IDs
	legMap := make(map[int][]int)
	for _, leg := range resp.GroupedItineraryResponse.LegDescs {
		var schedRefs []int
		for _, sched := range leg.Schedules {
			schedRefs = append(schedRefs, sched.Ref)
		}
		legMap[leg.ID] = schedRefs
	}

	// Process each itinerary group and its pricing information
	for _, group := range resp.GroupedItineraryResponse.ItineraryGroups {
		for _, itin := range group.Itineraries {
			for _, pricing := range itin.PricingInformation {
				// Skip if no passenger info available
				if len(pricing.Fare.PassengerInfoList) == 0 {
					continue
				}
				passengerInfo := pricing.Fare.PassengerInfoList[0].PassengerInfo
				totalPrice := pricing.Fare.TotalFare.TotalPrice
				numLegs := len(itin.Legs)

				// Build baggage information mapping per segment
				baggageInfo := make(map[int]string)
				for _, bag := range passengerInfo.BaggageInformation {
					if allowance, ok := baggageMap[bag.Allowance.Ref]; ok {
						for _, seg := range bag.Segments {
							baggageInfo[seg.ID] = allowance
						}
					}
				}

				// Process each leg of the itinerary
				for i, legRef := range itin.Legs {
					legID := legRef.Ref
					schedRefs, ok := legMap[legID]
					if !ok {
						continue
					}

					// Process each schedule in the leg
					for idx, schedRef := range schedRefs {
						if flightData, ok := scheduleMap[schedRef]; ok {
							// Build baggage info for each passenger type
							baggage := []DTO.BaggageInfo{}
							for _, p := range req.Passengers {
								if allowance, exists := baggageInfo[idx]; exists {
									baggage = append(baggage, DTO.BaggageInfo{
										PassengerType: p.Type,
										Allowance:     allowance,
									})
								}
							}

							// Format departure time with date
							departureDate := group.GroupDescription.LegDescriptions[i].DepartureDate
							departureTime := fmt.Sprintf("%sT%s", departureDate, flightData.DepartureTime[:8])

							// Add flight to response
							flights.Flights = append(flights.Flights, DTO.Flight{
								FlightNumber:  flightData.FlightNumber,
								Origin:        flightData.Origin,
								Destination:   flightData.Destination,
								DepartureTime: departureTime,
								ArrivalTime:   fmt.Sprintf("%sT%s", departureDate, flightData.ArrivalTime[:8]),
								Price:         totalPrice / float64(numLegs), // Split total price across legs
								Baggage:       baggage,
							})
						}
					}
				}
			}
		}
	}

	// Log the number of flights parsed
	fmt.Println("Parsed Flights:", len(flights.Flights))
	return &flights, nil
}

// BuildSabreRequest constructs the request payload for Sabre API
// Args:
//
//	req - The flight search request from the client
//	PCC - Pseudo City Code for authentication
//
// Returns:
//
//	Formatted Sabre request structure
func BuildSabreRequest(req *DTO.FlightSearchRequest, PCC string) DTO.SabreRequestFormat {
	// Convert passenger info to Sabre format
	passengers := []DTO.PassengerTypeQuantity{}
	for _, p := range req.Passengers {
		passengers = append(passengers, DTO.PassengerTypeQuantity{
			Code:     p.Type,
			Quantity: p.Count,
		})
	}

	// Build origin-destination information for one-way trip
	originDest := []DTO.OriginDest{
		{
			OriginLocation: DTO.Location{
				LocationCode: req.Origin,
				LocationType: "A", // Airport
			},
			DestinationLocation: DTO.Location{
				LocationCode: req.Destination,
				LocationType: "A",
			},
			DepartureDateTime: req.DepartureDateTime,
		},
	}

	// Add return leg if round trip
	if req.TripType == "round_trip" {
		originDest = append(originDest, DTO.OriginDest{
			OriginLocation:      DTO.Location{LocationCode: req.Destination, LocationType: "A"},
			DestinationLocation: DTO.Location{LocationCode: req.Origin, LocationType: "A"},
			DepartureDateTime:   req.ReturnDateTime,
		})
	}

	// Construct and return the complete Sabre request
	return DTO.SabreRequestFormat{
		OTA_AirLowFareSearchRQ: DTO.OTA_AirLowFareSearchRQ{
			Version: "5",
			POS: DTO.POS{
				Source: []DTO.Source{
					{
						PseudoCityCode: PCC,
						RequestorID: DTO.RequestorID{
							CompanyName: DTO.CompanyName{
								Code: "TN",
							},
							ID:   "1",
							Type: "1",
						},
					},
				},
			},
			OriginDestinationInformation: originDest,
			TravelerInfoSummary: DTO.TravelerInfoSummary{
				AirTravelerAvail: []DTO.AirTravelerAvail{
					{
						PassengerTypeQuantity: passengers,
					},
				},
			},
			TPA_Extensions: DTO.TPAExtensions{
				IntelliSellTransaction: DTO.IntelliSellTransaction{
					RequestType: DTO.RequestType{
						Name: "50ITINS", // Request up to 50 itineraries
					},
				},
			},
		},
	}
}
