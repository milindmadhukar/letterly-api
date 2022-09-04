package utils

import "math/rand"

func GenerateRoomCode(length int) string {
  // Generate a random room code
  var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
  b := make([]rune, length)
  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)

}
