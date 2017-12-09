package view

import (
	"bytes"
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	"log"
	"strings"
	"github.com/proshik/goswe/yandex"
	"github.com/proshik/goswe/model"
	_"bufio"
)

const (
	ALL_VIEWS = ""

	TEXT_VIEW      = "text"
	TRANSLATE_VIEW = "translate"
	HISTORY_VIEW   = "history"
)

var VIEW_TITLES = map[string]string{
	TEXT_VIEW:      "Text",
	TRANSLATE_VIEW: "Dictionary",
	HISTORY_VIEW:   "History",
}

var VIEWS = []string{
	TEXT_VIEW,
	TRANSLATE_VIEW,
	HISTORY_VIEW,
}

var activeIndex = 0

var optionsDetectLang = whatlanggo.Options{
	Whitelist: map[whatlanggo.Lang]bool{
		whatlanggo.Eng: true,
		whatlanggo.Rus: true,
	},
}

type TranslateLangOpt struct {
	source      string
	destination string
}

type UI struct {
	YDict *yandex.YDictionary
}

func NewUI(yd *yandex.YDictionary) *UI {
	return &UI{yd}
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

	if err := g.SetKeybinding(ALL_VIEWS, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(TEXT_VIEW, gocui.KeyCtrlY, gocui.ModNone, ui.handleText); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(TEXT_VIEW, gocui.KeyEsc, gocui.ModNone, cleanView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(ALL_VIEWS, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (ui *UI) handleText(g *gocui.Gui, v *gocui.View) error {
	//extract text from TEXT view
	textFromTextView := getViewValue(g, TEXT_VIEW)
	if textFromTextView == "" {
		return nil
	}
	//detect a language in text
	info := whatlanggo.DetectLangWithOptions(textFromTextView, optionsDetectLang)
	langOpt, err := createTranslateLangOpt(info)
	if err != nil {
		return err
	}
	// update TRANSLATE view. Exactly search translate word in db or translate with yandex dictionary
	go g.Update(func(g *gocui.Gui) error {
		translateView, err := g.View(TRANSLATE_VIEW)
		if err != nil {
			return err
		}

		//_ := strings.NewReader(textFromTextView)
		////scaner := bufio.NewScanner(reader)
		//
		//scanner := bufio.NewScanner(strings.NewReader(textFromTextView))
		//scanner.Split(bufio.ScanLines)
		//for scanner.Scan() {
		//	fmt.Println(scanner.Text())
		//}


		mayBeWord := maybeWord(textFromTextView)

		fmt.Printf("value: %s, len=%d, maybe=%s\n", textFromTextView, len(textFromTextView), mayBeWord)

		textFromTextView = strings.Replace(textFromTextView, "\n", " ", -1)

		fmt.Printf("value: %s, len=%d, maybe=%s\n", textFromTextView, len(textFromTextView), mayBeWord)

		fmt.Printf("contains=%s", strings.Contains(textFromTextView, "\n"))

		//fmt.Fprintln(translateView, "...translated...")
		word, err := ui.translate(textFromTextView, langOpt.source, langOpt.destination)
		if err != nil {
			fmt.Fprintln(translateView, "...error on translate text...")
		}

		translateView.Clear()

		if !word.IsEmpty() {
			fmt.Fprintln(translateView, printWord(word))
		} else {
			fmt.Fprintln(translateView, "...translate not found...")
		}

		return nil
	})
	//push text to HISTORY view
	go func() {
		g.Update(func(g *gocui.Gui) error {
			historyView, err := g.View(HISTORY_VIEW)
			if err != nil {
				return err
			}

			fmt.Fprintln(historyView, textFromTextView)
			return nil
		})
	}()
	return nil
}
func maybeWord(text string) bool {
	if len(text) < 27 {
		return true
	}
	return false
}

func (ui *UI) translate(text string, langFrom string, langTo string) (*model.Word, error) {
	//word, err := ui.DBConnect.GetWords(text)
	//if err != nil {
	//	return nil, err
	//}

	//if !word.isEmpty() {
	//	return word, nil
	//}

	tr, err := ui.YDict.Translate(text, langFrom, langTo)
	if err != nil {
		log.Printf("Error on translate word, err=%v", err)
		return nil, err
	}
	//
	//word, err = ui.DBConnect.AddWord(Word{text, "Default", "Default", tr.Def})
	//if err != nil {
	//	log.Printf("Error on add word in DB, err=%v", err)
	//	return nil, err
	//}

	return &model.Word{text, "Default", "Default", tr.Def}, nil
}

func cleanView(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		return nil
	})
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (activeIndex + 1) % len(VIEW_TITLES)

	if _, err := g.SetCurrentView(VIEWS[nextIndex]); err != nil {
		return err
	}

	//todo check this place when will be need to enable cursor on view
	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	activeIndex = nextIndex
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func createTranslateLangOpt(lang whatlanggo.Lang) (*TranslateLangOpt, error) {
	switch lang {
	case whatlanggo.Rus:
		return &TranslateLangOpt{"ru", "en"}, nil
	case whatlanggo.Eng:
		return &TranslateLangOpt{"en", "ru"}, nil
	}
	return nil, errors.New(fmt.Sprintf("Unrecognize unexpected language=%s", whatlanggo.LangToString(lang)))
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(TEXT_VIEW, 0, 0, maxX/2-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = VIEW_TITLES[TEXT_VIEW]
		v.Editable = true
		v.Wrap = true

		if _, err = g.SetCurrentView(TEXT_VIEW); err != nil {
			return err
		}
	}

	if v, err := g.SetView(TRANSLATE_VIEW, maxX/2, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = VIEW_TITLES[TRANSLATE_VIEW]
		v.Wrap = true
		v.Autoscroll = true
	}
	if v, err := g.SetView(HISTORY_VIEW, 0, maxY/2, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = VIEW_TITLES[HISTORY_VIEW]
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil
}

func getViewValue(g *gocui.Gui, name string) string {
	v, err := g.View(name)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(v.Buffer())
}

func printWord(word *model.Word) string {
	var buf bytes.Buffer

	buf.WriteString(word.Translate[0].Tr[0].Text + "\n")
	buf.WriteString("\n")

	for _, w := range word.Translate {
		buf.WriteString(fmt.Sprintf("%s [%s] %s\n", w.Text, w.Ts, w.Pos))
		for i, t := range w.Tr {
			buf.WriteString(fmt.Sprintf("%d  %s %s", i+1, t.Text, t.Gen))
			if len(t.Syn) > 0 {
				for _, s := range t.Syn {
					buf.WriteString(fmt.Sprintf(", %s %s", s.Text, t.Gen))
				}
				fmt.Printf("\n")
				if len(t.Mean) > 0 {
					buf.WriteString(fmt.Sprintf("("))
					for im, s := range t.Mean {
						if im+1 != len(t.Mean) {
							buf.WriteString(fmt.Sprintf("%s, ", s.Text))
						} else {
							buf.WriteString(fmt.Sprintf("%s", s.Text))
						}
					}
					buf.WriteString(fmt.Sprintf(")\n"))
				}
			} else {
				buf.WriteString(fmt.Sprintf("\n"))
			}

		}
		buf.WriteString(fmt.Sprintf("\n"))
	}
	return buf.String()
}
