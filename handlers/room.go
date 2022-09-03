package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/touch-some-grass-bro/letterly-api/utils"
)

func CreateRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp map[string]interface{} = make(map[string]interface{})
    // TODO: Add to database?
		hostSessionID := r.Header.Get("sessionID")
		log.Println("Session ID", hostSessionID)
		room, err := utils.CreateHopChannel()
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
    var state map[string]string = make(map[string]string)
    if err := json.Unmarshal([]byte(stateValue), &state); err != nil {
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

