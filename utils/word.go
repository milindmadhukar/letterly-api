package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/touch-some-grass-bro/letterly-api/models"
)

func IsLastLetterMatching(previousWord, inputWord string) bool {
  if inputWord[0] == previousWord[len(previousWord) - 1] {
    return true
  }
  return false
}

const dictionaryAPIURL = "https://api.dictionaryapi.dev/api/v2/entries/en/"

func GetMeaning(word string) (*models.MeaningResponse, error) {
  resp, err := http.Get(dictionaryAPIURL + word)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  var data models.DictionaryAPIResponse
  err = json.NewDecoder(resp.Body).Decode(&data)
  if err != nil {
    return nil, err
  }

  if len(data) == 0 {
    return nil, errors.New("No meaning found.")
  }
  meaning := models.MeaningResponse{
  	Word:       word,
  	Definition: data[0].Meanings[0].Definitions[0].Definition,
  	Synonyms:   data[0].Meanings[0].Synonyms,
  	Antonyms:   data[0].Meanings[0].Antonyms,
  }

  return &meaning , nil
}

