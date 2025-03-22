package DTO

type ResponseBaggageInfo struct {
	PassengerType string `json:"passenger_type"`
	Allowance     string `json:"allowance"`
}

type Flight struct {
	FlightData    ScheduleDesc          `json:"schedule_desc"`
	DepartureTime string                `json:"departure_time"`
	ArrivalTime   string                `json:"arrival_time"`
	Price         float64               `json:"price"`
	Baggage       []ResponseBaggageInfo `json:"baggage"`
}

type FlightSearchResponse struct {
	Flights []Flight `json:"flights"`
}

type SabreResponse struct {
	GroupedItineraryResponse GroupedItineraryResponse `json:"groupedItineraryResponse"`
}

// GroupedItineraryResponse represents the main response content
type GroupedItineraryResponse struct {
	Version               string             `json:"version"`
	Messages              []Message          `json:"messages"`
	Statistics            Statistics         `json:"statistics"`
	ScheduleDescs         []ScheduleDesc     `json:"scheduleDescs"`
	BaggageAllowanceDescs []BaggageAllowance `json:"baggageAllowanceDescs"`
	LegDescs              []LegDesc          `json:"legDescs"`
	ItineraryGroups       []ItineraryGroup   `json:"itineraryGroups"`
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

// BaggageAllowance describes baggage allowance details
type BaggageAllowance struct {
	ID         int `json:"id"`
	PieceCount int `json:"pieceCount"`
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

// PassengerDetails contains detailed passenger fare information
type PassengerDetails struct {
	PassengerType      string          `json:"passengerType"`
	PassengerNumber    int             `json:"passengerNumber"`
	NonRefundable      bool            `json:"nonRefundable"`
	FareComponents     []FareComponent `json:"fareComponents"`
	Taxes              []TaxRef        `json:"taxes"`
	TaxSummaries       []TaxSummaryRef `json:"taxSummaries"`
	PassengerTotalFare TotalFare       `json:"passengerTotalFare"`
	BaggageInformation []BaggageInfo   `json:"baggageInformation"`
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
type BaggageInfo struct {
	ProvisionType string       `json:"provisionType"`
	AirlineCode   string       `json:"airlineCode"`
	Segments      []SegmentRef `json:"segments"`
	Allowance     AllowanceRef `json:"allowance"`
}

// SegmentRef references a segment in baggage info
type SegmentRef struct {
	ID int `json:"id"`
}

// AllowanceRef references a baggage allowance
type AllowanceRef struct {
	Ref int `json:"ref"`
}
