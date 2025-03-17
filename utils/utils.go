package utils

import (
	"fmt"

	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
)

// parseSabreResponse converts Sabre's response to our Flight model
func ParseSabreResponse(resp DTO.SabreResponse, req *DTO.FlightSearchRequest) (*DTO.FlightSearchResponse, error) {
	var flights DTO.FlightSearchResponse

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
							baggage := []DTO.BaggageInfo{}
							for _, p := range req.Passengers {
								if allowance, exists := baggageInfo[idx]; exists {
									baggage = append(baggage, DTO.BaggageInfo{
										PassengerType: p.Type,
										Allowance:     allowance,
									})
								}
							}

							departureDate := group.GroupDescription.LegDescriptions[i].DepartureDate
							departureTime := fmt.Sprintf("%sT%s", departureDate, flightData.DepartureTime[:8])

							flights.Flights = append(flights.Flights, DTO.Flight{
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

	fmt.Println("Parsed Flights:", len(flights.Flights))
	return &flights, nil
}

// buildSabreRequest constructs the Sabre API request payload
func BuildSabreRequest(req *DTO.FlightSearchRequest, PCC string) DTO.SabreRequestFormat {
	passengers := []DTO.PassengerTypeQuantity{}
	for _, p := range req.Passengers {
		passengers = append(passengers, DTO.PassengerTypeQuantity{
			Code:     p.Type,
			Quantity: p.Count,
		})
	}

	originDest := []DTO.OriginDest{
		{
			OriginLocation: DTO.Location{
				LocationCode: req.Origin,
				LocationType: "A",
			},
			DestinationLocation: DTO.Location{
				LocationCode: req.Destination,
				LocationType: "A",
			},
			DepartureDateTime: req.DepartureDateTime,
		},
	}
	if req.TripType == "round_trip" {
		originDest = append(originDest, DTO.OriginDest{
			OriginLocation:      DTO.Location{LocationCode: req.Destination, LocationType: "A"},
			DestinationLocation: DTO.Location{LocationCode: req.Origin, LocationType: "A"},
			DepartureDateTime:   req.ReturnDateTime,
		})
	}

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
			TravelPreferences: DTO.TravelPreferences{
				MaxStopsQuantity: 0,
				VendorPref: []DTO.VendorPref{
					{
						Code: "LO",
					},
				},
			},
			TPA_Extensions: DTO.TPAExtensions{
				IntelliSellTransaction: DTO.IntelliSellTransaction{
					RequestType: DTO.RequestType{
						Name: "50ITINS",
					},
				},
			},
		},
	}
}
