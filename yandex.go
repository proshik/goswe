package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type YDict struct {
	Token string
}

func NewYDict(token string) *YDict {
	return &YDict{token}
}

func (yDict *YDict) translate(text string, langFrom string, langTo string) (*Translate, error) {

	url := fmt.Sprintf("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?"+
		"lang=%s-%s&key=%s&text=%s", langFrom, langTo, yDict.Token, text)

	resp, err := http.Get(url)
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
