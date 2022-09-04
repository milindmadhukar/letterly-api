package handlers

import (
	"net/http"
	"time"

	"github.com/touch-some-grass-bro/letterly-api/utils"
)

func StartGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp map[string]interface{} = make(map[string]interface{})
		hostSessionID := r.Header.Get("sessionID")
    channelID := r.URL.Query().Get("channelID")
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
    newState := map[string]interface{}{
      "game":"started",
      "startTime" : time.Now(),
    }
		err := utils.UpdateHopChannel(channelID, newState)
		if err != nil {
      resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}
    resp["success"] = "Game started successfully"
		utils.JSON(w, http.StatusOK, resp)
	}
}

