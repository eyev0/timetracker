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

func updateCurrentEntry(user *model.User, payload *model.UpdateEntryInput, now time.Time) (entry *model.Entry, code int, status string, err error) {
	entry, tx, err := db.UpdateEntry(user, payload)
	status = "fail"
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = errors.Join(rollbackErr, err)
			log.Logger.Errorf("%v", err)
		}
		if strings.Contains(err.Error(), "no rows in result set") {
			code = http.StatusNotFound
			status = "not found"
		} else {
			code = http.StatusInternalServerError
		}
		return
	}

	entry.CalcElapsed(now)

	if !payload.NoSync {
		err = calendar.PostEvent(user, entry)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = errors.Join(rollbackErr, err)
				log.Logger.Errorf("%v", err)
			}
			code = http.StatusInternalServerError
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	code = http.StatusOK
	status = "success"
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
		var code int
		var status string
		currentEntry, code, status, err = updateCurrentEntry(user, &model.UpdateEntryInput{}, now)
		if err != nil {
			ctx.JSON(code, gin.H{"status": status, "message": err.Error()})
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
			payload = &model.UpdateEntryInput{}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
			return
		}
	}

	entry, code, status, err := updateCurrentEntry(user, payload, now)
	if err != nil {
		ctx.JSON(code, gin.H{"status": status, "message": err.Error()})
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
			ctx.JSON(http.StatusNotFound, gin.H{"status": "not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
	}

	entry.CalcElapsed(now)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"entry": entry}})
}
