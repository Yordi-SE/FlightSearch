package router

import (
	interfaces "github.com/Yordi-SE/FlightSearch/Interfaces"
	"github.com/Yordi-SE/FlightSearch/delivery/controller"
	"github.com/gin-gonic/gin"
)

// New returns a new Router instance
func NewRouter(FlightClient interfaces.UseScase) {
	router := gin.Default()

	Controller := controller.NewController(FlightClient)
	router.POST("/flight/search", Controller.SearchFlights)
	router.Run(":8080")

}
