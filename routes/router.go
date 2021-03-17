package routes

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rsnd/junion-backend/config"
	"github.com/rsnd/junion-backend/models"
)

func handleBaseRoute(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the Junion API")
}

// New returns router for app routes configuration.
func New() *echo.Echo {
	currentConfig := config.GetConfig()
	jwtSecret := currentConfig["JWT_SECRET"]

	// Setup base router
	router := echo.New()

	// Initialize request validators
	router.Validator = &models.CustomValidator{Validator: validator.New()}

	// Create route groups
	v1Group := router.Group("/v1")

	accountGroup := v1Group.Group("/accounts")
	accountGroupRestricted := v1Group.Group("/accounts")

	eventGroup := v1Group.Group("/events")
	eventGroupRestricted := v1Group.Group("/events")

	conversationGroup := v1Group.Group("/conversations")
	conversationGroupRestricted := v1Group.Group("/conversations")

	pollGroup := v1Group.Group("/polls")
	pollGroupRestricted := v1Group.Group("/polls")

	// Attach general middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Gzip())
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// Route specific middleware
	accountGroupRestricted.Use(middleware.JWT([]byte(jwtSecret)))
	eventGroupRestricted.Use(middleware.JWT([]byte(jwtSecret)))
	conversationGroupRestricted.Use(middleware.JWT([]byte(jwtSecret)))
	pollGroupRestricted.Use(middleware.JWT([]byte(jwtSecret)))

	// Initialize Routes
	v1Group.GET("", handleBaseRoute)
	v1Group.GET("/", handleBaseRoute)

	AccountRoutes(accountGroup, accountGroupRestricted)
	EventRoutes(eventGroup, eventGroupRestricted)
	ConversationRoutes(conversationGroup, conversationGroupRestricted)
	PollRoutes(pollGroup, pollGroupRestricted)

	// Handle general routes
	router.GET("/", func(c echo.Context) error {
		user := c.Get("user")
		fmt.Println("user ", user)
		return c.String(http.StatusOK, "Refer to the route /v1'")
	})
	router.GET("*", func(c echo.Context) error {
		return c.String(http.StatusNotFound, "A beast ate this page...")
	})

	return router
}
