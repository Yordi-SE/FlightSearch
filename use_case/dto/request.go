package DTO

import "fmt"

type Passenger struct {
	Type  string `json:"type" binding:"required,oneof=ADT CNN INF"` // Adult, Child, Infant
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
