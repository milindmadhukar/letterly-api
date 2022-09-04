package handlers

import (
	"net/http"
	"strconv"
	"time"

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

    state, err := utils.GetChannelState(channelID)
    if err != nil {
      resp["error"] = err.Error()
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }

    if state.Game == "started" {
      resp["error"] = "Game already started."
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }

    for idx := range state.Players {
      state.Players[idx].IsPlaying = true
      state.Players[idx].Score = 0
    }

    state.Game = "started"
    state.Round = 1
    state.RoundsPerStage = roundsInt
    state.Stage = 1
    state.StartTime = time.Now()

    state.YetToPlay = utils.GetPlayingPlayers(state.Players)
    // Pick a random player from Players
    state.CurrentPlayer, state.YetToPlay = utils.GetCurrentPlayer(state.YetToPlay)

		if err := utils.UpdateHopChannel(channelID, state); err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}
		resp["success"] = "Game started successfully"
		utils.JSON(w, http.StatusOK, resp)
	}
}
