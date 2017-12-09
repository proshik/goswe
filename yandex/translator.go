package yandex

type YTranslator struct {
	Token string
}

func NewYTranslator(token string) *YTranslator {
	return &YTranslator{token}
}

func (yDict *YTranslator) translate(text string, langFrom string, langTo string) error {

	//url := fmt.Sprintf("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?"+
	//	"lang=%s-%s&key=%s&text=%s", langFrom, langTo, yDict.Token, text)
	//
	//resp, err := http.Get(url)
	//if err != nil {
	//	return nil, err
	//}
	//
	//defer resp.Body.Close()
	//
	//d := json.NewDecoder(resp.Body)
	//
	//var tr Translate
	//err = d.Decode(&tr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &tr, nil

	return nil
}
