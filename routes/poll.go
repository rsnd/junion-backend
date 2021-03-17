package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rsnd/junion-backend/controllers"
)

// PollRoutes configures all endpoints for the poll routes
func PollRoutes(
	pollGroup *echo.Group,
	pollGroupRestricted *echo.Group,
) {
	pollGroup.GET("/", controllers.PollFind)
	pollGroupRestricted.POST("/", controllers.PollCreate)

	pollGroup.GET("/vote", controllers.PollVote)
	pollGroup.GET("/:id", controllers.PollFindOne)

	pollGroupRestricted.PUT("/:id", controllers.PollUpdate)
	pollGroupRestricted.DELETE("/:id", controllers.PollRemove)
}
