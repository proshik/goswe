package yandex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/proshik/goswe/model"
)

type YDictionary struct {
	token string
}

func NewYDictionary(token string) *YDictionary {
	return &YDictionary{token}
}

func (yDict *YDictionary) Translate(text string, langFrom string, langTo string) (model.TranslatedText, error) {

	url := fmt.Sprintf("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?"+
		"lang=%s-%s&key=%s&text=%s", langFrom, langTo, yDict.token, text)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	var tr model.Dictionary
	err = d.Decode(&tr)
	if err != nil {
		return nil, err
	}

	return &tr, nil
}
