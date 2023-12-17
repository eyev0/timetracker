package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/eyev0/timetracker/internal/config"
	"github.com/eyev0/timetracker/internal/db"
	"github.com/eyev0/timetracker/internal/utils"
)

func DeserializeUser(ctx *gin.Context) {
	var token string
	cookie, err := ctx.Cookie("token")

	authorizationHeader := ctx.Request.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		token = fields[1]
	} else if err == nil {
		token = cookie
	}

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
		return
	}

	config, _ := config.LoadConfig(".")
	sub, err := utils.ValidateJwt(token, config.JWTTokenSecret)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := db.GetUserById(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	ctx.Set("currentUser", user)
	ctx.Next()
}
