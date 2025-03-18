package controller

import (
	"fmt"

	interfaces "github.com/Yordi-SE/FlightSearch/Interfaces" // Package defining use case interfaces
	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"      // Package containing data transfer objects
	"github.com/gin-gonic/gin"                               // Gin web framework for HTTP handling
)

// Controller handles HTTP requests related to flight searches
type Controller struct {
	FlightClient interfaces.UseCase // Interface for interacting with flight search use case
}

// NewController creates and initializes a new Controller instance
// Args:
//
//	client - An implementation of the UseScase interface for flight operations
//
// Returns:
//
//	Pointer to a new Controller instance
func NewController(client interfaces.UseCase) *Controller {
	return &Controller{
		FlightClient: client, // Inject the flight client dependency
	}
}

// SearchFlights handles the HTTP POST request to search for flights
// It parses the request, validates it, and returns flight search results
// Args:
//
//	c - Gin context containing the HTTP request and response
func (ctrl *Controller) SearchFlights(c *gin.Context) {
	var req DTO.FlightSearchRequest

	// Bind JSON request body to FlightSearchRequest struct
	// Returns 400 Bad Request if binding fails (e.g., invalid JSON)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate the request data (e.g., required fields, formats)
	// Returns 400 Bad Request if validation fails
	err := req.Validate()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return // Note: This return is redundant due to the missing 'return' in the original
	}

	// Call the use case to search for flights
	result, err := ctrl.FlightClient.SearchFlights(&req)
	fmt.Println(result) // Log the result for debugging purposes

	// If the search fails, return 500 Internal Server Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Success: return 200 OK with the flight search results
	c.JSON(200, result)
}
