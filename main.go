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

	"github.com/jroimartin/gocui"
)

func main() {
	//yToken := os.Getenv("Y_TOKEN")
	//if yToken == "" {
	//	panic("Y_TOKEN is required variable")
	//}
	//yToken := "dict.1.1.20171201T214544Z.c0def3859d70a33d.88c6cf03a4e01eae8d732fce76205af4d35a7956"

	//dbPath := os.Getenv("DB_PATH")
	//if dbPath == "" {
	//	panic("DB_PATH is required variable")
	//}

	//dbPath := "database.db"

	//yandexDict := NewYandex(yToken)
	//dbConnect := NewDB(dbPath)

	//for {
	//	prompt := promptui.Select{
	//		Label: "Select command",
	//		Items: []string{"translate", "learn", "managment"},
	//	}
	//
	//	_, result, err := prompt.Run()
	//	if err != nil {
	//		fmt.Printf("Prompt failed %v\n", err)
	//		return
	//	}
	//
	//	switch result {
	//	case "translate":
	//		translate(yandexDict, dbConnect)
	//	case "learn":
	//		learn(dbConnect)
	//	case "managment":
	//		managment(dbConnect)
	//	}
	//}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorYellow
	g.InputEsc = true

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("text", gocui.KeyEnter, gocui.ModNone, handleText); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("text", gocui.KeyEsc, gocui.ModNone, cleanView); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
func cleanView(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		return nil
	})

	return nil
}

var (
	viewArr = []string{"text", "translate", "history"}
	active  = 0
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	out, err := g.View("translate")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func handleText(g *gocui.Gui, v *gocui.View) error {

	g.Update(func(g *gocui.Gui) error {
		//_, err := g.View("text")
		//if err != nil {
		//	return err
		//}

		translate, err := g.View("translate")
		if err != nil {
			return err
		}


		value := getViewValue(g, "text")
		fmt.Fprintln(translate, value)
		//cleanView(g, text)
		return nil
	})

	return nil
}

func getViewValue(g *gocui.Gui, name string) string {
	v, err := g.View(name)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(v.Buffer())
}

func setViewDefaults(v *gocui.View) {
	v.Frame = true
	v.Wrap = false
}

func setViewTextAndCursor(v *gocui.View, s string) {
	v.Clear()
	fmt.Fprint(v, s)
	v.SetCursor(len(s), 0)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("text", 0, 0, maxX/2-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Text"
		v.Editable = true
		v.Wrap = true
		//v.Highlight = true

		if _, err = setCurrentViewOnTop(g, "text"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("translate", maxX/2, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Translate"
		v.Wrap = true
		v.Autoscroll = true
	}
	if v, err := g.SetView("history", 0, maxY/2, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "History"
		v.Wrap = true
		v.Autoscroll = true
		fmt.Fprint(v, "Press TAB to change current view")
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

/*
 Follow need functions. Not remove!
 */

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

		word, err := db.GetWords(result)
		if err != nil {
			fmt.Printf("Error on translate word=%s\n", result)
		}

		if word != nil && word.Translate != nil && len(word.Translate) > 0 {
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

		word, err = db.AddWord(Word{result, "Default", "Default", tr.Def})
		if err != nil {
			fmt.Printf("Error on save word=%s in db\n", result)
			continue
		}

		if word == nil || word.Translate == nil || len(word.Translate) == 0 {
			fmt.Printf("Translate not found\n")
			continue
		}

		printWord(word)
	}

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

	switch result {
	case "words":

	case "rules":
	}

	fmt.Printf("Selected %s", result)

}

func managment(db *DBConnect) {
	fmt.Printf("Not implemented")
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
