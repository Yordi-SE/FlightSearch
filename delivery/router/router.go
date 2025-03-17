package router

import (
	"github.com/Yordi-SE/FlightSearch/delivery/controller"
	"github.com/Yordi-SE/FlightSearch/use_case"
	"github.com/gin-gonic/gin"
)

// New returns a new Router instance
func NewRouter(Sabre *use_case.SabreClient) {
	router := gin.Default()

	Controller := controller.NewController(Sabre)
	router.POST("/flight/search", Controller.SearchFlights)
	router.Run(":8080")

}
