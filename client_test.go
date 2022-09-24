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

func TestClient(t *testing.T) {
	ctx := context.Background()

	if err := os.RemoveAll("testout"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("testout", 0777); err != nil {
		t.Fatal(err)
	}

	if err := godotenv.Load("credentials.env"); err != nil {
		t.Fatal(err)
	}

	client := NewClient(os.Getenv("NOTION_API_KEY"))

	tests := map[string]func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error){
		"RetrieveDatabase": func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
			return client._RetrieveDatabase(ctx, "8b6685786cc647ecb614dbd9b3ee5113", buffer)
		},
		"QueryDatabase": func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
			return client._QueryDatabase(ctx, "8b6685786cc647ecb614dbd9b3ee5113", &QueryDatabaseOptions{}, buffer)
		},
		"RetrievePage": func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
			return client._RetrievePage(ctx, "7827e04dd13a4a1682744ec55bd85c56", buffer)
		},
		"UpdatePage": func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
			emojis := []string{"🍰", "🍣", "🍜", "🍤", "🥗"}
			opt := &UpdatePageOptions{Icon: &FileOrEmoji{Emoji: &Emoji{Type: "emoji", Emoji: emojis[rand.Intn(len(emojis))]}}}
			return client._UpdatePage(ctx, "7827e04dd13a4a1682744ec55bd85c56", opt, buffer)
		},
		"RetrieveBlockChildren": func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
			return client._RetrieveBlockChildren(ctx, "22a5412dd0ab4167930cb644d11fffea", buffer)
		},
	}

	if page, err := client.RetrievePage(ctx, "7827e04dd13a4a1682744ec55bd85c56"); err != nil {
		t.Fatal(err)
	} else {
		for _, pv := range page.Properties {
			pv := pv
			tests["RetrievePagePropertyItem."+pv.Type] = func(ctx context.Context, buffer *bytes.Buffer) (interface{}, error) {
				return client._RetrievePagePropertyItem(ctx, "7827e04dd13a4a1682744ec55bd85c56", pv.ID, buffer)
			}
		}
	}

	for testName, testFunc := range tests {
		testName := testName
		testFunc := testFunc
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			buffer := bytes.NewBuffer(nil)
			result, err := testFunc(ctx, buffer)
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
				if err := ioutil.WriteFile(fmt.Sprintf("testout/%v.want.json", testName), want, 0666); err != nil {
					t.Error(err)
				}
				if err := ioutil.WriteFile(fmt.Sprintf("testout/%v.got.json", testName), got, 0666); err != nil {
					t.Error(err)
				}
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
