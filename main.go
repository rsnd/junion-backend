package main

import (
	"github.com/rsnd/junion-backend/config"
	"github.com/rsnd/junion-backend/models"
	"github.com/rsnd/junion-backend/routes"
)

func main() {
	currentConfig := config.GetConfig();

	models.ConnectDB()
	app := routes.New()

	app.Logger.Fatal(app.Start(":" + currentConfig["PORT"]))
}