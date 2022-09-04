package handlers

import (
	"net/http"
	"strconv"
	"time"

	db "github.com/touch-some-grass-bro/letterly-api/db/sqlc"
	"github.com/touch-some-grass-bro/letterly-api/utils"
)

func StartGame(queries *db.Queries) http.HandlerFunc {
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

		randomWord, err := queries.GetRandomWord(r.Context(), "6")
		if err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}

		state.Stage1Word = randomWord

		state.PlayerStartTime = time.Now()
		state.PlayerEndTime = state.StartTime.Add(time.Second * 10)

		if err := utils.UpdateHopChannel(channelID, state); err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}
		resp["success"] = "Game started successfully"
		utils.JSON(w, http.StatusOK, resp)
	}
}

func AnswerQuestion(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp map[string]interface{} = make(map[string]interface{})
		sessionID := r.Header.Get("sessionID")
		channelID := r.URL.Query().Get("channelID")
		word := r.URL.Query().Get("word")
		state, err := utils.GetChannelState(channelID)
		if err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}
		if sessionID == "" {
			resp["error"] = "No sessionID provided."
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}

		if sessionID != state.CurrentPlayer {
			resp["error"] = "It's not your turn."
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}

    if word == "" {
      resp["error"] = "No word provided."
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }

		timeRemaining := state.PlayerEndTime.Sub(time.Now())
		if timeRemaining < 0 {
			resp["error"] = "Time's up."
			// Pick a random player from Players
			state.CurrentPlayer, state.YetToPlay = utils.GetCurrentPlayer(state.YetToPlay)
			err = utils.SendMessageToHopChannel("PLAYER_TIME_UP", channelID, map[string]interface{}{
				"player": state.CurrentPlayer,
			})
			if err != nil {
				resp["error"] = err.Error()
				utils.JSON(w, http.StatusBadRequest, resp)
				return
			}
			utils.JSON(w, http.StatusOK, resp)
			return
		}

		err = utils.SendMessageToHopChannel("PLAYER_ANSWER", channelID, map[string]interface{}{
			"player": state.CurrentPlayer,
			"word":   word,
		})
		if err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}

		if state.Stage == 1 {
			if utils.IsLastLetterMatching(state.Stage1Word, word) {
				player_idx := utils.FindPlayer(state.CurrentPlayer, state.Players)
				state.Players[player_idx].Score += int(timeRemaining.Seconds())*10 + len(word)
				state.Stage1Word = word
				resp["status"] = "correct"
			}
		}

    if state.Stage == 2 {

    }

    if state.Stage == 3 {

    }

		if len(state.YetToPlay) == 0 {
			if state.Round >= state.RoundsPerStage {
				if state.Stage >= 3 {
					state.Game = "finished"
					state.Stage = 0
					state.Round = 0
				} else {
					state.Stage++
          resp["stage"] = state.Stage
					state.Round = 0
				}
			}
			state.Round++
      resp["round"] = state.Round
			state.YetToPlay = utils.GetPlayingPlayers(state.Players)
		} else {
			// Pick a random player from Players
			state.CurrentPlayer, state.YetToPlay = utils.GetCurrentPlayer(state.YetToPlay)
		}

		state.PlayerStartTime = time.Now().Add(time.Second * 5)
		state.PlayerEndTime = state.PlayerStartTime.Add(time.Second * 10)

		if err := utils.UpdateHopChannel(channelID, state); err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}
		resp["success"] = "You answered."
		resp["status"] = "incorrect"
		utils.JSON(w, http.StatusOK, resp)
	}
}

func GetNextPlayer(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp map[string]interface{} = make(map[string]interface{})
		channelID := r.URL.Query().Get("channelID")

		state, err := utils.GetChannelState(channelID)
		if err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}

		// Pick a random player from Players

		if len(state.YetToPlay) == 0 {
			if state.Round >= state.RoundsPerStage {
				if state.Stage >= 3 {
					state.Game = "finished"
					state.Stage = 0
					state.Round = 0
				} else {
					state.Stage++
					state.Round = 0
				}
			}
			state.Round++
			if state.Stage == 1 {

				randomWord, err := queries.GetRandomWord(r.Context(), "6")
				if err != nil {
					resp["error"] = err.Error()
					utils.JSON(w, http.StatusBadRequest, resp)
					return
				}
				state.Stage1Word = randomWord
			}
		} else {
			state.CurrentPlayer, state.YetToPlay = utils.GetCurrentPlayer(state.YetToPlay)
		}

		state.PlayerStartTime = time.Now().Add(time.Second * 5)
		state.PlayerEndTime = state.PlayerStartTime.Add(time.Second * 10)

		if err := utils.UpdateHopChannel(channelID, state); err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}

		resp["success"] = "Next player's turn."
		utils.JSON(w, http.StatusOK, resp)
	}
}
