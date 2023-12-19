package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"

	"github.com/eyev0/timetracker/internal/cfg"
	"github.com/eyev0/timetracker/internal/db"
	"github.com/eyev0/timetracker/internal/google_api"
	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
	"github.com/eyev0/timetracker/internal/utils"
)

func SignUpUser(ctx *gin.Context) {
	var payload *model.RegisterUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := model.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  payload.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := db.CreateUser(&newUser)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
			return
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Internal error"})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": model.FilteredResponse(&newUser)}})
}

func SignInUser(ctx *gin.Context) {
	var payload *model.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := db.GetUserByEmail(strings.ToLower(payload.Email))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if user.Password != "" && user.Password != payload.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	log.Logger.Debugf("Expires in: %+v", cfg.C.TokenExpiresIn)

	token, err := utils.GenerateJwt(cfg.C.TokenExpiresIn, user.Id, cfg.C.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, int(cfg.C.TokenExpiresIn/time.Second), "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

func LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GoogleOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
		return
	}

	googleOauthToken, err := utils.GetGoogleOauthToken(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	googleUser, err := utils.GetGoogleUser(googleOauthToken.AccessToken, googleOauthToken.IdToken)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	email := strings.ToLower(googleUser.Email)

	user_data := model.User{
		Name:      googleUser.Name,
		Email:     email,
		Password:  "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	user, err := db.GetUserByEmail(email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			db.CreateUser(&user_data)
			user, err = db.GetUserByEmail(email)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
	}

	err = db.UpsertUserToken(user, googleOauthToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	token, err := utils.GenerateJwt(cfg.C.TokenExpiresIn, user.Id, cfg.C.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, cfg.C.TokenMaxAge*60, "/", "localhost", false, true)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(cfg.C.FrontEndOrigin, pathUrl))
}

func RefreshToken(ctx *gin.Context) {
	user := ctx.MustGet("currentUser").(*model.User)
	var err error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	creds, err := os.ReadFile("credentials.json")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	config, err := google.ConfigFromJSON(creds, calendar.CalendarEventsScope)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	oauth2Token, err := google_api.GetOAuthToken(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	oauth2Token, err = google_api.RefreshToken(config, user, oauth2Token.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	log.Logger.Infof("New Access Token: %s", oauth2Token)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"token": *oauth2Token}})
}

func GetToken(ctx *gin.Context) {
	user := ctx.MustGet("currentUser").(*model.User)

	oauth2Token, err := google_api.GetOAuthToken(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"token": *oauth2Token}})
}
