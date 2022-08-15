package codegen

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ParseCallback func(title string, desc string, props []DocProp) error

func cache(url string) (*goquery.Document, error) {
	cacheName := fmt.Sprintf(".cache/%x.html", md5.Sum([]byte(url)))

	if _, err := os.Stat(cacheName); err != nil {
		resp, err := http.DefaultClient.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(cacheName, data, 0666); err != nil {
			return nil, err
		}
	}

	f, err := os.Open(cacheName)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return goquery.NewDocumentFromReader(f)
}

func Parse(url string, callback ParseCallback) error {
	doc, err := cache(url)
	if err != nil {
		return err
	}

	var loopError error
	doc.Find(`.markdown-body`).ChildrenFiltered(`h1,h2,h3`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Text()
		desc := strings.Join(s.NextFilteredUntil(`p,ul,blockquote`, `.rdmd-table,h1,h2,h3`).Map(func(i int, s *goquery.Selection) string { return s.Text() }), "\n")

		fields := []DocProp{}
		s.NextFilteredUntil(`.rdmd-table`, `h1,h2,h3`).Each(func(i int, s *goquery.Selection) {
			columnMap := map[string]int{}
			s.Find(`thead tr th`).Each(func(i int, s *goquery.Selection) {
				columnMap[s.Text()] = i
			})
			s.Find(`tbody tr`).Each(func(i int, s *goquery.Selection) {
				fields = append(fields, DocProp{
					Name:        s.Find(`td`).Eq(columnMap["Property"]).Text(),
					Type:        s.Find(`td`).Eq(columnMap["Type"]).Text(),
					Description: s.Find(`td`).Eq(columnMap["Description"]).Text(),
				})
			})
		})

		if err := callback(title, desc, fields); err != nil {
			loopError = err
			return false
		}
		return true
	})

	return loopError
}
