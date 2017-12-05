package main

import (
	"github.com/jroimartin/gocui"
	"strings"
	"fmt"
	"log"
	"encoding/json"
)

var (
	viewArr = []string{"text", "translate", "history"}
	active  = 0
)

type UI struct {
	YandexDict *YandexDict
	DBConnect  *DBConnect
}

func NewUI(yd *YandexDict, dbc *DBConnect) *UI {
	return &UI{yd, dbc}
}

func (ui *UI) Run() {
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

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, ui.nextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("text", gocui.KeyEnter, gocui.ModNone, ui.handleText); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("text", gocui.KeyEsc, gocui.ModNone, ui.cleanView); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (ui *UI) cleanView(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		return nil
	})
	return nil
}

func (ui *UI) nextView(g *gocui.Gui, v *gocui.View) error {
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

func (ui *UI) handleText(g *gocui.Gui, v *gocui.View) error {

	g.Update(func(g *gocui.Gui) error {

		translateView, err := g.View("translate")
		if err != nil {
			return err
		}

		word, err := translate(ui.YandexDict, ui.DBConnect, getViewValue(g, "text"))
		_, err = json.MarshalIndent(&word, "", "\t")
		if err != nil {
			fmt.Println("error:", err)
		}

		translateView.Clear()

		//valueFromTextView := getViewValue(g, "text")
		fmt.Fprintln(translateView, word.Translate[0].Tr[0].Text)
		//cleanView(g, text)
		return nil
	})
	return nil
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

func (ui *UI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
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

func translate(yandex *YandexDict, db *DBConnect, text string) (*Word, error) {
	//for {
	//	validate := func(input string) error {
	//		if len(input) < 2 || len(input) > 25 {
	//			return errors.New("From 2 to 25 symbols\n")
	//		}
	//		_, err := strconv.ParseFloat(input, 64)
	//		if err == nil {
	//			return errors.New("Not numbers\n")
	//		}
	//		return nil
	//	}

	//validate(word)
	//
	//prompt := promptui.Prompt{
	//	Label:    "Word",
	//	Validate: validate,
	//}
	//
	//result, err := prompt.Run()
	//if err != nil {
	//	fmt.Printf("Error on read word, %v\n", err)
	//	return
	//}

	word, err := db.GetWords(text)
	if err != nil {
		return nil, err
		//fmt.Printf("Error on translate word=%s\n", text)
	}

	if word != nil && word.Translate != nil && len(word.Translate) > 0 {
		//printWord(word)
		return word, nil
		//continue
	}

	fmt.Printf("Will be translate throw YandexDict\n")

	tr, err := yandex.translate(text)
	if err != nil {
		fmt.Printf("Error on translate word=%s\n", text)
		return nil, err
		//continue
		//todo тут можно сделать зарпос на попытку повторного перевода
	}

	word, err = db.AddWord(Word{text, "Default", "Default", tr.Def})
	if err != nil {
		fmt.Printf("Error on save word=%s in db\n", text)
		return nil, err
		//continue
	}

	if word == nil || word.Translate == nil || len(word.Translate) == 0 {
		fmt.Printf("Translate not found\n")
		return nil, nil
		//continue
	}

	return word, nil
	//printWord(word)
	//}
}
