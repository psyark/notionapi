package codegen

import (
	"io"
	"strings"
)

// ObjectDocSectionはNotionのObject関連ドキュメントで隣接して出現する要素のグループです
// 最大1つのHeadingTokenを含みます
type ObjectDocSection struct {
	Heading  *HeadingElement
	Elements []ObjectDocElement
}

// TODO: AllParagraphTextにリネーム
func (s *ObjectDocSection) ParagraphText() string {
	texts := []string{}
	for _, e := range s.Elements {
		if e, ok := e.(*ParagraphElement); ok {
			texts = append(texts, e.Content)
		}
	}
	return strings.TrimSpace(strings.Join(texts, "\n"))
}

func (s *ObjectDocSection) FirstParagraphText() string {
	for _, e := range s.Elements {
		if e, ok := e.(*ParagraphElement); ok {
			return e.Content
		}
	}
	return ""
}

func (s *ObjectDocSection) Parameters() []ObjectDocParameter {
	params := []ObjectDocParameter{}
	for _, e := range s.Elements {
		if e, ok := e.(*BlockParametersElement); ok {
			params = append(params, *e...)
		}
	}
	return params
}

// ParseObjectDoc は NotionのObject関連ドキュメントをパースしてセクションのスライスを返します
func ParseObjectDoc(url string) ([]ObjectDocSection, error) {
	tokenizer, err := newObjectDocTokenizer(url)
	if err != nil {
		return nil, err
	}

	sections := []ObjectDocSection{{}}
	for {
		element, err := tokenizer.next()
		if err == io.EOF {
			return sections, nil
		} else if err != nil {
			return nil, err
		}
		if heading, ok := element.(*HeadingElement); ok {
			sections = append(sections, ObjectDocSection{Heading: heading})
		} else {
			current := &sections[len(sections)-1]
			current.Elements = append(current.Elements, element)
		}
	}
}
