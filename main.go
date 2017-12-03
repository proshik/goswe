package main

import (
	"io/ioutil"
	"os"
	"bufio"
	"strings"
	"encoding/json"
	"fmt"
	"log"
	"github.com/manifoldco/promptui"
	"errors"
	"strconv"
	"github.com/fatih/color"
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

	yandexDict := NewYandex(yToken)
	dbConnect := NewDB(dbPath)

	for {
		prompt := promptui.Select{
			Label: "Select command",
			Items: []string{"translate", "learn", "managment"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case "translate":
			translate(yandexDict, dbConnect)
		case "learn":
			learn(dbConnect)
		case "managment":
			managment(dbConnect)
		}
	}
}

func translate(yandex *YandexDict, db *DBConnect) {
	for {
		validate := func(input string) error {
			if len(input) < 2 || len(input) > 25 {
				return errors.New("From 2 to 25 symbols\n")
			}
			_, err := strconv.ParseFloat(input, 64)
			if err == nil {
				return errors.New("Not numbers\n")
			}
			return nil
		}

		prompt := promptui.Prompt{
			Label:    "Word",
			Validate: validate,
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Error on read word, %v\n", err)
			return
		}

		word, err := db.GetEnv(result)
		if err != nil {
			fmt.Printf("Error on translate word=%s\n", result)
		}

		if word != nil {
			printWord(word)
			continue
		}

		fmt.Printf("Will be translate throw YandexDict\n")

		tr, err := yandex.translate(result)
		if err != nil {
			fmt.Printf("Error on translate word=%s\n", result)
			continue
			//todo тут можно сделать зарпос на попытку повторного перевода
		}

		word, err = db.AddWOrd(Word{result, "Default", "Default", tr.Def})
		if err != nil {
			fmt.Printf("Error on save word=%s in db\n", result)
			continue
		}

		if word.Translate == nil || len(word.Translate) == 0 {
			fmt.Printf("Translate not found")
			continue
		}

		printWord(word)
	}

}

func printWord(word *Word) {
	color.Red("%s\n", word.Translate[0].Tr[0].Text)
	fmt.Printf("\n")

	for _, w := range word.Translate {
		fmt.Printf("%s [%s] %s\n", w.Text, w.Ts, w.Pos)
		for i, t := range w.Tr {
			fmt.Printf("%d  %s %s", i+1, t.Text, t.Gen)
			if len(t.Syn) > 0 {
				for _, s := range t.Syn {
					fmt.Printf(", %s %s", s.Text, t.Gen)
				}
				fmt.Printf("\n")
				if len(t.Mean) > 0 {
					fmt.Printf("(")
					for im, s := range t.Mean {
						if im+1 != len(t.Mean) {
							fmt.Printf("%s, ", s.Text)
						} else {
							fmt.Printf("%s", s.Text)
						}
					}
					fmt.Printf(")\n")
				}
			} else {
				fmt.Printf("\n")
			}

		}
		fmt.Printf("\n")
	}

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func learn(db *DBConnect) {
	prompt := promptui.Select{
		Label: "Select command",
		Items: []string{"words", "rules"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("Selected %s", result)

}

func managment(db *DBConnect) {

}

func fillBasicEnglishWords(yandex *YandexDict, db *DBConnect) {
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
		tr, err := yandex.translate(rw.Text)
		if err != nil {
			log.Fatalf("Error on word=%s, with error=%v", rw.Text, err)
		}

		word, err := db.AddWOrd(Word{rw.Text, rw.Category, rw.Subcategory, tr.Def})
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
