package model

import (
	"bytes"
	"fmt"
)

type TranslatedText interface {
	Print() (string, string)
	IsEmpty() bool
}

//For Yandex translator API
type Translate struct {
	Code int      `json:"code"`
	Lang string   `json:"lang"`
	Text []string `json:"text"`
}

func (w *Translate) Print() (string, string) {
	if len(w.Text[0]) > 25 {
		return w.Text[0][0:25], w.Text[0]
	} else {
		return w.Text[0], w.Text[0]
	}
}

func (w *Translate) IsEmpty() bool {
	if w != nil && len(w.Text) > 0 {
		return false
	}
	return true
}

//For Yandex dictionary API
type Dictionary struct {
	Head Head  `json:"head"`
	Def  []Def `json:"def"`
}

func (w *Dictionary) Print() (string, string) {
	var buf bytes.Buffer

	buf.WriteString(w.Def[0].Tr[0].Text + "\n")
	buf.WriteString("\n")

	for _, w := range w.Def {
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
	return w.Def[0].Text, buf.String()
}

func (w *Dictionary) IsEmpty() bool {
	if w != nil && w.Def != nil && len(w.Def) > 0 {
		return false
	}
	return true
}

type Head struct{}

type Def struct {
	Text string `json:"text"`
	Pos  string `json:"pos"`
	Ts   string `json:"ts"`
	Tr   []Tr   `json:"tr"`
}

type Tr struct {
	Attr
	Syn  []Syn  `json:"syn"`
	Mean []Mean `json:"mean"`
	Ex   []Ex   `json:"ex"`
}

type Syn struct {
	Attr
}

type Mean struct {
	Attr
}

type Ex struct {
	Attr
	Tr
}

type Attr struct {
	Text string `json:"text"`
	Pos  string `json:"pos"`
	Gen  string `json:"gen"`
}
