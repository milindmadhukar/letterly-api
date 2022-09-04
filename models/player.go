package models

type Player struct {
  UserName string `json:"userName"`
  Score int `json:"score"`
  IsPlaying bool `json:"isPlaying"`
}
