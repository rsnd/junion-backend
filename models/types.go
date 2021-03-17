package models

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PasswordRecovery struct
type PasswordRecovery struct {
	Token  int       `json:"token"`
	Expiry time.Time `json:"expiry"`
}

// User struct
type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Fullname  string             `json:"fullname" bson:"fullname"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	IsAdmin   bool               `json:"isAdmin" bson:"isAdmin"`
	Password  string             `json:"-" bson:"password"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`

	PasswordRecovery PasswordRecovery `json:"-" bson:"passwordRecovery,omitempty"`
}

// SigninCred struct for signing in.
type SigninCred struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SigninPayload response after signin successful.
type SigninPayload struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// SignupVerCred struct for signup verification.
type SignupVerCred struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
}

// EmailVerification struct for email verification.
type EmailVerification struct {
	Email            string `json:"email" bson:"email"`
	VerificationCode int    `json:"verificationCode" bson:"verificationCode"`
}

// SignupCred struct for user registration.
type SignupCred struct {
	Fullname string `json:"fullname" bson:"fullname"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	IsAdmin  bool   `json:"isAdmin" bson:"isAdmin"`
	Password string `json:"password" bson:"password"`
	Code     string `json:"code" bson:"code"`
}

// Email struct.
type Email struct {
	Email string `json:"email" bson:"email"`
}

// UpdateUser struct.
type UpdateUser struct {
	Fullname string `json:"fullname" bson:"fullname"`
}

// PasswordResetCred struct for reseting a user's password.
type PasswordResetCred struct {
	Email       string `json:"email" bson:"email"`
	NewPassword string `json:"newPassword" bson:"newPassword"`
	Code        string `json:"code" bson:"code"`
}

// UpcomingEvent struct.
type UpcomingEvent struct {
	Title string    `json:"title" bson:"title" validate:"required"`
	Date  time.Time `json:"date" bson:"date" validate:"required"`
}

// SocialLinks struct.
type SocialLinks struct {
	Facebook   string `json:"facebook" bson:"facebook"`
	Twitter    string `json:"twitter" bson:"twitter"`
	Linkedin   string `json:"linkedin" bson:"linkedin"`
	Eventbrite string `json:"eventbrite" bson:"eventbrite"`
}

// Event struct.
type Event struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	Title          string             `json:"title" bson:"title" validate:"required"`
	Description    string             `json:"description" bson:"description" validate:"required"`
	Passcode       string             `json:"passcode" bson:"passcode"`
	Date           time.Time          `json:"date" bson:"date" validate:"required"`
	URL            string             `json:"url" bson:"url"`
	AudeinceSize   int                `json:"audeinceSize" bson:"audeinceSize"`
	CreatedBy      primitive.ObjectID `json:"createdBy" bson:"createdBy" validate:"required"`
	UpcomingEvents []UpcomingEvent    `json:"upcomingEvents" bson:"upcomingEvents"`
	SocialLinks    SocialLinks        `json:"socialLinks" bson:"socialLinks"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
}

// UpdateEvent struct.
type UpdateEvent struct {
	Title          string          `json:"title" bson:"title" validate:"required"`
	Description    string          `json:"description" bson:"description" validate:"required"`
	Date           time.Time       `json:"date" bson:"date" validate:"required"`
	URL            string          `json:"url" bson:"url"`
	UpcomingEvents []UpcomingEvent `json:"upcomingEvents" bson:"upcomingEvents"`
	SocialLinks    SocialLinks     `json:"socialLinks" bson:"socialLinks"`
}

// EventActions struct.
type EventActions struct {
	Type     string `json:"type" query:"type"`
	ID       string `json:"id" query:"id"`
	Passcode string `json:"passcode" query:"passcode"`
}

// ConversationCreatedBy struct.
type ConversationCreatedBy struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

// ConversationReply struct.
type ConversationReply struct {
	Text      string                `json:"text" bson:"text"`
	Date      time.Time             `json:"date" bson:"date"`
	CreatedBy ConversationCreatedBy `json:"createdBy" bson:"createdBy"`
}

// Conversation struct.
type Conversation struct {
	Text      string                `json:"text" bson:"text" validate:"required"`
	Likes     int                   `json:"likes" bson:"likes"`
	Date      time.Time             `json:"date" bson:"date" validate:"required"`
	EventId   string                `json:"eventId" bson:"eventId" validate:"required"`
	CreatedBy ConversationCreatedBy `json:"createdBy" bson:"createdBy" validate:"required"`
	Replies   []ConversationReply   `json:"replies" bson:"replies" validate:"required"`
}

// ConversationUpdate struct.
type ConversationUpdate struct {
	Text    string            `json:"text" bson:"text" validate:"required"`
	Replies ConversationReply `json:"replies" bson:"replies" validate:"required"`
}

// PollEvent struct.
type PollEvent struct {
	ID        string `json:"id" bson:"id" validate:"required"`
	CreatedBy string `json:"createdBy" bson:"createdBy" validate:"required"`
}

// PollOption struct.
type PollOption struct {
	Position int    `json:"postion" bson:"position"`
	Text     string `json:"text" bson:"text"`
	Votes    int    `json:"votes" bson:"votes"`
}

// Poll struct.
type Poll struct {
	Image     string       `json:"image" bson:"image" validate:"required"`
	Question  string       `json:"question" bson:"question" validate:"required"`
	Event     PollEvent    `json:"event" bson:"event" validate:"required"`
	Votes     int          `json:"votes" bson:"votes"`
	Options   []PollOption `json:"options" bson:"options" validate:"required"`
	CreatedAt time.Time    `json:"createdAt" bson:"createdAt"`
}

// PollUpdate struct.
type PollUpdate struct {
	Image    string       `json:"image" bson:"image" validate:"required"`
	Question string       `json:"question" bson:"question" validate:"required"`
	Options  []PollOption `json:"options" bson:"options" validate:"required"`
}

// CustomValidator is a custom request validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate is CustomValidator validator function
func (cv *CustomValidator) Validate(i interface{}) error {
	errorRes, err := cv.Validator.Struct(i).(validator.ValidationErrors)
	if !err {
		return nil
	}
	errorFields := []string{}
	for _, k := range errorRes {
		errorFields = append(errorFields, k.StructField())
	}
	if len(errorFields) == 1 {
		return errors.New(strings.Join(errorFields, ", ") + " field is invalid or missing.")
	}
	return errors.New(strings.Join(errorFields, ", ") + " fields are invalid or missing.")
}
