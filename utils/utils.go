package utils

import (
	"fmt"
	"sort"
	"strconv"

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
	flights := &DTO.FlightSearchResponse{}

	// Precompute mappings with capacity hints
	baggageMap := make(map[int]string, len(resp.GroupedItineraryResponse.BaggageAllowanceDescs))
	for _, allowance := range resp.GroupedItineraryResponse.BaggageAllowanceDescs {
		baggageMap[allowance.ID] = fmt.Sprintf("%dPC", allowance.PieceCount)
	}

	scheduleMap := make(map[int]DTO.ScheduleDesc, len(resp.GroupedItineraryResponse.ScheduleDescs))
	for _, sched := range resp.GroupedItineraryResponse.ScheduleDescs {
		scheduleMap[sched.ID] = sched
	}

	legMap := make(map[int][]int, len(resp.GroupedItineraryResponse.LegDescs))
	for _, leg := range resp.GroupedItineraryResponse.LegDescs {
		schedRefs := make([]int, 0, len(leg.Schedules))
		for _, sched := range leg.Schedules {
			schedRefs = append(schedRefs, sched.Ref)
		}
		legMap[leg.ID] = schedRefs
	}

	// Process itineraries
	for _, group := range resp.GroupedItineraryResponse.ItineraryGroups {
		for _, itin := range group.Itineraries {
			processItinerary(itin, group, req, baggageMap, scheduleMap, legMap, flights)
		}
	}

	// Sort flights by price
	sort.Slice(flights.Flights, func(i, j int) bool {
		priceI, _ := strconv.ParseFloat(flights.Flights[i].FlightData[0].TotalPrice[:len(flights.Flights[i].FlightData[0].TotalPrice)-4], 64) // Extract numeric part
		priceJ, _ := strconv.ParseFloat(flights.Flights[j].FlightData[0].TotalPrice[:len(flights.Flights[j].FlightData[0].TotalPrice)-4], 64)
		return priceI < priceJ
	})

	return flights, nil
}

// processItinerary processes a single itinerary and appends flights to the response.
func processItinerary(itin DTO.Itinerary, group DTO.ItineraryGroup, req *DTO.FlightSearchRequest,
	baggageMap map[int]string, scheduleMap map[int]DTO.ScheduleDesc, legMap map[int][]int,
	flights *DTO.FlightSearchResponse) {
	for _, pricing := range itin.PricingInformation {
		if len(pricing.Fare.PassengerInfoList) == 0 {
			continue
		}

		totalPrice := pricing.Fare.TotalFare.TotalPrice
		priceStr := fmt.Sprintf("%f %s", totalPrice, pricing.Fare.TotalFare.Currency)

		// Build itinerary-wide baggage info
		baggageInfo := make(map[string]struct {
			PassengerNumber   int
			Allowance         string
			PricePerPassenger string
			PassengerType     string
		})

		for idx, passenger := range pricing.Fare.PassengerInfoList {
			for _, bag := range passenger.PassengerInfo.BaggageInformation {
				fmt.Println("passengerType", passenger.PassengerInfo.PassengerType)
				if allowance, ok := baggageMap[bag.Allowance.Ref]; ok {
					for _, seg := range bag.Segments {
						baggageInfo[strconv.Itoa(seg.ID)+strconv.Itoa(idx)] = struct {
							PassengerNumber   int
							Allowance         string
							PricePerPassenger string
							PassengerType     string
						}{
							PassengerType:     passenger.PassengerInfo.PassengerType,
							PassengerNumber:   passenger.PassengerInfo.PassengerNumber,
							Allowance:         allowance,
							PricePerPassenger: fmt.Sprintf("%f %s", passenger.PassengerInfo.PassengerTotalFare.TotalFare, passenger.PassengerInfo.PassengerTotalFare.Currency),
						}
					}
				}
			}
		}

		globalSegIdx := 0
		for i, legRef := range itin.Legs {
			schedRefs, ok := legMap[legRef.Ref]
			if !ok {
				globalSegIdx += len(schedRefs)
				continue
			}

			schedules := make([]DTO.FlightDataScheduleDesc, 0, len(schedRefs))
			departureDate := group.GroupDescription.LegDescriptions[i].DepartureDate

			for _, schedRef := range schedRefs {
				if flightData, ok := scheduleMap[schedRef]; ok {
					baggage := make([]DTO.ResponseBaggageInfo, 0, len(req.Passengers))

					for idx := range req.Passengers {
						if allowance, exists := baggageInfo[strconv.Itoa(globalSegIdx)+strconv.Itoa(idx)]; exists {
							baggage = append(baggage, DTO.ResponseBaggageInfo{
								Allowance:         allowance.Allowance,
								PricePerPassenger: allowance.PricePerPassenger,
								PassengerNumber:   allowance.PassengerNumber,
								PassengerType:     allowance.PassengerType,
							})
						}
					}

					schedules = append(schedules, DTO.FlightDataScheduleDesc{
						ScheduleDesc: flightData,
						Baggage:      baggage,
						TotalPrice:   priceStr,
					})
				}
				globalSegIdx++
			}

			if len(schedules) > 0 {
				flights.Flights = append(flights.Flights, DTO.Flight{
					DepartureDate: departureDate,
					FlightData:    schedules,
				})
			}
		}
	}
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
			OriginLocation: DTO.RequestLocation{
				LocationCode: req.Origin,
				LocationType: "A", // Airport
			},
			DestinationLocation: DTO.RequestLocation{
				LocationCode: req.Destination,
				LocationType: "A",
			},
			DepartureDateTime: req.DepartureDateTime,
		},
	}

	// Add return leg if round trip
	if req.TripType == "round_trip" {
		originDest = append(originDest, DTO.OriginDest{
			OriginLocation:      DTO.RequestLocation{LocationCode: req.Destination, LocationType: "A"},
			DestinationLocation: DTO.RequestLocation{LocationCode: req.Origin, LocationType: "A"},
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
