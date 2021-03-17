package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rsnd/junion-backend/controllers"
)

// AccountRoutes configures all endpoints for the accounts route
func AccountRoutes(
	accountGroup *echo.Group,
	accountGroupRestricted *echo.Group,
) {
	accountGroup.POST("/signin", controllers.AuthSignIn)
	accountGroup.POST("/signup/verification", controllers.AuthSignupVerification)
	accountGroup.POST("/signup", controllers.AuthSignup)
	accountGroup.POST("/forgotpassword/verification", controllers.AuthPasswordResetVerification)
	accountGroup.POST("/forgotpassword", controllers.AuthPasswordReset)

	accountGroupRestricted.GET("/:id", controllers.AccountFindOne)
	accountGroupRestricted.PUT("/:id", controllers.AccountUpdate)
	accountGroupRestricted.DELETE("/:id", controllers.AccountRemove)
}
