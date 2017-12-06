package main

import (
	"encoding/json"
	"net/http"
)

type YDict struct {
	Token string
}

func NewYDict(token string) *YDict {
	return &YDict{token}
}

func (yDict *YDict) translate(text string) (*Translate, error) {

	resp, err := http.Get("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?lang=en-ru&key=" + yDict.Token + "&text=" + text)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	var tr Translate
	err = d.Decode(&tr)
	if err != nil {
		return nil, err
	}

	return &tr, nil
}
