package codegen

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type objectDocTokenizer struct {
	lines []string
	index int
}

func (t *objectDocTokenizer) next() (ObjectDocElement, error) {
	if t.index >= len(t.lines) {
		return nil, io.EOF
	}

	switch t.lines[t.index] {
	case "[block:code]", "[block:callout]", "[block:parameters]":
		var block ObjectDocElement
		switch t.lines[t.index] {
		case "[block:code]":
			block = &BlockCodeElement{}
		case "[block:callout]":
			block = &BlockCalloutElement{}
		case "[block:parameters]":
			block = &BlockParametersElement{}
		}

		startIndex := t.index
		for t.index < len(t.lines) {
			if t.lines[t.index] == "[/block]" {
				t.index++

				content := []byte(strings.Join(t.lines[startIndex+1:t.index-1], "\n"))
				if err := json.Unmarshal(content, block); err != nil {
					return nil, err
				}
				return block, nil
			}
			t.index++
		}
		return nil, fmt.Errorf("[/block] not exists")
	default:
		if strings.HasPrefix(t.lines[t.index], "#") {
			token := &HeadingElement{t.lines[t.index]}
			t.index++
			return token, nil
		} else {
			startIndex := t.index
			for t.index < len(t.lines) && t.isParagraph(t.lines[t.index]) {
				t.index++
			}
			return &ParagraphElement{strings.Join(t.lines[startIndex:t.index], "\n")}, nil
		}
	}
}

func (t *objectDocTokenizer) isParagraph(line string) bool {
	switch line {
	case "[block:code]", "[block:callout]", "[block:parameters]":
		return false
	}
	return !strings.HasPrefix(line, "#")
}

func newObjectDocTokenizer(url string) (*objectDocTokenizer, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	ssrPropsBytes := []byte(doc.Find(`#ssr-props`).AttrOr("data-initial-props", ""))
	ssrProps := struct {
		Doc struct {
			Body string `json:"body"`
		} `json:"doc"`
	}{}
	if err := json.Unmarshal(ssrPropsBytes, &ssrProps); err != nil {
		return nil, err
	}

	lines := strings.Split(ssrProps.Doc.Body, "\n")
	return &objectDocTokenizer{lines, 0}, nil
}
