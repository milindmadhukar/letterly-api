package utils

import (
	"math/rand"
	"time"

	"github.com/touch-some-grass-bro/letterly-api/models"
)

func GetPlayingPlayers(players []models.Player) []models.Player {
  playing := make([]models.Player, 0)
  for _, player := range players {
    if player.IsPlaying {
      playing = append(playing, player)
    }
  }
  return playing
}

func GetCurrentPlayer(players []models.Player) (models.Player, []models.Player) {
  rand.Seed(time.Now().UnixNano())
  idx := rand.Intn(len(players))
  // Remove element at idx
  remaining := removeIndex(players, idx)

  return players[idx], remaining

}

func removeIndex(s []models.Player, index int) []models.Player {
  return append(s[:index], s[index+1:]...)
}
