package codegen

import "fmt"

type ObjectDocElement interface {
	objectDocElement()
}

type HeadingElement struct {
	text string
}

type ParagraphElement struct {
	content string
}

type BlockCodeElement struct {
	Codes []struct {
		Name     string `json:"string"`
		Language string `json:"language"`
		Code     string `json:"code"`
	} `json:"codes"`
}

type BlockCalloutElement struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type BlockParametersElement struct {
	Data map[string]string `json:"data"`
	Cols int               `json:"cols"`
	Rows int               `json:"rows"`
}

func (t *BlockParametersElement) MapSlice() []map[string]string {
	slice := make([]map[string]string, t.Rows)

	for r := range slice {
		m := map[string]string{}
		for c := 0; c < t.Cols; c++ {
			m[t.Data[fmt.Sprintf("h-%d", c)]] = t.Data[fmt.Sprintf("%d-%d", r, c)]
		}
		slice[r] = m
	}

	return slice
}

func (t *HeadingElement) objectDocElement()         {}
func (t *ParagraphElement) objectDocElement()       {}
func (t *BlockCodeElement) objectDocElement()       {}
func (t *BlockCalloutElement) objectDocElement()    {}
func (t *BlockParametersElement) objectDocElement() {}
