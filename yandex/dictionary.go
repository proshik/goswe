package yandex

import (
	"encoding/json"
	"fmt"
	"github.com/proshik/gotrew/model"
	"net/http"
)

type YDictionary struct {
	token string
}

func NewYDictionary(token string) *YDictionary {
	return &YDictionary{token}
}

type HttpError struct {
	ErrorResponse
	StatusCode int
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HttpError) Error() string {
	switch e.Code {
	case 401, 402:
		return fmt.Sprintf("ERROR! API key is invalid!")
	case 403:
		return fmt.Sprintf("ERROR! Exceeded limit on number of requests!")
	case 413:
		return fmt.Sprintf("ERROR! Exceeded max size of text!")
	default:
		return fmt.Sprintf("Error with StatusCode=%d from translate service!", e.StatusCode)
	}
}

//Throw standart error or Error
func (yDict *YDictionary) Translate(text string, langFrom string, langTo string) (model.TranslatedText, error) {

	url := fmt.Sprintf("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?"+
		"lang=%s-%s&key=%s&text=%s", langFrom, langTo, yDict.token, text)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		d := json.NewDecoder(resp.Body)

		var erResp ErrorResponse
		err = d.Decode(&erResp)
		if err != nil {

			return nil, err
		}
		return nil, &HttpError{erResp, resp.StatusCode}
	}

	d := json.NewDecoder(resp.Body)

	var tr model.Dictionary
	err = d.Decode(&tr)
	if err != nil {

		return nil, err
	}

	return &tr, nil
}
