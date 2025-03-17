package controller

import (
	"fmt"

	interfaces "github.com/Yordi-SE/FlightSearch/Interfaces"
	DTO "github.com/Yordi-SE/FlightSearch/use_case/dto"
	"github.com/gin-gonic/gin"
)

// SearchFlights is a controller that handles flight search requests
type Controller struct {
	FlightClient interfaces.UseScase
}

// NewController returns a new Controller instance
func NewController(client interfaces.UseScase) *Controller {
	return &Controller{FlightClient: client}
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
	result, err := ctrl.FlightClient.SearchFlights(&req)
	fmt.Println(result)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)

}
