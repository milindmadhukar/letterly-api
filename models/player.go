package models

type Player struct {
  SessionID string `json:"sessionID"`
  UserName string `json:"userName"`
  Score int `json:"score"`
  IsPlaying bool `json:"isPlaying"`
}
