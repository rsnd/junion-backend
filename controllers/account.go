package controllers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rsnd/junion-backend/models"
	"github.com/rsnd/junion-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AccountFindOne is an example controller for the "/accounts" route
func AccountFindOne(c echo.Context) error {
	userID := c.Param("id")
	bsonObjectID, _ := primitive.ObjectIDFromHex(userID)

	user := new(models.User)
	err := models.UserCollection.FindOne(
		models.Ctx,
		bson.M{"_id": bsonObjectID},
	).Decode(&user)
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "This user does not exist")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, user)
}

// AccountUpdate is an example controller for the "/accounts" route
func AccountUpdate(c echo.Context) error {
	userID := c.Param("id")
	bsonObjectID, _ := primitive.ObjectIDFromHex(userID)
	updateUser := new(models.UpdateUser)
	utils.BindJSON(c, updateUser)

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	user := new(models.User)
	err := models.UserCollection.FindOneAndUpdate(
		models.Ctx,
		bson.M{"_id": bsonObjectID},
		bson.M{
			"$set": bson.M{
				"fullname": updateUser.Fullname,
			},
		},
		&opt,
	).Decode(&user)
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "This user does not exist")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(201, user)
}

// AccountRemove is an example controller for the "/accounts" route
func AccountRemove(c echo.Context) error {
	userID := c.Param("id")
	bsonObjectID, _ := primitive.ObjectIDFromHex(userID)

	_, err := models.UserCollection.DeleteOne(
		models.Ctx,
		bson.M{"_id": bsonObjectID},
	)
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "This user does not exist")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	return c.JSON(200, map[string]string{"message": "User deleted successfully"})
}
