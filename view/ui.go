package view

import (
	"errors"
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/jroimartin/gocui"
	"github.com/proshik/goswe/model"
	"github.com/proshik/goswe/yandex"
	"log"
	"strings"
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

var OPTIONS_DETECT_LANG = whatlanggo.Options{
	Whitelist: map[whatlanggo.Lang]bool{
		whatlanggo.Eng: true,
		whatlanggo.Rus: true,
	},
}

var activeIndex = 0

var history = make([]string, 0)

type TranslateLangOpt struct {
	source      string
	destination string
}

type UI struct {
	YDict *yandex.YDictionary
	//YTr   *yandex.YTranslator
}

//func NewUI(yd *yandex.YDictionary, yt *yandex.YTranslator) *UI {
//	return &UI{yd, yt}
//}

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
	if err := g.SetKeybinding(TEXT_VIEW, gocui.KeyEnter, gocui.ModNone, ui.handleText); err != nil {
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
	info := whatlanggo.DetectLangWithOptions(textFromTextView, OPTIONS_DETECT_LANG)
	langOpt, err := createTranslateLangOpt(info)
	if err != nil {
		return err
	}
	go func() {
		g.Update(func(g *gocui.Gui) error {
			translateView, err := g.View(TRANSLATE_VIEW)
			if err != nil {
				return err
			}

			translateView.Clear()

			fmt.Fprintln(translateView, "...translated...")
			return nil
		})
	}()

	go func() {
		word, err := ui.translate(textFromTextView, langOpt.source, langOpt.destination)
		if err != nil {
			g.Update(func(g *gocui.Gui) error {
				translateView, err := g.View(TRANSLATE_VIEW)
				if err != nil {
					return err
				}
				fmt.Fprintln(translateView, "...error on translate text...")
				return nil
			})
			return
		}

		// update TRANSLATE view. Exactly search translate word in db or translate with yandex dictionary
		g.Update(func(g *gocui.Gui) error {
			translateView, err := g.View(TRANSLATE_VIEW)
			if err != nil {
				return err
			}

			translateView.Clear()

			if !word.IsEmpty() {
				fmt.Fprintln(translateView, word.Print())
			} else {
				fmt.Fprintln(translateView, "...translate not found...")
			}

			return nil
		})
	}()
	//push text to HISTORY view
	go func() {
		g.Update(func(g *gocui.Gui) error {
			historyView, err := g.View(HISTORY_VIEW)
			if err != nil {
				return err
			}

			history = append(history, textFromTextView)

			historyView.Clear()
			for i := len(history) - 1; i >= 0; i-- {
				fmt.Fprintln(historyView, history[i])
			}

			return nil
		})
	}()
	return nil
}

func (ui *UI) translate(text string, langFrom string, langTo string) (model.TranslatedText, error) {
	//mayBeWord := maybeWord(text)
	//containsCaret := strings.Contains(text, "\n")
	//check on word and not contains \n symbol
	//if mayBeWord && !containsCaret {
	tr, err := ui.YDict.Translate(text, langFrom, langTo)
	if err != nil {
		log.Printf("Error on translate word throw Dictionary, err=%v", err)
		return nil, err
	}
	//if !tr.IsEmpty() {
	//	return tr, nil
	//}

	//} else {
	//
	//}
	//now try translate throw yandex Translator
	//tr, err := ui.YTr.Translate(text, langFrom, langTo)
	//if err != nil {
	//	log.Printf("Error on translate word throw Translator, err=%v", err)
	//	return nil, err
	//}
	return tr, nil
}

func cleanView(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		return nil
	})
	return nil
}

func nextView(g *gocui.Gui, _ *gocui.View) error {
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

func quit(_ *gocui.Gui, _ *gocui.View) error {
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
		v.Autoscroll = true
		//v.Editor = &VimEditor{}

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

//func maybeWord(text string) bool {
//	if len(text) < 27 {
//		return true
//	}
//	return false
//}

//b, err := json.MarshalIndent(&word, "", "\t")
//if err != nil {
//	fmt.Println("error:", err)
//}
//
//fmt.Println(string(b))
