package controllers

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/rsnd/junion-backend/models"
	"github.com/rsnd/junion-backend/utils"
)

// AuthSignIn handles user sign in.
func AuthSignIn(c echo.Context) error {
	signinCred := new(models.SigninCred)
	if err := c.Bind(signinCred); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}

	if err := c.Validate(signinCred); err != nil {
		return utils.ErrorHandler(http.StatusBadRequest, fmt.Sprint(err))
	}

	user := new(models.User)
	err := models.UserCollection.FindOne(
		models.Ctx,
		bson.M{"email": signinCred.Email},
	).Decode(&user)

	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "This user does not exist")
	}
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(signinCred.Password),
	)

	if err != nil {
		return utils.ErrorHandler(401, "Password is incorrect")
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	res := &models.SigninPayload{
		Token: token,
		User:  *user,
	}

	return c.JSON(201, res)
}

// AuthSignupVerification guards against duplicate user signups
func AuthSignupVerification(c echo.Context) error {
	signupVerCred := new(models.SignupVerCred)
	if err := c.Bind(&signupVerCred); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}
	if err := c.Validate(signupVerCred); err != nil {
		return utils.ErrorHandler(http.StatusBadRequest, fmt.Sprint(err))
	}

	user := new(models.User)
	err := models.UserCollection.FindOne(
		models.Ctx,
		bson.M{
			"$or": []bson.M{
				{"email": strings.ToLower(signupVerCred.Email)},
				{"username": strings.ToLower(signupVerCred.Username)},
			},
		},
	).Decode(&user)
	if err != nil && fmt.Sprint(err) != "mongo: no documents in result" {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	if err == nil {
		return utils.ErrorHandler(409, "This email or username already exists.")
	}

	userEmail := strings.ToLower(signupVerCred.Email)
	verifEntity := new(models.EmailVerification)
	err = models.EmailverificationsCollection.FindOne(
		models.Ctx,
		bson.M{"email": userEmail},
	).Decode(&verifEntity)

	verificationCode := int(math.Floor(rand.Float64()*(9999-1000+1)) + 1000)

	_, err = models.EmailverificationsCollection.ReplaceOne(
		models.Ctx,
		bson.M{"email": userEmail},
		bson.M{
			"email":            userEmail,
			"verificationCode": verificationCode,
		},
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	return c.JSON(201, map[string]string{
		"code":    strconv.Itoa(verificationCode),
		"message": "Verification code sent",
	})
}

// AuthSignup registers a new user
func AuthSignup(c echo.Context) error {
	signupCred := new(models.SignupCred)
	if err := c.Bind(&signupCred); err != nil {
		return utils.ErrorHandler(500, "An error occured")
	}
	if err := c.Validate(signupCred); err != nil {
		return utils.ErrorHandler(http.StatusBadRequest, fmt.Sprint(err))
	}

	userEmail := strings.ToLower(signupCred.Email)
	verificationCode, err := strconv.Atoi(signupCred.Code)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	emailVerif := new(models.EmailVerification)
	err = models.EmailverificationsCollection.FindOne(
		models.Ctx,
		bson.M{
			"email":            userEmail,
			"verificationCode": verificationCode,
		},
	).Decode(&emailVerif)
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "Email verification unsuccessful.")
	}
	// fmt.Println("Reached-------")

	_, err = models.EmailverificationsCollection.DeleteOne(
		models.Ctx,
		bson.M{"email": userEmail},
	)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signupCred.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	user := bson.M{
		"fullname":  signupCred.Fullname,
		"username":  strings.ToLower(signupCred.Username),
		"email":     userEmail,
		"isAdmin":   false,
		"password":  string(passwordHash),
		"createdAt": time.Now().UTC(),
		"updatedAt": time.Now().UTC(),
	}
	_, err = models.UserCollection.InsertOne(
		models.Ctx,
		user,
	)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	delete(user, "password")

	return c.JSON(201, user)
}

// AuthPasswordResetVerification verifies a password reset
func AuthPasswordResetVerification(c echo.Context) error {
	userEmail := new(models.Email)
	utils.BindJSON(c, userEmail)

	fmt.Println("userEmail", userEmail)
	user := new(models.User)
	err := models.UserCollection.FindOne(
		models.Ctx,
		bson.M{"email": strings.ToLower(userEmail.Email)},
	).Decode(&user)
	if err != nil && fmt.Sprint(err) != "mongo: no documents in result" {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	if fmt.Sprint(err) == "mongo: no documents in result" {
		return utils.ErrorHandler(404, "User not found")
	}

	token := int(math.Floor(rand.Float64()*(9999-1000+1)) + 1000)
	passwordRecovery := bson.M{
		"token":  token,
		"expiry": time.Now().Local().Add(360000),
	}

	_, err = models.UserCollection.UpdateOne(
		models.Ctx,
		bson.M{"email": strings.ToLower(userEmail.Email)},
		bson.M{
			"$set": bson.M{"passwordRecovery": passwordRecovery},
		},
	)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	return c.JSON(201, map[string]interface{}{
		"message": "Verification code sent",
		"token":   token,
	})
}

// AuthPasswordReset resets a user's password
func AuthPasswordReset(c echo.Context) error {
	passwordResetCred := new(models.PasswordResetCred)
	utils.BindJSON(c, passwordResetCred)

	user := new(models.User)
	err := models.UserCollection.FindOne(
		models.Ctx,
		bson.M{"email": strings.ToLower(passwordResetCred.Email)},
	).Decode(&user)
	if err != nil {
		return utils.ErrorHandler(404, "Verification code not found or incorrect.")
	}

	if user.PasswordRecovery.Expiry.After(time.Now()) {
		return utils.ErrorHandler(500, "Sorry, your reset code has expired, please resend another.")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwordResetCred.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}
	_, err = models.UserCollection.UpdateOne(
		models.Ctx,
		bson.M{"email": strings.ToLower(passwordResetCred.Email)},
		bson.M{
			"$set": bson.M{"password": string(passwordHash)},
		},
	)
	if err != nil {
		return utils.ErrorHandler(500, "An error occured, pls try again.")
	}

	return c.JSON(201, map[string]interface{}{
		"message": "Password reset successful",
	})
}
