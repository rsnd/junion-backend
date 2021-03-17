package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rsnd/junion-backend/controllers"
)

// ConversationRoutes configures all endpoints for the conversations route
func ConversationRoutes(
	conversationGroup *echo.Group,
	conversationGroupRestricted *echo.Group,
) {
	conversationGroup.GET("/", controllers.ConversationFind)
	conversationGroupRestricted.POST("/", controllers.ConversationCreate)

	conversationGroup.GET("/action", controllers.ConversationActions)
	conversationGroup.GET("/:id", controllers.ConversationFindOne)

	conversationGroupRestricted.PATCH("/:id", controllers.ConversationUpdate)
	conversationGroupRestricted.DELETE("/:id", controllers.ConversationRemove)
}
