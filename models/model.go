package models




// BaggageInfo represents baggage allowance for a passenger type
type BaggageInfo struct {
	PassengerType string `json:"passenger_type"` 
	Allowance     string `json:"allowance"`      
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
