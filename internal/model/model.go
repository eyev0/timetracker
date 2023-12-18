package model

import (
	"time"

	"github.com/google/uuid"
)

type Result[T any] struct {
	Value T
	Err   error
}

type User struct {
	Id         string    `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	TimeZone   string    `db:"time_zone"`
	CalendarId *string   `db:"calendar_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Token struct {
	Id           string    `db:"id" json:"id"`
	UserId       string    `db:"user_id" json:"user_id"`
	AccessToken  string    `db:"access_token" json:"access_token"`
	IdToken      string    `db:"id_token" json:"id_token"`
	ExpiresIn    int       `db:"expires_in" json:"expires_in"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	Scope        string    `db:"scope" json:"scope"`
	TokenType    string    `db:"token_type" json:"token_type"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Entry struct {
	Id             uuid.UUID  `db:"id" json:"id,omitempty"`
	UserId         string     `db:"user_id" json:"user_id"`
	StartDateTime  time.Time  `db:"start_timestamp" json:"start_datetime,omitempty"`
	ElapsedSeconds int        `json:"elapsed_seconds"`
	ElapsedMinutes int        `json:"elapsed_minutes"`
	ElapsedHours   int        `json:"elapsed_hours"`
	EndDateTime    *time.Time `db:"end_timestamp" json:"end_datetime,omitempty"`
	Note           string     `db:"note" json:"note"`
	CalendarId     string     `db:"calendar_id" json:"calendar_id"`
}

func (entry *Entry) CalcElapsed(now time.Time) {
	elapsed := now.Sub(entry.StartDateTime)
	entry.ElapsedSeconds = int(elapsed/time.Second) % 60
	entry.ElapsedMinutes = int(elapsed/time.Minute) % 60
	entry.ElapsedHours = int(elapsed / time.Hour)
}

type CreateEntryInput struct {
	Note string `json:"note" type:"string" doc:"Description of activity"`
}

type UpdateEntryInput struct {
	Id          *string    `json:"id,omitempty" type:"string"`
	Note        *string    `json:"note,omitempty" type:"string"`
	EndDateTime *time.Time `json:"end_datetime,omitempty" type:"string"`
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"update_at,omitempty"`
}

func FilteredResponse(user *User) UserResponse {
	return UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type Settings struct {
	CalendarId string `json:"calendar_id,omitempty"`
}

type SettingsInput = Settings
