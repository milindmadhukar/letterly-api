package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/touch-some-grass-bro/letterly-api/models"
	"github.com/touch-some-grass-bro/letterly-api/utils"
)

func StartGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp map[string]interface{} = make(map[string]interface{})
		hostSessionID := r.Header.Get("sessionID")
		channelID := r.URL.Query().Get("channelID")
		roundsPerStage := r.URL.Query().Get("roundsPerStage")
		if hostSessionID == "" {
			resp["error"] = "No sessionID provided."
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}
		if channelID == "" {
			resp["error"] = "No channelID provided."
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}
		roundsInt, err := strconv.Atoi(roundsPerStage)
		if err != nil {
			resp["error"] = "Invalid roundsPerStage provided."
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}

    newState := &models.ChannelState{
    	Game:           "started",
    	Round:          1,
    	RoundsPerStage: roundsInt,
    	Stage:          1,
    	StartTime:      time.Now(),
    }


		if err := utils.UpdateHopChannel(channelID, newState); err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}
		resp["success"] = "Game started successfully"
		utils.JSON(w, http.StatusOK, resp)
	}
}
