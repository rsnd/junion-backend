package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/rsnd/junion-backend/config"
	"github.com/rsnd/junion-backend/models"
)

// BindJSON return a server error
func BindJSON(c echo.Context, val interface{}) error {
	if err := c.Bind(val); err != nil {
		return ErrorHandler(500, "An error occured")
	}
	return nil
}

// ErrorHandler is a global error handler
func ErrorHandler(code int, message string) error {
	return echo.NewHTTPError(code, message)
}

// GenerateToken generate a JWT
func GenerateToken(user *models.User) (string, error) {

	// Create token
	tokenSeed := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := tokenSeed.Claims.(jwt.MapClaims)
	claims["data"] = user

	// Generate encoded token and send it as response.
	token, err := tokenSeed.SignedString([]byte(config.GetConfig()["JWT_SECRET"]))
	if err != nil {
		return "", ErrorHandler(500, "An error occured, pls try again.")
	}
	return token, nil
}
