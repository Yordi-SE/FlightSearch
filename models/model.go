package models



// FlightSearchResponse represents the output of a flight search
type FlightSearchResponse struct {
	Flights []Flight `json:"flights"` // List of flight options
}

// Flight represents a single flight option
type Flight struct {
	FlightNumber  string        `json:"flight_number"`  // e.g., "AA123"
	Origin        string        `json:"origin"`         // IATA airport code (e.g., "JFK")
	Destination   string        `json:"destination"`    // IATA airport code (e.g., "LAX")
	DepartureTime string        `json:"departure_time"` // Format: "YYYY-MM-DDTHH:MM:SS"
	ArrivalTime   string        `json:"arrival_time"`   // Format: "YYYY-MM-DDTHH:MM:SS"
	Price         float64       `json:"price"`          // Price per flight in USD (or currency from API)
	Baggage       []BaggageInfo `json:"baggage"`        // Baggage allowance per passenger type
}

// BaggageInfo represents baggage allowance for a passenger type
type BaggageInfo struct {
	PassengerType string `json:"passenger_type"` // e.g., "ADT", "CNN"
	Allowance     string `json:"allowance"`      // e.g., "1PC 23kg"
}

type Location struct {
	LocationCode string `json:"LocationCode"`
	LocationType string `json:"LocationType"`
}

type OriginDest struct {
	OriginLocation      Location `json:"OriginLocation"`
	DestinationLocation Location `json:"DestinationLocation"`
	DepartureDateTime   string   `json:"DepartureDateTime"`
}

type OTA_AirLowFareSearchRQ struct {
	Version                      string              `json:"Version"`
	POS                          POS                 `json:"POS"`
	OriginDestinationInformation []OriginDest        `json:"OriginDestinationInformation"`
	TravelPreferences            TravelPreferences   `json:"TravelPreferences"`
	TravelerInfoSummary          TravelerInfoSummary `json:"TravelerInfoSummary"`
	TPA_Extensions               TPAExtensions       `json:"TPA_Extensions"`
}
type POS struct {
	Source []Source `json:"Source"`
}

type Source struct {
	PseudoCityCode string      `json:"PseudoCityCode"`
	RequestorID    RequestorID `json:"RequestorID"`
}

type RequestorID struct {
	Type        string      `json:"Type"`
	ID          string      `json:"ID"`
	CompanyName CompanyName `json:"CompanyName"`
}

type CompanyName struct {
	Code string `json:"Code"`
}

// OriginDestinationInformation related structs
type OriginDestinationInfo struct {
	DepartureDateTime   string   `json:"DepartureDateTime"`
	OriginLocation      Location `json:"OriginLocation"`
	DestinationLocation Location `json:"DestinationLocation"`
}

// TravelPreferences related structs
type TravelPreferences struct {
	MaxStopsQuantity int          `json:"MaxStopsQuantity"`
	VendorPref       []VendorPref `json:"VendorPref"`
}

type VendorPref struct {
	Code string `json:"Code"`
}

// TravelerInfoSummary related structs
type TravelerInfoSummary struct {
	AirTravelerAvail []AirTravelerAvail `json:"AirTravelerAvail"`
}

type AirTravelerAvail struct {
	PassengerTypeQuantity []PassengerTypeQuantity `json:"PassengerTypeQuantity"`
}

type PassengerTypeQuantity struct {
	Code     string `json:"Code"`
	Quantity int    `json:"Quantity"`
}

// TPA_Extensions related structs
type TPAExtensions struct {
	IntelliSellTransaction IntelliSellTransaction `json:"IntelliSellTransaction"`
}

type IntelliSellTransaction struct {
	RequestType RequestType `json:"RequestType"`
}

type RequestType struct {
	Name string `json:"Name"`
}
