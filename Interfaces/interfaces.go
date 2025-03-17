package interfaces

import (
	"github.com/Yordi-SE/FlightSearch/models"
	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
)

type UseScase interface {
	SearchFlights(req *DTO.FlightSearchRequest) ([]models.Flight, error)
	GetToken() error
}
