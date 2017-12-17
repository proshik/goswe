package view

import (
	"bytes"
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
	INFO_VIEW      = "info"
)

var VIEW_TITLES = map[string]string{
	TEXT_VIEW:      "Text",
	TRANSLATE_VIEW: "Dictionary",
	HISTORY_VIEW:   "History",
	INFO_VIEW:      "Keyboard shortcuts",
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

var translateCache = make(map[string]string)

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
	if err := g.SetKeybinding(ALL_VIEWS, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(TEXT_VIEW, gocui.KeyEnter, gocui.ModNone, ui.handleInputText); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(TEXT_VIEW, gocui.KeyEsc, gocui.ModNone, cleanView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(HISTORY_VIEW, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(HISTORY_VIEW, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(HISTORY_VIEW, gocui.KeyEnter, gocui.ModNone, handleHistoryItem); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (ui *UI) handleInputText(g *gocui.Gui, v *gocui.View) error {
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
		if val, ok := translateCache[textFromTextView]; ok {
			g.Update(func(g *gocui.Gui) error {
				translateView, err := g.View(TRANSLATE_VIEW)
				if err != nil {
					return err
				}

				translateView.Clear()

				fmt.Fprintln(translateView, val)
				return nil
			})
		} else {
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
					short, full := word.Print()
					translateCache[short] = full
					fmt.Fprintln(translateView, full)
				} else {
					fmt.Fprintln(translateView, "...translate not found...")
				}

				return nil
			})
		}
	}()
	//push text to HISTORY view
	go func() {
		g.Update(func(g *gocui.Gui) error {
			historyView, err := g.View(HISTORY_VIEW)
			if err != nil {
				return err
			}

			if len(history) > 0 {
				if history[len(history)-1] == textFromTextView {
					return nil
				}
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

	tr, err := ui.YDict.Translate(text, langFrom, langTo)
	if err != nil {
		log.Printf("Error on translate word throw Dictionary, err=%v", err)
		return nil, err
	}

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
	nextIndex := (activeIndex + 1) % (len(VIEW_TITLES) - 1)

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

func cursorDown(_ *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(_ *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleHistoryItem(g *gocui.Gui, v *gocui.View) error {

	_, cy := v.Cursor()

	line, err := v.Line(cy)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	g.Update(func(g *gocui.Gui) error {
		translateView, err := g.View(TRANSLATE_VIEW)
		if err != nil {
			return err
		}

		word := translateCache[line]

		translateView.Clear()

		fmt.Fprintln(translateView, word)
		return nil
	})

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

	if v, err := g.SetView(TRANSLATE_VIEW, maxX/2, 0, maxX-1, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = VIEW_TITLES[TRANSLATE_VIEW]
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView(INFO_VIEW, maxX/2, maxY-7, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = VIEW_TITLES[INFO_VIEW]
		v.Wrap = true

		go g.Update(func(g *gocui.Gui) error {
			infoView, err := g.View(INFO_VIEW)
			if err != nil {
				return err
			}

			buf := bytes.Buffer{}
			buf.WriteString("<Tab> - change view\n")
			buf.WriteString("<Enter> - translate text on Text view\n")
			buf.WriteString("<Enter> - select text on History view\n")
			buf.WriteString("<↑,↓> - navigate inside view\n")
			buf.WriteString("<Ctrl+C> - exit\n")
			fmt.Fprintln(infoView, buf.String())

			return nil
		})
	}

	if v, err := g.SetView(HISTORY_VIEW, 0, maxY/2, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = VIEW_TITLES[HISTORY_VIEW]
		v.Wrap = true
		v.Autoscroll = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
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
