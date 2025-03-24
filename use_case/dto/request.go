package DTO

import "fmt"

type Passenger struct {
	Type  string `json:"type" binding:"required,oneof=ADT CNN INF C06"` // Adult, Child, Infant
	Count int    `json:"count" binding:"required,min=1"`
}

type FlightSearchRequest struct {
	TripType          string      `json:"trip_type" binding:"required,oneof=one_way round_trip"`
	Origin            string      `json:"origin" binding:"required,len=3"`      // IATA code
	Destination       string      `json:"destination" binding:"required,len=3"` // IATA code
	DepartureDateTime string      `json:"departure_date" binding:"required"`    // Format: YYYY-MM-DD
	ReturnDateTime    string      `json:"return_date"`                          // Required for round_trip
	Passengers        []Passenger `json:"passengers" binding:"required,dive"`
}

// Validate ensures business rules (e.g., no child/infant traveling alone)
func (r *FlightSearchRequest) Validate() error {
	hasAdult := false
	for _, p := range r.Passengers {
		if p.Type == "ADT" {
			hasAdult = true
			break
		}
	}
	if !hasAdult && len(r.Passengers) > 0 {
		return fmt.Errorf("at least one adult (ADT) is required when traveling with children or infants")
	}
	if r.TripType == "round_trip" && r.ReturnDateTime == "" {
		return fmt.Errorf("return_date is required for round_trip")
	}
	return nil
}

type RequestLocation struct {
	LocationCode string `json:"LocationCode"`
	LocationType string `json:"LocationType"`
}

type OriginDest struct {
	OriginLocation      RequestLocation `json:"OriginLocation"`
	DestinationLocation RequestLocation `json:"DestinationLocation"`
	DepartureDateTime   string          `json:"DepartureDateTime"`
}

type OTA_AirLowFareSearchRQ struct {
	Version                      string              `json:"Version"`
	POS                          POS                 `json:"POS"`
	OriginDestinationInformation []OriginDest        `json:"OriginDestinationInformation"`
	TravelerInfoSummary          TravelerInfoSummary `json:"TravelerInfoSummary"`
	TPA_Extensions               TPAExtensions       `json:"TPA_Extensions"`
}
type SabreRequestFormat struct {
	OTA_AirLowFareSearchRQ OTA_AirLowFareSearchRQ `json:"OTA_AirLowFareSearchRQ"`
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
	DepartureDateTime   string          `json:"DepartureDateTime"`
	OriginLocation      RequestLocation `json:"OriginLocation"`
	DestinationLocation RequestLocation `json:"DestinationLocation"`
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
