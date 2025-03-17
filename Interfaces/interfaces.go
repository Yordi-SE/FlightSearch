package interfaces

import (
	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
)

type UseScase interface {
	SearchFlights(req *DTO.FlightSearchRequest) (*DTO.FlightSearchResponse, error)
	GetToken() error
}
