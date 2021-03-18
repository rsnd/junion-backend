package config

import (
	"os"
	// Blank import for loading environment variables
	_ "github.com/joho/godotenv/autoload"
)

var common = map[string]string{
	"ENVIRONMENT":      os.Getenv("ENVIRONMENT"),
	"SENTRY_DSN":       os.Getenv("SENTRY_DSN"),
	"SENDGRID_APY_KEY": os.Getenv("SENDGRID_APY_KEY"),
}

// Config exports configuration variables from a .env file
var config = map[string](map[string]string){
	"test": map[string]string{
		"TEST_USER_USERNAME": os.Getenv("TEST_USER_USERNAME"),
		"TEST_USER_EMAIL":    os.Getenv("TEST_USER_EMAIL"),
		"TEST_USER_PASSWORD": os.Getenv("TEST_USER_PASSWORD"),
		"PORT":               os.Getenv("TEST_PORT"),
		"DATABASE_URL":       os.Getenv("TEST_DATABASE_URL"),
		"DATABASE_NAME":      os.Getenv("TEST_DATABASE_NAME"),
		"JWT_SECRET":         os.Getenv("DEV_JWT_SECRET"),
		"JWT_EXPIERSIN":      os.Getenv("DEV_JWT_EXPIERSIN"),
	},
	"local": map[string]string{
		"PORT":          os.Getenv("DEV_PORT"),
		"DATABASE_URL":  os.Getenv("LOCAL_DATABASE_URL"),
		"DATABASE_NAME": os.Getenv("LOCAL_DATABASE_NAME"),
		"JWT_SECRET":    os.Getenv("DEV_JWT_SECRET"),
		"JWT_EXPIERSIN": os.Getenv("DEV_JWT_EXPIERSIN"),
	},
	"development": map[string]string{
		// "PORT":          os.Getenv("DEV_PORT"),
		"PORT":          os.Getenv("PORT"),
		"DATABASE_URL":  os.Getenv("DEV_DATABASE_URL"),
		"DATABASE_NAME": os.Getenv("DEV_DATABASE_NAME"),
		"JWT_SECRET":    os.Getenv("DEV_JWT_SECRET"),
		"JWT_EXPIERSIN": os.Getenv("DEV_JWT_EXPIERSIN"),
	},
	"staging": map[string]string{
		"PORT":          os.Getenv("DEV_PORT"),
		"DATABASE_URL":  os.Getenv("STAGING_DATABASE_URL"),
		"DATABASE_NAME": os.Getenv("STAGING_DATABASE_NAME"),
		"JWT_SECRET":    os.Getenv("DEV_JWT_SECRET"),
		"JWT_EXPIERSIN": os.Getenv("DEV_JWT_EXPIERSIN"),
	},
	"production": map[string]string{
		"PORT":          os.Getenv("PORT"),
		"DATABASE_URL":  os.Getenv("PROD_DATABASE_URL"),
		"DATABASE_NAME": os.Getenv("PROD_DATABASE_NAME"),
		"JWT_SECRET":    os.Getenv("PROD_JWT_SECRET"),
		"JWT_EXPIERSIN": os.Getenv("PROD_JWT_EXPIERSIN"),
	},
}

// GetConfig return a map of environment variables
// based on the environment.
func GetConfig() map[string]string {
	currentConfig := config[os.Getenv("ENVIRONMENT")]
	for k, v := range common {
		currentConfig[k] = v
	}
	return currentConfig
}
