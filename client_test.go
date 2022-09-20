package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/yudai/gojsondiff"
)

var (
	client *Client
)

func init() {
	if err := godotenv.Load("credentials.env"); err != nil {
		panic(err)
	}
	client = NewClient(os.Getenv("NOTION_API_KEY"))
}

func TestClient(t *testing.T) {
	if err := os.RemoveAll("testout"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("testout", 0666); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name string
		Call func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error)
	}{
		{
			Name: "RetrieveDatabase",
			Call: func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
				return client._RetrieveDatabase(ctx, "8b6685786cc647ecb614dbd9b3ee5113", buffer)
			},
		},
		{
			Name: "QueryDatabase",
			Call: func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
				return client._QueryDatabase(ctx, "8b6685786cc647ecb614dbd9b3ee5113", &QueryDatabaseOptions{}, buffer)
			},
		},
		{
			Name: "RetrieveBlockChildren",
			Call: func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
				return client._RetrieveBlockChildren(ctx, "22a5412dd0ab4167930cb644d11fffea", buffer)
			},
		},
		{
			Name: "UpdatePage",
			Call: func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
				emojis := []string{"🍰", "🍣", "🍜", "🍤", "🥗"}
				opt := &UpdatePageOptions{Icon: &FileOrEmoji{Type: "emoji"}}
				opt.Icon.Emoji = emojis[rand.Intn(len(emojis))]
				return client._UpdatePage(ctx, "7827e04dd13a4a1682744ec55bd85c56", opt, buffer)
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			ctx := context.Background()
			buffer := bytes.NewBuffer(nil)
			result, err := test.Call(ctx, buffer)
			if err != nil {
				t.Fatal(err)
			}

			remarshal, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			diff, err := gojsondiff.New().Compare(buffer.Bytes(), remarshal)
			if err != nil {
				t.Fatal(err)
			}

			if diff.Modified() {
				dm := diffMap{}

				// res, _ := formatter.NewDeltaFormatter().Format(e.diff)
				for _, delta := range diff.Deltas() {
					dm.add(delta, "")
				}

				diffBytes, _ := json.MarshalIndent(dm, "", "  ")

				ioutil.WriteFile(fmt.Sprintf("testout/%v.want.json", test.Name), normalize(buffer.Bytes()), 0666)
				ioutil.WriteFile(fmt.Sprintf("testout/%v.got.json", test.Name), normalize(remarshal), 0666)
				ioutil.WriteFile(fmt.Sprintf("testout/%v.diff.json", test.Name), diffBytes, 0666)
				t.Error("validation failed")
			}
		})
	}
}

func normalize(src []byte) []byte {
	tmp := map[string]interface{}{}
	json.Unmarshal(src, &tmp)
	out, _ := json.MarshalIndent(tmp, "", "  ")
	return out
}

func TestRetrievePagePropertyItem(t *testing.T) {
	ctx := context.Background()
	const pageID = "7827e04dd13a4a1682744ec55bd85c56"
	data, err := useCache(fmt.Sprintf(".cache/%v.json", pageID), func() ([]byte, error) {
		buffer := bytes.NewBuffer(nil)
		_, err := client._RetrievePage(ctx, pageID, buffer)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
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
			buffer := bytes.NewBuffer(nil)
			_, err := client._RetrievePagePropertyItem(ctx, pageID, pv.ID, buffer)
			if err != nil {
				return nil, err
			}
			return buffer.Bytes(), nil
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

type validationError struct {
	source      []byte
	remarshaled []byte
	diff        gojsondiff.Diff
}

func (e validationError) Error() string {
	dm := diffMap{}

	// res, _ := formatter.NewDeltaFormatter().Format(e.diff)
	for _, delta := range e.diff.Deltas() {
		dm.add(delta, "")
	}

	diffStr, _ := json.MarshalIndent(dm, "", "  ")
	buffer := bytes.NewBuffer(nil)
	json.Indent(buffer, e.source, "", "  ")
	return fmt.Sprintf("validation failed. diff: %v\nwant: %v\ngot: %v", string(diffStr), buffer.String(), string(e.remarshaled))
}

type diffMap map[string]interface{}

func (dm diffMap) add(delta gojsondiff.Delta, prefix string) {
	switch delta := delta.(type) {
	case *gojsondiff.Added:
		dm[prefix+delta.Position.String()] = struct {
			Value interface{} `json:"ADDED"`
		}{delta.Value}
	case *gojsondiff.Deleted:
		dm[prefix+delta.Position.String()] = struct {
			Value interface{} `json:"DELETED"`
		}{delta.Value}
	case *gojsondiff.Modified:
		dm[prefix+delta.Position.String()] = struct {
			OldValue interface{} `json:"MODIFIED_FROM"`
			NewValue interface{} `json:"MODIFIED_TO"`
		}{delta.OldValue, delta.NewValue}
	case *gojsondiff.Array:
		for _, d := range delta.Deltas {
			dm.add(d, prefix+delta.Position.String()+".")
		}
	case *gojsondiff.Object:
		for _, d := range delta.Deltas {
			dm.add(d, prefix+delta.Position.String()+".")
		}
	default:
		panic(reflect.TypeOf(delta))
	}
}
