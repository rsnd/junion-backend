package controllers

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/rsnd/junion-backend/models"
	"github.com/rsnd/junion-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EventFind is a controller
func EventFind(c echo.Context) error {
	user := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userClaims := user["data"].(map[string]interface{})
	bsonObjectID, _ := primitive.ObjectIDFromHex(userClaims["_id"].(string))

	var events []bson.M
	cursor, err := models.EventsCollection.Find(
		models.Ctx,
		bson.M{"createdBy": bsonObjectID},
	)
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return c.JSON(200, [0]string{})
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	if err = cursor.All(models.Ctx, &events); err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	if len(events) == 0 {
		return c.JSON(200, [0]string{})
	}

	return c.JSON(200, events)
}

// EventCreate is a controller
func EventCreate(c echo.Context) error {
	eventData := new(models.Event)
	if err := c.Bind(eventData); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}
	if err := c.Validate(eventData); err != nil {
		return utils.ErrorHandler(http.StatusBadRequest, fmt.Sprint(err))
	}

	existingEvent := new(models.Event)
	userClaims := c.Get("user").(*jwt.Token).
		Claims.(jwt.MapClaims)["data"].(map[string]interface{})
	bsonObjectID, _ := primitive.ObjectIDFromHex(userClaims["_id"].(string))
	err := models.EventsCollection.FindOne(
		models.Ctx,
		bson.M{
			"createdBy": bsonObjectID,
			"title":     eventData.Title,
		},
	).Decode(&existingEvent)
	if existingEvent.Title != "" {
		return c.JSON(409, "This event already exists.")
	}

	passcode := int(math.Floor(rand.Float64()*(9999-1000+1)) + 1000)
	newEvent := bson.M{
		"title":          eventData.Title,
		"description":    eventData.Description,
		"passcode":       strconv.Itoa(passcode),
		"date":           eventData.Date,
		"url":            eventData.URL,
		"audeinceSize":   0,
		"createdBy":      bsonObjectID,
		"upcomingEvents": eventData.UpcomingEvents,
		"socialLinks":    eventData.SocialLinks,
		"createdAt":      time.Now().UTC(),
		"updatedAt":      time.Now().UTC(),
	}
	res, err := models.EventsCollection.InsertOne(
		models.Ctx,
		newEvent,
	)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	newEvent["_id"] = res.InsertedID

	return c.JSON(201, newEvent)
}

// EventActions is a controller
func EventActions(c echo.Context) error {
	queryParams := new(models.EventActions)
	if err := c.Bind(queryParams); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}

	event := new(models.Event)
	bsonObjectID, _ := primitive.ObjectIDFromHex(queryParams.ID)
	err := models.EventsCollection.FindOne(
		models.Ctx,
		bson.M{"_id": bsonObjectID},
	).Decode(&event)

	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "Event does not exist.")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	if event.Passcode != queryParams.Passcode {
		return utils.ErrorHandler(403, "Incorrect passcode.")
	}

	if time.Now().After(event.Date) {
		return utils.ErrorHandler(403, "This event has expired.")
	}

	operation := 1
	updatedEvent := new(models.Event)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	if queryParams.Type == "leave" {
		operation = -1
	}
	if event.AudeinceSize == 0 && queryParams.Type == "leave" {
		operation = 0
	}

	err = models.EventsCollection.FindOneAndUpdate(
		models.Ctx,
		bson.M{"_id": bsonObjectID},
		bson.M{
			"$inc": bson.M{
				"audeinceSize": operation,
			},
		},
		&opt,
	).Decode(&updatedEvent)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	return c.JSON(200, updatedEvent)
}

// EventFindOne is a controller
func EventFindOne(c echo.Context) error {
	eventID := c.Param("id")
	event := new(models.Event)
	bsonObjectID, _ := primitive.ObjectIDFromHex(eventID)
	err := models.EventsCollection.FindOne(
		models.Ctx,
		bson.M{"_id": bsonObjectID},
	).Decode(&event)

	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "Event does not exist.")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, event)
}

// EventUpdate is a controller
func EventUpdate(c echo.Context) error {
	eventData := new(models.UpdateEvent)
	if err := c.Bind(eventData); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}
	if err := c.Validate(eventData); err != nil {
		return utils.ErrorHandler(http.StatusBadRequest, fmt.Sprint(err))
	}

	eventID := c.Param("id")
	event := new(models.Event)
	eventObjectID, _ := primitive.ObjectIDFromHex(eventID)
	userClaims := c.Get("user").(*jwt.Token).
		Claims.(jwt.MapClaims)["data"].(map[string]interface{})
	userObjectID, _ := primitive.ObjectIDFromHex(userClaims["_id"].(string))
	err := models.EventsCollection.FindOne(
		models.Ctx,
		bson.M{
			"_id":       eventObjectID,
			"createdBy": userObjectID,
		},
	).Decode(&event)

	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "Event does not exist.")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	if time.Now().After(event.Date) {
		return utils.ErrorHandler(403, "This event has expired.")
	}

	updatedEvent := new(models.Event)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err = models.EventsCollection.FindOneAndUpdate(
		models.Ctx,
		bson.M{"_id": eventObjectID},
		bson.M{
			"$set": bson.M{
				"title":          eventData.Title,
				"description":    eventData.Description,
				"date":           eventData.Date,
				"url":            eventData.URL,
				"upcomingEvents": eventData.UpcomingEvents,
				"socialLinks":    eventData.SocialLinks,
			},
		},
		&opt,
	).Decode(&updatedEvent)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	return c.JSON(200, updatedEvent)
}

// EventRemove is a controller
func EventRemove(c echo.Context) error {
	eventID := c.Param("id")

	bsonObjectID, _ := primitive.ObjectIDFromHex(eventID)
	_, err := models.EventsCollection.DeleteOne(
		models.Ctx,
		bson.M{"_id": bsonObjectID},
	)

	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "Event does not exist.")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, map[string]string{"message": "Event deleted successfully"})
}
