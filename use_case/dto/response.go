package DTO

type BaggageInfo struct {
	PassengerType string `json:"passenger_type"`
	Allowance     string `json:"allowance"`
}

type Flight struct {
	FlightNumber  string        `json:"flight_number"`
	Origin        string        `json:"origin"`
	Destination   string        `json:"destination"`
	DepartureTime string        `json:"departure_time"`
	ArrivalTime   string        `json:"arrival_time"`
	Price         float64       `json:"price"`
	Baggage       []BaggageInfo `json:"baggage"`
}

type FlightSearchResponse struct {
	Flights []Flight `json:"flights"`
}

type SabreResponse struct {
	GroupedItineraryResponse struct {
		Version  string `json:"version"`
		Messages []struct {
			Severity string `json:"severity"`
			Type     string `json:"type"`
			Code     string `json:"code"`
			Text     string `json:"text"`
		} `json:"messages"`
		Statistics struct {
			ItineraryCount int `json:"itineraryCount"`
		} `json:"statistics"`
		ScheduleDescs []struct {
			ID              int    `json:"id"`
			Frequency       string `json:"frequency"`
			StopCount       int    `json:"stopCount"`
			ETicketable     bool   `json:"eTicketable"`
			TotalMilesFlown int    `json:"totalMilesFlown"`
			ElapsedTime     int    `json:"elapsedTime"`
			Departure       struct {
				Airport string `json:"airport"`
				City    string `json:"city"`
				Country string `json:"country"`
				Time    string `json:"time"`
			} `json:"departure"`
			Arrival struct {
				Airport string `json:"airport"`
				City    string `json:"city"`
				Country string `json:"country"`
				Time    string `json:"time"`
			} `json:"arrival"`
			Carrier struct {
				Marketing             string `json:"marketing"`
				MarketingFlightNumber int    `json:"marketingFlightNumber"`
				Operating             string `json:"operating"`
				OperatingFlightNumber int    `json:"operatingFlightNumber"`
				Equipment             struct {
					Code            string `json:"code"`
					TypeForFirstLeg string `json:"typeForFirstLeg"`
					TypeForLastLeg  string `json:"typeForLastLeg"`
				} `json:"equipment"`
			} `json:"carrier"`
		} `json:"scheduleDescs"`
		BaggageAllowanceDescs []struct {
			ID         int `json:"id"`
			PieceCount int `json:"pieceCount"`
		} `json:"baggageAllowanceDescs"`
		LegDescs []struct {
			ID          int `json:"id"`
			ElapsedTime int `json:"elapsedTime"`
			Schedules   []struct {
				Ref int `json:"ref"`
			} `json:"schedules"`
		} `json:"legDescs"`
		ItineraryGroups []struct {
			GroupDescription struct {
				LegDescriptions []struct {
					DepartureDate     string `json:"departureDate"`
					DepartureLocation string `json:"departureLocation"`
					ArrivalLocation   string `json:"arrivalLocation"`
				} `json:"legDescriptions"`
			} `json:"groupDescription"`
			Itineraries []struct {
				ID            int    `json:"id"`
				PricingSource string `json:"pricingSource"`
				Legs          []struct {
					Ref int `json:"ref"`
				} `json:"legs"`
				PricingInformation []struct {
					PricingSubsource string `json:"pricingSubsource"`
					Fare             struct {
						OfferItemID           string `json:"offerItemId"`
						MandatoryInd          bool   `json:"mandatoryInd"`
						ServiceID             string `json:"serviceId"`
						ValidatingCarrierCode string `json:"validatingCarrierCode"`
						Vita                  bool   `json:"vita"`
						ETicketable           bool   `json:"eTicketable"`
						LastTicketDate        string `json:"lastTicketDate"`
						LastTicketTime        string `json:"lastTicketTime"`
						GoverningCarriers     string `json:"governingCarriers"`
						PassengerInfoList     []struct {
							PassengerInfo struct {
								PassengerType   string `json:"passengerType"`
								PassengerNumber int    `json:"passengerNumber"`
								NonRefundable   bool   `json:"nonRefundable"`
								FareComponents  []struct {
									Ref          int    `json:"ref"`
									BeginAirport string `json:"beginAirport"`
									EndAirport   string `json:"endAirport"`
									Segments     []struct {
										Segment struct {
											BookingCode       string `json:"bookingCode"`
											CabinCode         string `json:"cabinCode"`
											MealCode          string `json:"mealCode"`
											SeatsAvailable    int    `json:"seatsAvailable"`
											AvailabilityBreak bool   `json:"availabilityBreak"`
										} `json:"segment"`
									} `json:"segments"`
								} `json:"fareComponents"`
								Taxes []struct {
									Ref int `json:"ref"`
								} `json:"taxes"`
								TaxSummaries []struct {
									Ref int `json:"ref"`
								} `json:"taxSummaries"`
								PassengerTotalFare struct {
									TotalFare            float64 `json:"totalFare"`
									TotalTaxAmount       float64 `json:"totalTaxAmount"`
									Currency             string  `json:"currency"`
									BaseFareAmount       float64 `json:"baseFareAmount"`
									BaseFareCurrency     string  `json:"baseFareCurrency"`
									EquivalentAmount     float64 `json:"equivalentAmount"`
									EquivalentCurrency   string  `json:"equivalentCurrency"`
									ConstructionAmount   float64 `json:"constructionAmount"`
									ConstructionCurrency string  `json:"constructionCurrency"`
									ExchangeRateOne      float64 `json:"exchangeRateOne"`
								} `json:"passengerTotalFare"`
								BaggageInformation []struct {
									ProvisionType string `json:"provisionType"`
									AirlineCode   string `json:"airlineCode"`
									Segments      []struct {
										ID int `json:"id"`
									} `json:"segments"`
									Allowance struct {
										Ref int `json:"ref"`
									} `json:"allowance"`
								} `json:"baggageInformation"`
							} `json:"passengerInfo"`
						} `json:"passengerInfoList"`
						TotalFare struct {
							TotalPrice           float64 `json:"totalPrice"`
							TotalTaxAmount       float64 `json:"totalTaxAmount"`
							Currency             string  `json:"currency"`
							BaseFareAmount       float64 `json:"baseFareAmount"`
							BaseFareCurrency     string  `json:"baseFareCurrency"`
							ConstructionAmount   float64 `json:"constructionAmount"`
							ConstructionCurrency string  `json:"constructionCurrency"`
							EquivalentAmount     float64 `json:"equivalentAmount"`
							EquivalentCurrency   string  `json:"equivalentCurrency"`
						} `json:"totalFare"`
					} `json:"fare"`
				} `json:"pricingInformation"`
			} `json:"itineraries"`
		} `json:"itineraryGroups"`
	} `json:"groupedItineraryResponse"`
}
