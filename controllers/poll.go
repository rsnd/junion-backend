package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rsnd/junion-backend/models"
	"github.com/rsnd/junion-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PollFind is a controller
func PollFind(c echo.Context) error {
	var polls []bson.M
	cursor, err := models.PollsCollection.Find(
		models.Ctx,
		bson.M{"event.id": c.QueryParam("eventID")},
	)
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return c.JSON(200, [0]string{})
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	if err = cursor.All(models.Ctx, &polls); err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	if len(polls) == 0 {
		return c.JSON(200, [0]string{})
	}
	return c.JSON(200, polls)
}

// PollCreate is a controller
func PollCreate(c echo.Context) error {
	pollData := new(models.Poll)
	if err := c.Bind(pollData); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}
	if err := c.Validate(pollData); err != nil {
		return utils.ErrorHandler(http.StatusBadRequest, fmt.Sprint(err))
	}

	jsonData, _ := json.Marshal(pollData)
	var poll map[string]interface{}
	json.Unmarshal(jsonData, &poll)

	res, err := models.PollsCollection.InsertOne(
		models.Ctx,
		poll,
	)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	poll["_id"] = res.InsertedID
	return c.JSON(201, poll)
}

// PollVote is a controller
func PollVote(c echo.Context) error {
	pollId, _ := primitive.ObjectIDFromHex(c.QueryParam("id"))
	optionPosition := c.QueryParam("position")

	updatedPoll := new(models.Poll)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := models.PollsCollection.FindOneAndUpdate(
		models.Ctx,
		bson.M{"_id": pollId},
		bson.M{
			"$inc": bson.M{
				"votes":                                1,
				"options." + optionPosition + ".votes": 1,
			},
		},
		&opt,
	).Decode(&updatedPoll)
	if err != nil {
		fmt.Println(err)
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, updatedPoll)
}

// PollFindOne is a controller
func PollFindOne(c echo.Context) error {
	pollID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	poll := new(models.Poll)
	err := models.PollsCollection.FindOne(
		models.Ctx,
		bson.M{"_id": pollID},
	).Decode(&poll)
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "Poll does not exist.")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, poll)
}

// PollUpdate is a controller
func PollUpdate(c echo.Context) error {
	pollData := new(models.PollUpdate)
	if err := c.Bind(pollData); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}
	if err := c.Validate(pollData); err != nil {
		return utils.ErrorHandler(http.StatusBadRequest, fmt.Sprint(err))
	}
	pollID, _ := primitive.ObjectIDFromHex(c.Param("id"))

	updatedPoll := new(models.Poll)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := models.PollsCollection.FindOneAndUpdate(
		models.Ctx,
		bson.M{"_id": pollID},
		bson.M{
			"$set": bson.M{
				"image":    pollData.Image,
				"question": pollData.Question,
				"options":  pollData.Options,
			},
		},
		&opt,
	).Decode(&updatedPoll)
	if err != nil {
		fmt.Println(err)
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, updatedPoll)
}

// PollRemove is a controller
func PollRemove(c echo.Context) error {
	pollID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	_, err := models.PollsCollection.DeleteOne(
		models.Ctx,
		bson.M{"_id": pollID},
	)

	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "Event does not exist.")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, map[string]string{"message": "Poll deleted successfully"})
}
