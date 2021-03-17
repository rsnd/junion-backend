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

// ConversationFind is a controller
func ConversationFind(c echo.Context) error {
	eventObjectID, _ := primitive.ObjectIDFromHex(c.QueryParam("eventID"))

	var conversations []bson.M
	cursor, err := models.ConversationsCollection.Find(
		models.Ctx,
		bson.M{"eventId": eventObjectID},
	)
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return c.JSON(200, [0]string{})
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	if err = cursor.All(models.Ctx, &conversations); err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	if len(conversations) == 0 {
		return c.JSON(200, [0]string{})
	}
	return c.JSON(200, conversations)
}

// ConversationCreate is a controller
func ConversationCreate(c echo.Context) error {
	conversationData := new(models.Conversation)
	if err := c.Bind(conversationData); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}
	if err := c.Validate(conversationData); err != nil {
		return utils.ErrorHandler(http.StatusBadRequest, fmt.Sprint(err))
	}

	jsonData, _ := json.Marshal(conversationData)
	var conversation map[string]interface{}
	json.Unmarshal(jsonData, &conversation)

	res, err := models.ConversationsCollection.InsertOne(
		models.Ctx,
		conversation,
	)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	conversation["_id"] = res.InsertedID
	return c.JSON(201, conversation)
}

// ConversationActions is a controller
func ConversationActions(c echo.Context) error {
	conversationID, _ := primitive.ObjectIDFromHex(c.QueryParam("id"))
	actionType := c.QueryParam("type")

	like := 1
	conversation := new(models.Conversation)
	err := models.ConversationsCollection.FindOne(
		models.Ctx,
		bson.M{"_id": conversationID},
	).Decode(&conversation)

	if actionType == "dislike" {
		like = -1
	}
	if conversation.Likes == 0 && actionType == "dislike" {
		like = 0
	}

	updatedConversation := new(models.Conversation)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err = models.ConversationsCollection.FindOneAndUpdate(
		models.Ctx,
		bson.M{"_id": conversationID},
		bson.M{
			"$inc": bson.M{
				"likes": like,
			},
		},
		&opt,
	).Decode(&updatedConversation)
	if err != nil {
		fmt.Println(err)
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, updatedConversation)
}

// ConversationFindOne is a controller
func ConversationFindOne(c echo.Context) error {
	conversationID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	conversation := new(models.Conversation)
	err := models.ConversationsCollection.FindOne(
		models.Ctx,
		bson.M{"_id": conversationID},
	).Decode(&conversation)
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "Conversation does not exist.")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, conversation)
}

// ConversationUpdate is a controller
func ConversationUpdate(c echo.Context) error {
	conversationData := new(models.ConversationUpdate)
	if err := c.Bind(conversationData); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}
	if err := c.Validate(conversationData); err != nil {
		return utils.ErrorHandler(http.StatusBadRequest, fmt.Sprint(err))
	}
	conversationID, _ := primitive.ObjectIDFromHex(c.Param("id"))

	updatedConversation := new(models.Conversation)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := models.ConversationsCollection.FindOneAndUpdate(
		models.Ctx,
		bson.M{"_id": conversationID},
		bson.M{
			"$set": bson.M{
				"text": conversationData.Text,
			},
			"$push": bson.M{
				"replies": conversationData.Replies,
			},
		},
		&opt,
	).Decode(&updatedConversation)
	if err != nil {
		fmt.Println(err)
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, updatedConversation)
}

// ConversationRemove is a controller
func ConversationRemove(c echo.Context) error {
	conversationID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	_, err := models.ConversationsCollection.DeleteOne(
		models.Ctx,
		bson.M{"_id": conversationID},
	)

	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "Event does not exist.")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, map[string]string{"message": "Conversation deleted successfully"})
}
