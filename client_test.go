package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/joho/godotenv"
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
	ctx := context.Background()

	if err := os.RemoveAll("testout"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("testout", 0666); err != nil {
		t.Fatal(err)
	}

	type TestCase struct {
		Name string
		Call func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error)
	}

	tests := []TestCase{
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
			Name: "RetrievePage",
			Call: func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
				return client._RetrievePage(ctx, "7827e04dd13a4a1682744ec55bd85c56", buffer)
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
		{
			Name: "RetrieveBlockChildren",
			Call: func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
				return client._RetrieveBlockChildren(ctx, "22a5412dd0ab4167930cb644d11fffea", buffer)
			},
		},
	}

	page, err := client.RetrievePage(ctx, "7827e04dd13a4a1682744ec55bd85c56")
	if err != nil {
		t.Fatal(err)
	}

	for _, pv := range page.Properties {
		pv := pv
		tests = append(tests, TestCase{Name: "RetrievePagePropertyItem_" + pv.Type, Call: func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
			return client._RetrievePagePropertyItem(ctx, "7827e04dd13a4a1682744ec55bd85c56", pv.ID, buffer)
		}})
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			buffer := bytes.NewBuffer(nil)
			result, err := test.Call(ctx, buffer)
			if err != nil {
				t.Fatal(err)
			}

			want := buffer.Bytes()

			got, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			want = normalize(want)
			got = normalize(got)

			if !bytes.Equal(want, got) {
				ioutil.WriteFile(fmt.Sprintf("testout/%v.want.json", test.Name), want, 0666)
				ioutil.WriteFile(fmt.Sprintf("testout/%v.got.json", test.Name), got, 0666)
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
