package main

type RawWord struct {
	Text        string `json:"text"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}

type Word struct {
	Text        string `json:"text"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	Translate   []Def  `json:"translate"`
}

func (w *Word) isEmpty() bool {
	if w != nil && w.Translate != nil && len(w.Translate) > 0 {
		return false
	}
	return true
}

type Translate struct {
	Head Head  `json:"head"`
	Def  []Def `json:"def"`
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
