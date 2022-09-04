package utils

import (
	"math/rand"
	"time"
)

func GenerateRoomCode(length int) string {
  rand.Seed(time.Now().UnixNano())
  var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
  b := make([]rune, length)
  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)

}
