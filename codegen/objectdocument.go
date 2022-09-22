package codegen

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseObjectDocument は NotionのObject関連ドキュメントをパースしてトークンのスライスを返します
func ParseObjectDocument(url string) ([]token, error) {
	tokenizer, err := newObjectDocumentTokenizer(url)
	if err != nil {
		return nil, err
	}

	tokens := []token{}
	for {
		t, err := tokenizer.next()
		if err == io.EOF {
			return tokens, nil
		} else if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
}

type objectDocumentTokenizer struct {
	lines []string
	index int
}

type token interface {
	token()
}

type headingToken struct {
	text string
}

type paragraphToken struct {
	content string
}

type blockCodeToken struct {
	Codes []struct {
		Name     string `json:"string"`
		Language string `json:"language"`
		Code     string `json:"code"`
	} `json:"codes"`
}
type blockCalloutToken struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
type blockParametersToken struct {
	Data map[string]string `json:"data"`
	Cols int               `json:"cols"`
	Rows int               `json:"rows"`
}

func (t *blockParametersToken) MapSlice() []map[string]string {
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

func (t *headingToken) token()         {}
func (t *paragraphToken) token()       {}
func (t *blockCodeToken) token()       {}
func (t *blockCalloutToken) token()    {}
func (t *blockParametersToken) token() {}

func (t *objectDocumentTokenizer) next() (token, error) {
	if t.index >= len(t.lines) {
		return nil, io.EOF
	}

	switch t.lines[t.index] {
	case "[block:code]", "[block:callout]", "[block:parameters]":
		var block token
		switch t.lines[t.index] {
		case "[block:code]":
			block = &blockCodeToken{}
		case "[block:callout]":
			block = &blockCalloutToken{}
		case "[block:parameters]":
			block = &blockParametersToken{}
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
			token := &headingToken{t.lines[t.index]}
			t.index++
			return token, nil
		} else {
			startIndex := t.index
			for t.index < len(t.lines) && t.isParagraph(t.lines[t.index]) {
				t.index++
			}
			return &paragraphToken{strings.Join(t.lines[startIndex:t.index], "\n")}, nil
		}
	}
}

func (t *objectDocumentTokenizer) isParagraph(line string) bool {
	switch line {
	case "[block:code]", "[block:callout]", "[block:parameters]":
		return false
	}
	return !strings.HasPrefix(line, "#")
}

func newObjectDocumentTokenizer(url string) (*objectDocumentTokenizer, error) {
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
	return &objectDocumentTokenizer{lines, 0}, nil
}
