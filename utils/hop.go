package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/touch-some-grass-bro/letterly-api/models"
)

const hopURL string = "https://api.hop.io/v1/"

type hopResponse struct {
	Success bool `json:"success,omitempty"`
	Data    struct {
		Channel hopChannel `json:"channel"`
	} `json:"data,omitempty"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type hopChannel struct {
	ID        string          `json:"id"`
	State     json.RawMessage `json:"state"`
	CreatedAt string          `json:"created_at"`
	Type      string          `json:"type"`
}

func ExecuteHopRequest(endpoint, reqMethod string, reqBody io.Reader, params map[string]string) (*hopResponse, error) {
	url := hopURL + endpoint
	req, err := http.NewRequest(reqMethod, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", models.Config.API.HopToken)
	req.Header.Add("Content-Type", "application/json")

	for key, value := range params {
		req.URL.Query().Add(key, value)
	}

	var resp *hopResponse

	hopResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer hopResp.Body.Close()
	// bodyBytes, err := io.ReadAll(hopResp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bodyString := string(bodyBytes)
	// log.Println("Response from hop", bodyString)

	if err := json.NewDecoder(hopResp.Body).Decode(&resp); err != nil {
		if err.Error() == "EOF" {
			return &hopResponse{}, nil
		}
		log.Println("Error decoding hop response:", err)
		return nil, err
	}

	return resp, nil
}

func CreateHopChannel(hostSessionID, username string) (*hopChannel, error) {
	// Create a map of players with session ID to username
	reqBody, err := json.Marshal(map[string]interface{}{
		"type": "unprotected",
		"state": map[string]interface{}{
			"game": "created",
			"host": hostSessionID,
			// Create a state players object as a map of sessionID to player name
			"players": map[string]string{
				hostSessionID: username,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := ExecuteHopRequest(
		"channels",
		"POST",
		bytes.NewBuffer(reqBody),
		map[string]string{
			"project": models.Config.API.HopProjectID,
		},
	)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Error.Message)
	}

	return &resp.Data.Channel, err
}

func GetHopChannel(roomID string) (*hopChannel, error) {
	resp, err := ExecuteHopRequest(
		"channels/"+roomID,
		"GET",
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Error.Message)
	}

	return &resp.Data.Channel, err
}

func DeleteHopChannel(roomID string) error {
	resp, err := ExecuteHopRequest(
		"channels/"+roomID,
		"DELETE",
		nil,
		nil,
	)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Error.Message)
	}

	return nil
}

func UpdateHopChannel(roomID string, state map[string]string) error {
	reqBody, err := json.Marshal(map[string]interface{}{
		"type":  "unprotected",
		"state": state,
	})
	if err != nil {
		return err
	}

	resp, err := ExecuteHopRequest(
		"channels/"+roomID,
		"PATCH",
		bytes.NewBuffer(reqBody),
		map[string]string{
			"project": models.Config.API.HopProjectID,
			"channel": roomID,
		},
	)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Error.Message)
	}

	return nil
}

func SendMessageToHopChannel(event string, roomID string, data map[string]interface{}) error {
	reqBody, err := json.Marshal(map[string]interface{}{
		"e": event,
		"d": data,
	})
	if err != nil {
		return err
	}

	_, err = ExecuteHopRequest(
		"channels/"+roomID+"/messages",
		"POST",
		bytes.NewBuffer(reqBody),
		map[string]string{
			"channel": roomID,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
