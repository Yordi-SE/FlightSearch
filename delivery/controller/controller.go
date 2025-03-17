package controller

import (
	"fmt"

	"github.com/Yordi-SE/FlightSearch/use_case"
	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
	"github.com/gin-gonic/gin"
)

// SearchFlights is a controller that handles flight search requests
type Controller struct {
	SabreClient *use_case.SabreClient
}

// NewController returns a new Controller instance
func NewController(client *use_case.SabreClient) *Controller {
	return &Controller{SabreClient: client}
}

func (ctrl *Controller) SearchFlights(c *gin.Context) {
	var req DTO.FlightSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := req.Validate()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	result, err := ctrl.SabreClient.SearchFlights(&req)
	fmt.Println(result)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)

}
