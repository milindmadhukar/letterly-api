package models

import "time"

type ChannelState struct {
	Game            string    `json:"game,omitempty"`
	Host            string    `json:"host,omitempty"`
	PlayerCount     int       `json:"playerCount,omitempty"`
	Players         []Player  `json:"players,omitempty"`
	Round           int       `json:"round,omitempty"`
	RoundsPerStage  int       `json:"roundsPerStage,omitempty"`
	Stage           int       `json:"stage,omitempty"`
	StartTime       time.Time `json:"startTime,omitempty"`
	CurrentPlayer   string    `json:"currentPlayer,omitempty"`
	YetToPlay       []string  `json:"yetToPlay,omitempty"`
	PlayerStartTime time.Time `json:"playerStartTime,omitempty"`
	PlayerEndTime   time.Time `json:"playerEndTime,omitempty"`
	Stage1Word      string    `json:"stage1Word,omitempty"`
}
