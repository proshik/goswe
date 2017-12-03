package main

import (
	"io/ioutil"
	"os"
	"bufio"
	"strings"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	yToken := os.Getenv("Y_TOKEN")
	if yToken == "" {
		panic("Y_TOKEN is required variable")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		panic("DB_PATH is required variable")
	}

	file, err := os.Open("result.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)

	var rawWords = make([]RawWord, 0)
	err = decoder.Decode(&rawWords)
	if err != nil {
		panic(err)
	}

	yandex := NewYandex(yToken)
	db := NewDB(dbPath)

	for i := 0; i < 3; i++ {
		tr, err := yandex.translate(rawWords[i].Text)
		if err != nil {
			log.Fatalf("Error on iteration=%d, text=%s, with error=%v", i, rawWords[i].Text, err)
		}

		word, err := db.AddWOrd(Word{rawWords[i].Text, rawWords[i].Category, rawWords[i].Subcategory, tr.Def})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Success translate and save word=%s\n", word.Text)
	}

	//for _, rw := range rawWords {
	//
	//	yandex.
	//
	//}
}

func readRawWords() {
	fInfos, err := ioutil.ReadDir("words")
	if err != nil {
		panic(err)
	}

	var count int
	rawWords := make([]RawWord, 0)
	for _, fInfo := range fInfos {

		words := func(fInfo os.FileInfo) []RawWord {
			result := make([]RawWord, 0)

			file, err := os.Open("words" + "/" + fInfo.Name())
			if err != nil {
				panic(err)
			}

			defer file.Close()

			s := bufio.NewScanner(file)

			for s.Scan() {
				elem := strings.Split(s.Text(), " - ")

				subcategoryTitle := strings.TrimSuffix(fInfo.Name(), ".txt")

				result = append(result, RawWord{elem[0], "Basic English words", subcategoryTitle})
			}

			return result
		}(fInfo)

		for _, w := range words {
			rawWords = append(rawWords, w)
			count++
		}
	}

	b, err := json.MarshalIndent(&rawWords, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}

	err = ioutil.WriteFile("result.json", b, 06444)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result connt %d\n", count)
}
