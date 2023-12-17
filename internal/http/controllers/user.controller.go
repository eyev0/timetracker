package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eyev0/timetracker/internal/db"
	"github.com/eyev0/timetracker/internal/model"
)

func GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*model.User)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": model.FilteredResponse(currentUser)}})
}

func PatchSettings(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*model.User)

	var payload *model.SettingsInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := db.UpdateUserSettings(currentUser, payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
