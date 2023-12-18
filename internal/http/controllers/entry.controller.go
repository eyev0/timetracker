package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/eyev0/timetracker/internal/db"
	"github.com/eyev0/timetracker/internal/google_api/calendar"
	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
)

func CreateEntry(ctx *gin.Context) {
	user := ctx.MustGet("currentUser").(*model.User)

	var payload *model.CreateEntryInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	entryId, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	entry := &model.Entry{Id: entryId, UserId: user.Id, Note: payload.Note}

	if err := db.CreateEntry(entry, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"entry": entry}})
}

func UpdateEntry(ctx *gin.Context) {
	user := ctx.MustGet("currentUser").(*model.User)

	var payload *model.UpdateEntryInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		if err.Error() == "EOF" {
			// no payload = ok
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
	}

	entry, tx, err := db.UpdateEntry(user, payload)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Logger.Error("%+v", rollbackErr)
			message := fmt.Sprintf("Error while rolling back transaction: %s, cause: %s", rollbackErr, err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": message})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	err = calendar.PostEvent(user, entry)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Logger.Error("%+v", rollbackErr)
			message := fmt.Sprintf("Error while rolling back transaction: %s, cause: %s", rollbackErr, err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": message})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"entry": entry}})
}

func GetEntry(ctx *gin.Context) {
	user := ctx.MustGet("currentUser").(*model.User)

	entry, err := db.GetCurrentUserEntry(user)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"entry": entry}})
}
