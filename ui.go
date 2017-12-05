package main

import (
	"github.com/jroimartin/gocui"
	"strings"
	"fmt"
)

type UI struct {
	YandexDict *YandexDict
	DBConnect  *DBConnect
}

func NewUI(yd *YandexDict, dbc *DBConnect) *UI {
	return NewUI(yd, dbc)
}

func Run

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
