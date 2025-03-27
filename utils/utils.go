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
	baggageMap := make(map[int]DTO.BaggageAllowanceType, len(resp.GroupedItineraryResponse.BaggageAllowanceDescs))
	for _, allowance := range resp.GroupedItineraryResponse.BaggageAllowanceDescs {
		baggageMap[allowance.ID] = allowance
	}

	fareComponentsMap := make(map[int]DTO.FareComponentType, len(resp.GroupedItineraryResponse.FareComponentDescs))
	for _, fare := range resp.GroupedItineraryResponse.FareComponentDescs {
		fareComponentsMap[fare.ID] = fare
	}

	baggageChargeMap := make(map[int]DTO.BaggageChargeType, len(resp.GroupedItineraryResponse.BaggageChargeDescs))
	for _, charge := range resp.GroupedItineraryResponse.BaggageChargeDescs {
		baggageChargeMap[charge.ID] = charge
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
			processItinerary(itin, group, baggageMap, baggageChargeMap, scheduleMap, legMap, flights, fareComponentsMap)
		}
	}

	// Sort flights by price
	sort.Slice(flights.Flights, func(i, j int) bool {
		priceI, _ := strconv.ParseFloat(flights.Flights[i].TotalPrice[:len(flights.Flights[i].TotalPrice)-4], 64) // Extract numeric part
		priceJ, _ := strconv.ParseFloat(flights.Flights[j].TotalPrice[:len(flights.Flights[j].TotalPrice)-4], 64)
		return priceI < priceJ
	})

	return flights, nil
}

// processItinerary processes a single itinerary and appends flights to the response.
func processItinerary(itin DTO.Itinerary, group DTO.ItineraryGroup,
	baggageMap map[int]DTO.BaggageAllowanceType, baggageChargeMap map[int]DTO.BaggageChargeType, scheduleMap map[int]DTO.ScheduleDesc, legMap map[int][]int,
	flights *DTO.FlightSearchResponse, fare map[int]DTO.FareComponentType) {
	for _, pricing := range itin.PricingInformation {
		if len(pricing.Fare.PassengerInfoList) == 0 {
			continue
		}

		totalPrice := pricing.Fare.TotalFare.TotalPrice
		priceStr := fmt.Sprintf("%f %s", totalPrice, pricing.Fare.TotalFare.Currency)
		// Build itinerary-wide baggage info
		fmt.Println(priceStr)
		baggageInfo := make(map[string][]DTO.BaggageAllowanceType)
		chargeInfo := make(map[string][]DTO.BaggageChargeType)
		for idx, passenger := range pricing.Fare.PassengerInfoList {
			for _, bag := range passenger.PassengerInfo.BaggageInformation {
				if bag.ProvisionType == "C" {
					if charge, ok := baggageChargeMap[bag.Charge.Ref]; ok {
						for _, seg := range bag.Segments {
							chargeInfo[strconv.Itoa(seg.ID)+strconv.Itoa(idx)] = append(chargeInfo[strconv.Itoa(seg.ID)+strconv.Itoa(idx)], charge)
						}
					}
				} else if allowance, ok := baggageMap[bag.Allowance.Ref]; ok {
					for _, seg := range bag.Segments {
						baggageInfo[strconv.Itoa(seg.ID)+strconv.Itoa(idx)] = append(baggageInfo[strconv.Itoa(seg.ID)+strconv.Itoa(idx)], allowance)
					}
				}
			}
		}
		fareInfo := make(map[string]DTO.FareComponentType)
		for idx, passenger := range pricing.Fare.PassengerInfoList {
			for id, passengerFare := range passenger.PassengerInfo.FareComponents {
				if fareComp, ok := fare[passengerFare.Ref]; ok {
					fareInfo[strconv.Itoa(id)+strconv.Itoa(idx)] = fareComp
				}
			}
		}

		globalSegIdx := 0
		ItinFlights := make([]DTO.Flight, 0, len(itin.Legs))

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
					var baggage = make([]struct {
						Baggage         []DTO.BaggageAllowanceType
						Charge          []DTO.BaggageChargeType
						PassengerNumber int
						PassengerType   string
						NonRefundable   bool
						FareComponent   struct {
							FareComponent DTO.FareComponentType
							BeginAirport  string
							EndAirport    string
						}
					}, len(pricing.Fare.PassengerInfoList))

					for idx, passenger := range pricing.Fare.PassengerInfoList {
						if allowance, exists := baggageInfo[strconv.Itoa(globalSegIdx)+strconv.Itoa(idx)]; exists {
							baggage[idx].Baggage = append(baggage[idx].Baggage, allowance...)
						}
						if charge, exists := chargeInfo[strconv.Itoa(globalSegIdx)+strconv.Itoa(idx)]; exists {
							baggage[idx].Charge = append(baggage[idx].Charge, charge...)
						}
						baggage[idx].PassengerNumber = passenger.PassengerInfo.PassengerNumber
						baggage[idx].PassengerType = passenger.PassengerInfo.PassengerType
						baggage[idx].NonRefundable = passenger.PassengerInfo.NonRefundable
						available := false
						for _, farecomp := range passenger.PassengerInfo.FareComponents {
							if farecomp.BeginAirport == flightData.Departure.Airport && farecomp.EndAirport == flightData.Arrival.Airport {
								baggage[idx].FareComponent.FareComponent = fare[farecomp.Ref]
								baggage[idx].FareComponent.BeginAirport = farecomp.BeginAirport
								baggage[idx].FareComponent.EndAirport = farecomp.EndAirport
								available = true
								break
							}
						}
						if !available {
							for _, farecomp := range passenger.PassengerInfo.FareComponents {
								if farecomp.BeginAirport == group.GroupDescription.LegDescriptions[i].DepartureLocation && farecomp.EndAirport == group.GroupDescription.LegDescriptions[i].ArrivalLocation && len(farecomp.Segments) > 1 {
									baggage[idx].FareComponent.FareComponent = fare[farecomp.Ref]
									baggage[idx].FareComponent.BeginAirport = farecomp.BeginAirport
									baggage[idx].FareComponent.EndAirport = farecomp.EndAirport
									break
								}
							}
						}
					}
					schedules = append(schedules, DTO.FlightDataScheduleDesc{
						ScheduleDesc: flightData,
						Baggage:      baggage,
					})
				}
				globalSegIdx++
			}
			if len(schedules) > 0 {
				ItinFlights = append(ItinFlights, DTO.Flight{
					DepartureDate: departureDate,
					FlightData:    schedules,
				})
			}

		}
		if len(ItinFlights) > 0 {
			flights.Flights = append(flights.Flights, struct {
				Flights    []DTO.Flight "json:\"flights\""
				TotalPrice string
			}{
				Flights:    ItinFlights,
				TotalPrice: priceStr,
			})
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
			TravelPreferences: DTO.TravelPreferences{
				Baggage: DTO.Baggage{
					CarryOnInfo: true,
					Description: true,
					RequestType: "C",
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
