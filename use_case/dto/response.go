package DTO

type ResponseBaggageInfo struct {
	PassengerType     string `json:"passenger_type"`
	Allowance         string `json:"allowance"`
	PricePerPassenger string `json:"price_per_passenger"`
	PassengerNumber   int    `json:"passenger_number"`
}

type Flight struct {
	FlightData    []FlightDataScheduleDesc `json:"flight_data"`
	DepartureDate string                   `json:"departure_time"`
}

type FlightDataScheduleDesc struct {
	ScheduleDesc ScheduleDesc `json:"schedule_desc"`
	Baggage      interface{}  `json:"baggage"`
}

type FlightSearchResponse struct {
	Flights []struct {
		Flights    []Flight `json:"flights"`
		TotalPrice string
	}
}

type SabreResponse struct {
	GroupedItineraryResponse GroupedItineraryResponse `json:"groupedItineraryResponse"`
}

// GroupedItineraryResponse represents the main response content
type GroupedItineraryResponse struct {
	Version               string                 `json:"version"`
	Messages              []Message              `json:"messages"`
	Statistics            Statistics             `json:"statistics"`
	ScheduleDescs         []ScheduleDesc         `json:"scheduleDescs"`
	FareComponentDescs    []FareComponentType    `json:"fareComponentDescs"`
	BaggageAllowanceDescs []BaggageAllowanceType `json:"baggageAllowanceDescs"`
	BaggageChargeDescs    []BaggageChargeType    `json:"baggageChargeDescs"`
	LegDescs              []LegDesc              `json:"legDescs"`
	ItineraryGroups       []ItineraryGroup       `json:"itineraryGroups"`
}
type BaggageAllowanceType struct {
	ID           int     `json:"id"` // Required field
	Description1 *string `json:"description1,omitempty"`
	Description2 *string `json:"description2,omitempty"`
	PieceCount   *int    `json:"pieceCount,omitempty"`
	Unit         *string `json:"unit,omitempty"`
	Weight       *int    `json:"weight,omitempty"`
}
type BaggageChargeType struct {
	ID                   int      `json:"id"`
	Description1         *string  `json:"description1,omitempty"`
	Description2         *string  `json:"description2,omitempty"`
	EquivalentAmount     *float64 `json:"equivalentAmount,omitempty"`
	EquivalentCurrency   *string  `json:"equivalentCurrency,omitempty"`
	FirstPiece           *int     `json:"firstPiece,omitempty"`
	LastPiece            *int     `json:"lastPiece,omitempty"`
	NoChargeNotAvailable *string  `json:"noChargeNotAvailable,omitempty"`
}

// Message represents an individual message in the response
type Message struct {
	Severity string `json:"severity"`
	Type     string `json:"type"`
	Code     string `json:"code"`
	Text     string `json:"text"`
}

// Statistics holds statistical data about itineraries
type Statistics struct {
	ItineraryCount int `json:"itineraryCount"`
}

// ScheduleDesc describes a flight schedule
type ScheduleDesc struct {
	ID              int      `json:"id"`
	Frequency       string   `json:"frequency"`
	StopCount       int      `json:"stopCount"`
	ETicketable     bool     `json:"eTicketable"`
	TotalMilesFlown int      `json:"totalMilesFlown"`
	ElapsedTime     int      `json:"elapsedTime"`
	Departure       Location `json:"departure"`
	Arrival         Location `json:"arrival"`
	Carrier         Carrier  `json:"carrier"`
}

// Location represents a departure or arrival location
type Location struct {
	Airport string `json:"airport"`
	City    string `json:"city"`
	Country string `json:"country"`
	Time    string `json:"time"`
}

// Carrier contains flight carrier details
type Carrier struct {
	Marketing             string    `json:"marketing"`
	MarketingFlightNumber int       `json:"marketingFlightNumber"`
	Operating             string    `json:"operating"`
	OperatingFlightNumber int       `json:"operatingFlightNumber"`
	Equipment             Equipment `json:"equipment"`
}

// Equipment describes the aircraft equipment
type Equipment struct {
	Code            string `json:"code"`
	TypeForFirstLeg string `json:"typeForFirstLeg"`
	TypeForLastLeg  string `json:"typeForLastLeg"`
}

// LegDesc describes a leg of a journey
type LegDesc struct {
	ID          int           `json:"id"`
	ElapsedTime int           `json:"elapsedTime"`
	Schedules   []ScheduleRef `json:"schedules"`
}

// ScheduleRef references a schedule in a leg
type ScheduleRef struct {
	Ref int `json:"ref"`
}

// ItineraryGroup groups related itineraries
type ItineraryGroup struct {
	GroupDescription GroupDescription `json:"groupDescription"`
	Itineraries      []Itinerary      `json:"itineraries"`
}

// GroupDescription describes the group of itineraries
type GroupDescription struct {
	LegDescriptions []LegDescription `json:"legDescriptions"`
}

// LegDescription describes a leg within a group
type LegDescription struct {
	DepartureDate     string `json:"departureDate"`
	DepartureLocation string `json:"departureLocation"`
	ArrivalLocation   string `json:"arrivalLocation"`
}

// Itinerary represents a single travel itinerary
type Itinerary struct {
	ID                 int           `json:"id"`
	PricingSource      string        `json:"pricingSource"`
	Legs               []LegRef      `json:"legs"`
	PricingInformation []PricingInfo `json:"pricingInformation"`
}

// LegRef references a leg in an itinerary
type LegRef struct {
	Ref int `json:"ref"`
}

// PricingInfo contains pricing details for an itinerary
type PricingInfo struct {
	PricingSubsource string `json:"pricingSubsource"`
	Fare             Fare   `json:"fare"`
}

// Fare contains fare details
type Fare struct {
	OfferItemID           string          `json:"offerItemId"`
	MandatoryInd          bool            `json:"mandatoryInd"`
	ServiceID             string          `json:"serviceId"`
	ValidatingCarrierCode string          `json:"validatingCarrierCode"`
	Vita                  bool            `json:"vita"`
	ETicketable           bool            `json:"eTicketable"`
	LastTicketDate        string          `json:"lastTicketDate"`
	LastTicketTime        string          `json:"lastTicketTime"`
	GoverningCarriers     string          `json:"governingCarriers"`
	PassengerInfoList     []PassengerInfo `json:"passengerInfoList"`
	TotalFare             TotalFare       `json:"totalFare"`
}

// PassengerInfo wraps passenger-specific details
type PassengerInfo struct {
	PassengerInfo PassengerDetails `json:"passengerInfo"`
}
type PassengerTotalFare struct {
	TotalFare            float64 `json:"totalFare"`
	TotalTaxAmount       float64 `json:"totalTaxAmount"`
	Currency             string  `json:"currency"`
	BaseFareAmount       float64 `json:"baseFareAmount"`
	BaseFareCurrency     string  `json:"baseFareCurrency"`
	EquivalentAmount     float64 `json:"equivalentAmount"`
	EquivalentCurrency   string  `json:"equivalentCurrency"`
	ConstructionAmount   float64 `json:"constructionAmount"`
	ConstructionCurrency string  `json:"constructionCurrency"`
	CommissionPercentage float64 `json:"commissionPercentage"`
	CommissionAmount     float64 `json:"commissionAmount"`
	ExchangeRateOne      float64 `json:"exchangeRateOne"`
}

// PassengerDetails contains detailed passenger fare information
type PassengerDetails struct {
	PassengerType      string                   `json:"passengerType"`
	PassengerNumber    int                      `json:"passengerNumber"`
	NonRefundable      bool                     `json:"nonRefundable"`
	FareComponents     []FareComponent          `json:"fareComponents"`
	Taxes              []TaxRef                 `json:"taxes"`
	TaxSummaries       []TaxSummaryRef          `json:"taxSummaries"`
	PassengerTotalFare PassengerTotalFare       `json:"passengerTotalFare"`
	BaggageInformation []BaggageInformationType `json:"baggageInformation"`
}

// FareComponent describes a component of the fare
type FareComponent struct {
	Ref          int       `json:"ref"`
	BeginAirport string    `json:"beginAirport"`
	EndAirport   string    `json:"endAirport"`
	Segments     []Segment `json:"segments"`
}

// Segment wraps segment-specific details
type Segment struct {
	Segment SegmentDetails `json:"segment"`
}

// SegmentDetails contains details about a flight segment
type SegmentDetails struct {
	BookingCode       string `json:"bookingCode"`
	CabinCode         string `json:"cabinCode"`
	MealCode          string `json:"mealCode"`
	SeatsAvailable    int    `json:"seatsAvailable"`
	AvailabilityBreak bool   `json:"availabilityBreak"`
}

// TaxRef references a tax
type TaxRef struct {
	Ref int `json:"ref"`
}

// TaxSummaryRef references a tax summary
type TaxSummaryRef struct {
	Ref int `json:"ref"`
}

// TotalFare contains total fare information
type TotalFare struct {
	TotalPrice           float64 `json:"totalPrice"`
	TotalTaxAmount       float64 `json:"totalTaxAmount"`
	Currency             string  `json:"currency"`
	BaseFareAmount       float64 `json:"baseFareAmount"`
	BaseFareCurrency     string  `json:"baseFareCurrency"`
	EquivalentAmount     float64 `json:"equivalentAmount"`
	EquivalentCurrency   string  `json:"equivalentCurrency"`
	ConstructionAmount   float64 `json:"constructionAmount"`
	ConstructionCurrency string  `json:"constructionCurrency"`
	ExchangeRateOne      float64 `json:"exchangeRateOne,omitempty"`
}

// BaggageInfo contains baggage information for a passenger
type BaggageInformationType struct {
	AirlineCode   string        `json:"airlineCode"` // Required field
	Allowance     *Allowance    `json:"allowance,omitempty"`
	Charge        *Charge       `json:"charge,omitempty"`
	ProvisionType string        `json:"provisionType"` // Required field
	Segments      []SegmentType `json:"segments"`      // Required field (array with minItems: 1)
}

// Allowance represents a reference to a Baggage Allowance ID
type Allowance struct {
	Ref int `json:"ref"`
}

// Charge represents a reference to a Baggage Charge ID
type Charge struct {
	Ref int `json:"ref"`
}

// Segment represents a segment index in the itinerary
type SegmentType struct {
	ID int `json:"id"` // Assuming 'id' as the property based on context
}

// SegmentRef references a segment in baggage info
type SegmentRef struct {
	ID int `json:"id"`
}

// AllowanceRef references a baggage allowance
type AllowanceRef struct {
	Ref int `json:"ref"`
}

// FareComponentType represents a fare component in fareComponentDescs
type FareComponentType struct {
	ID                          int                    `json:"id"` // Required
	GoverningCarrier            *string                `json:"governingCarrier,omitempty"`
	FareAmount                  *float64               `json:"fareAmount,omitempty"`
	FareCurrency                *string                `json:"fareCurrency,omitempty"`
	FareBasisCode               *string                `json:"fareBasisCode,omitempty"`
	FarePassengerType           *string                `json:"farePassengerType,omitempty"`
	PublishedFareAmount         *float64               `json:"publishedFareAmount,omitempty"`
	PublishedFareCurrency       *string                `json:"publishedFareCurrency,omitempty"`
	OneWayFare                  *bool                  `json:"oneWayFare,omitempty"`
	Directionality              *string                `json:"directionality,omitempty"`
	Direction                   *string                `json:"direction,omitempty"`
	NotValidAfter               *string                `json:"notValidAfter,omitempty"`
	ApplicablePricingCategories *string                `json:"applicablePricingCategories,omitempty"`
	VendorCode                  *string                `json:"vendorCode,omitempty"`
	FareTypeBitmap              *string                `json:"fareTypeBitmap,omitempty"`
	FareType                    *string                `json:"fareType,omitempty"`
	FareTariff                  *string                `json:"fareTariff,omitempty"`
	FareRule                    *string                `json:"fareRule,omitempty"`
	CabinCode                   *string                `json:"cabinCode,omitempty"`
	Segments                    []FareComponentSegment `json:"segments,omitempty"`
}

// FareComponentSegment represents a segment within a fare component
type FareComponentSegment struct {
	Segment *FareSegmentDetails `json:"segment,omitempty"`
}

// FareSegmentDetails contains segment-specific details
type FareSegmentDetails struct {
	Stopover   *bool       `json:"stopover,omitempty"`
	Surcharges []Surcharge `json:"surcharges,omitempty"`
}

// Surcharge represents additional charges on a segment
type Surcharge struct {
	Amount      float64 `json:"amount"`
	Currency    *string `json:"currency,omitempty"`
	Description *string `json:"description,omitempty"`
	Type        *string `json:"type,omitempty"`
}
