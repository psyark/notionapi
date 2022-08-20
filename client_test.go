package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/yudai/gojsondiff"
)

var (
	listener = &TestListener{}
	client   *Client
)

func init() {
	if err := godotenv.Load("credentials.env"); err != nil {
		panic(err)
	}
	client = NewClient(os.Getenv("NOTION_API_KEY")).SetListener(listener)
}

func TestRetrieveDatabase(t *testing.T) {
	ctx := context.Background()
	const databaseID = "8b6685786cc647ecb614dbd9b3ee5113"
	data, err := useCache(fmt.Sprintf(".cache/%v.json", databaseID), func() ([]byte, error) {
		_, err := client.RetrieveDatabase(ctx, databaseID)
		if err != nil {
			return nil, err
		}
		return listener.body, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := check(data, &Database{}); err != nil {
		t.Fatal(err)
	}
}

func TestRetrievePagePropertyItem(t *testing.T) {
	ctx := context.Background()
	const pageID = "7827e04dd13a4a1682744ec55bd85c56"
	data, err := useCache(fmt.Sprintf(".cache/%v.json", pageID), func() ([]byte, error) {
		_, err := client.RetrievePage(ctx, pageID)
		if err != nil {
			return nil, err
		}
		return listener.body, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	page := &Page{}
	if err := check(data, page); err != nil {
		t.Fatal(err)
	}

	for k, pv := range page.Properties {
		data, err := useCache(fmt.Sprintf(".cache/%v_%v.json", pageID, k), func() ([]byte, error) {
			_, err := client.RetrievePagePropertyItem(ctx, pageID, pv.ID)
			if err != nil {
				return nil, err
			}
			return listener.body, nil
		})
		if err != nil {
			t.Fatal(err)
		}

		if err := check(data, &PropertyItemOrPagination{}); err != nil {
			t.Fatal(err)
		}
	}
}

func check(data []byte, result interface{}) error {
	if err := json.Unmarshal(data, result); err != nil {
		return err
	}

	data2, err := json.Marshal(result)
	if err != nil {
		return err
	}

	diff, err := gojsondiff.New().Compare(data, data2)
	if err != nil {
		return err
	}

	if diff.Modified() {
		return validationError{data, data2, diff}
	}
	return nil
}

func useCache(fileName string, ifNotExists func() ([]byte, error)) ([]byte, error) {
	if _, err := os.Stat(fileName); err != nil {
		data, err := ifNotExists()
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(fileName, data, 0666); err != nil {
			return nil, err
		}

		return data, nil
	}

	return ioutil.ReadFile(fileName)
}
