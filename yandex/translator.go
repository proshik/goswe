package yandex

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"net/url"
	"github.com/proshik/goswe/model"
)

type YTranslator struct {
	Token string
}

func NewYTranslator(token string) *YTranslator {
	return &YTranslator{token}
}

func (y *YTranslator) Translate(text string, langFrom string, langTo string) (model.TranslatedText, error) {

	urlPath := fmt.Sprintf("https://translate.yandex.net/api/v1.5/tr.json/translate?"+
		"lang=%s-%s&key=%s", langFrom, langTo, y.Token)

	urlValues := url.Values{}
	urlValues.Add("text", text)

	client := &http.Client{}
	r, err := http.NewRequest("POST", urlPath, strings.NewReader(urlValues.Encode()))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	var tr model.Translate
	err = d.Decode(&tr)
	if err != nil {
		return nil, err
	}

	return &tr, nil
}