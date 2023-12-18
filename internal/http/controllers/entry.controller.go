package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/eyev0/timetracker/internal/db"
	"github.com/eyev0/timetracker/internal/google_api/calendar"
	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
)

func stopCurrent(user *model.User, now time.Time) (entry *model.Entry, err error) {
	entry, tx, err := db.UpdateEntry(user, nil)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = errors.Join(rollbackErr, err)
		}
		log.Logger.Errorf("%v", err)
		return
	}

	entry.CalcElapsed(now)

	err = calendar.PostEvent(user, entry)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = errors.Join(rollbackErr, err)
		}
		log.Logger.Errorf("%v", err)
		return
	}

	err = tx.Commit()
	return
}

func CreateEntry(ctx *gin.Context) {
	now := time.Now()
	user := ctx.MustGet("currentUser").(*model.User)

	var payload *model.CreateEntryInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// if current entry exists, update it and start new entry
	currentEntry, err := db.GetCurrentUserEntry(user)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			// ok, create new entry
			currentEntry = nil
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
	} else {
		currentEntry, err = stopCurrent(user, now)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
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

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"new_entry": entry, "entry": currentEntry}})
}

func UpdateEntry(ctx *gin.Context) {
	now := time.Now()
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
			err = errors.Join(rollbackErr, err)
		}
		log.Logger.Error("%v", rollbackErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	entry.CalcElapsed(now)

	err = calendar.PostEvent(user, entry)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = errors.Join(rollbackErr, err)
		}
		log.Logger.Error("%v", rollbackErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
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
	now := time.Now()
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

	entry.CalcElapsed(now)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"entry": entry}})
}
