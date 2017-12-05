package main

import (
	"encoding/json"
	"net/http"
)

type YandexDict struct {
	Token string
}

func NewYandex(token string) *YandexDict {
	return &YandexDict{token}
}

func (yDict *YandexDict) translate(text string) (*Translate, error) {

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

//func main() {
//	yandex := YandexDict{""}
//
//	result, err := yandex.translate("success")
//	if err != nil {
//		panic(err)
//	}
//
//	b, err := json.MarshalIndent(&result, "", "\t")
//	if err != nil {
//		fmt.Println("error:", err)
//	}
//
//	fmt.Printf("%s\n", b)
//}
