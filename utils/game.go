package utils

import (
	"math/rand"
	"time"

	"github.com/touch-some-grass-bro/letterly-api/models"
)

func GetPlayingPlayers(players []models.Player) []string {
  playing := make([]string, 0)
  for _, player := range players {
    if player.IsPlaying {
      playing = append(playing, player.SessionID)
    }
  }
  return playing
}

func GetCurrentPlayer(players []string) (string, []string) {
  rand.Seed(time.Now().UnixNano())
  idx := rand.Intn(len(players))
  // Remove element at idx
  currPlayer := players[idx]
  remaining := removeIndex(players, idx)
  remainingSessions := make([]string, 0)
  for _, player := range remaining {
   remainingSessions = append(remainingSessions, player)
  }

  return currPlayer, remainingSessions

}

func FindPlayer(sessionID string, players []models.Player) int {
  for idx, player := range players {
    if player.SessionID == sessionID {
      return idx
    }
  }
  return -1
}

func IsUserInGame(sessionID string, players []models.Player) bool {
  for _, player := range players {
    if player.SessionID == sessionID {
      return true
    }
  }
  return false
}

func removeIndex(s []string, index int) []string {
  return append(s[:index], s[index+1:]...)
}
