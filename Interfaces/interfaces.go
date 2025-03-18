package interfaces

import (
	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
)

type UseCase interface {
	SearchFlights(req *DTO.FlightSearchRequest) (*DTO.FlightSearchResponse, error)
	GetToken() error
}
