package router

import (
	interfaces "github.com/Yordi-SE/FlightSearch/Interfaces"
	"github.com/Yordi-SE/FlightSearch/delivery/controller"
	"github.com/gin-gonic/gin"
)

// New returns a new Router instance
func NewRouter(Sabre interfaces.UseScase) {
	router := gin.Default()

	Controller := controller.NewController(Sabre)
	router.POST("/flight/search", Controller.SearchFlights)
	router.Run(":8080")

}
