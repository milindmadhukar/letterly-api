package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/touch-some-grass-bro/letterly-api/models"
	"github.com/touch-some-grass-bro/letterly-api/utils"
)

func CreateRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp map[string]interface{} = make(map[string]interface{})
		hostSessionID := r.Header.Get("sessionID")
    userName := r.URL.Query().Get("userName")
    if hostSessionID == "" {
      resp["error"] = "No sessionID provided."
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }
    if userName == "" {
      userName = "Anonymous"
    }
		room, err := utils.CreateHopChannel(hostSessionID, userName)
		if err != nil {
      resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
		}
		utils.JSON(w, http.StatusOK, room)
	}
}

func GetRoom() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
		var resp map[string]interface{} = make(map[string]interface{})
    roomID := chi.URLParam(r, "roomID")
    room, err := utils.GetHopChannel(roomID)
    if err != nil {
      resp["error"] = err.Error()
			utils.JSON(w, http.StatusBadRequest, resp)
			return
    }

    utils.JSON(w, http.StatusOK, room)
  }
}

func DeleteRoom() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var resp map[string]interface{} = make(map[string]interface{})
    roomID := chi.URLParam(r, "roomID")
    err := utils.DeleteHopChannel(roomID)
    if err != nil {
      resp["error"] = err.Error()
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }

    resp["success"] = "Room deleted successfully."
    utils.JSON(w, http.StatusOK, resp)
  }
}

// NOTE: Pass state as json in header.
func UpdateRoom() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var resp map[string]interface{} = make(map[string]interface{})
    roomID := chi.URLParam(r, "roomID")
    stateValue := r.Header.Get("state")
    var state *models.ChannelState
    if err := json.Unmarshal([]byte(stateValue), state); err != nil {
      resp["error"] = err.Error()
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }
    err := utils.UpdateHopChannel(roomID, state)
    if err != nil {
      resp["error"] = err.Error()
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }

    utils.JSON(w, http.StatusOK, resp)
  }
}
// WARN: This is for testing only
func SendMessageToRoom() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var resp map[string]interface{} = make(map[string]interface{})
    roomID := chi.URLParam(r, "roomID")
    message := chi.URLParam(r, "message")
    var data map[string]interface{} = make(map[string]interface{})
    data["message"] = message
    err := utils.SendMessageToHopChannel("SEND_DA_MESSAGE", roomID, data)
    if err != nil {
      resp["error"] = err.Error()
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }
    resp["success"] = "Message sent successfully."
    resp["content"] = message
    utils.JSON(w, http.StatusOK, resp)
  }
}

func JoinRoom() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var resp map[string]interface{} = make(map[string]interface{})
    roomID := chi.URLParam(r, "roomID")
    sessionID := r.Header.Get("sessionID")
    userName := r.URL.Query().Get("userName")
    state, err := utils.GetChannelState(roomID)
    if err != nil {
      resp["error"] = err.Error()
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }
    state.PlayerCount += 1
    state.Players = append(state.Players, models.Player{
    	SessionID: sessionID,
    	UserName:  userName,
    	Score:     0,
    	IsPlaying: false,
    })

    if err := utils.UpdateHopChannel(roomID, state); err != nil {
      resp["error"] = err.Error()
      utils.JSON(w, http.StatusBadRequest, resp)
      return
    }

    resp["success"] = "Player successfully joined."
    utils.JSON(w, http.StatusOK, resp)
  }
}
