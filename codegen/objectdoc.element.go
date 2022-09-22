package codegen

import (
	"encoding/json"
	"fmt"
)

type ObjectDocElement interface {
	objectDocElement()
}

type HeadingElement struct {
	Text string
}

type ParagraphElement struct {
	Content string
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

type BlockParameter struct {
	Property     string
	Type         string
	Description  string
	ExampleValue string `json:"Example value"`
}

type BlockParametersElement []BlockParameter

func (t *BlockParametersElement) UnmarshalJSON(data []byte) error {
	raw := struct {
		Data map[string]string `json:"data"`
		Cols int               `json:"cols"`
		Rows int               `json:"rows"`
	}{}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	mapSlice := make([]map[string]string, raw.Rows)
	for r := range mapSlice {
		m := map[string]string{}
		for c := 0; c < raw.Cols; c++ {
			m[raw.Data[fmt.Sprintf("h-%d", c)]] = raw.Data[fmt.Sprintf("%d-%d", r, c)]
		}
		mapSlice[r] = m
	}

	data, err := json.Marshal(mapSlice)
	if err != nil {
		return err
	}

	type Alias BlockParametersElement
	return json.Unmarshal(data, (*Alias)(t))
}

func (t *HeadingElement) objectDocElement()         {}
func (t *ParagraphElement) objectDocElement()       {}
func (t *BlockCodeElement) objectDocElement()       {}
func (t *BlockCalloutElement) objectDocElement()    {}
func (t *BlockParametersElement) objectDocElement() {}
