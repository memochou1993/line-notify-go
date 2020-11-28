package app

import (
	"encoding/json"
	"log"
)

type payload struct {
	Status      string `json:"Name"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

func Parse(raw []byte) *payload {
	payload := &payload{}

	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Println(err.Error())
	}

	return payload
}
