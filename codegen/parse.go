package codegen

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

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
