package codegen

import (
	"io"
)

// ObjectDocSectionはNotionのObject関連ドキュメントで隣接して出現する要素のグループです
// 最大1つのHeadingTokenを含みます
type ObjectDocSection struct {
	Heading  *HeadingElement
	Elements []ObjectDocElement
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
