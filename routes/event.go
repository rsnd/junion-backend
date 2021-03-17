package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rsnd/junion-backend/controllers"
)

// EventRoutes configures all endpoints for the events route
func EventRoutes(
	eventGroup *echo.Group,
	eventGroupRestricted *echo.Group,
) {
	eventGroupRestricted.GET("/", controllers.EventFind)
	eventGroupRestricted.POST("/", controllers.EventCreate)

	eventGroup.GET("/action", controllers.EventActions)
	eventGroup.GET("/:id", controllers.EventFindOne)

	eventGroupRestricted.PUT("/:id", controllers.EventUpdate)
	eventGroupRestricted.DELETE("/:id", controllers.EventRemove)
}
