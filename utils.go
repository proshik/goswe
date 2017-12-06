package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func fillBasicEnglishWords(yandex *YDict, db *DBConnect) {
	file, err := os.Open("result.json")
	if err != nil {
		panic(err)
	}

	var rawWords = make([]RawWord, 0)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&rawWords)
	if err != nil {
		panic(err)
	}

	for _, rw := range rawWords {
		tr, err := yandex.translate(rw.Text, "en", "ru")
		if err != nil {
			log.Fatalf("Error on word=%s, with error=%v", rw.Text, err)
		}

		word, err := db.AddWord(Word{rw.Text, rw.Category, rw.Subcategory, tr.Def})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Success translate and save word=%s\n", word.Text)
	}
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
