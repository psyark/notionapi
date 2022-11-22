package mapping

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/psyark/notionapi"
)

type Object struct {
	Title_String    string                  `notion:"title"`  // 名前
	Title_Raw       notionapi.RichTextArray `notion:"title"`  // 名前
	RichText_String string                  `notion:"%40RTE"` // テキスト
	// RichText_Raw    notionapi.RichTextArray `notion:"%40RTE"` // テキスト
	Email       string               `notion:"hclY"`     // メール
	URL         string               `notion:"Udz%3F"`   // URL
	PhoneNumber string               `notion:"qjI%3B"`   // 電話
	Number      float64              `notion:"p%7Bq%3E"` // 数値
	Checkbox    bool                 `notion:"%3DqUq"`   // チェックボックス
	Date_Time   time.Time            `notion:"OL%3C%3F"` // 日付
	Date_Raw    *notionapi.DateValue `notion:"OL%3C%3F"` // 日付
	// LastEditedTime_Interface interface{} `notion:"CHbM"`     // 最終更新日時
	// CreatedTime_Interface    interface{} `notion:"~gd%5C"`   // 作成日時
	// Select           interface{} `notion:"rMGi"`     // セレクト
	// Multi_select     interface{} `notion:"%3Ewkp"`   // マルチセレクト
	// Status           interface{} `notion:"GCe%3C"`   // ステータス
	// Files            interface{} `notion:"sAR_"`     // ファイル&メディア
	// Relation1        interface{} `notion:"wpAL"`     // リレーション1
	// Relation2        interface{} `notion:"NBWw"`     // リレーション2
	// Rollup           interface{} `notion:"Fo%7DT"`   // ロールアップ
	// Formula          interface{} `notion:"ZTlY"`     // 関数
	// People           interface{} `notion:"BvQz"`     // ユーザー
	// Created_by       interface{} `notion:"ig%40H"`   // 作成者
	// Last_edited_by   interface{} `notion:"FY%5B%7C"` // 最終更新者
}

const databaseID = "8b6685786cc647ecb614dbd9b3ee5113"

var client *notionapi.Client

func init() {
	if err := godotenv.Load("../credentials.env"); err != nil {
		panic(err)
	}
	client = notionapi.NewClient(os.Getenv("NOTION_API_KEY"))
}

func TestXxx(t *testing.T) {
	ctx := context.Background()

	// db, _ := client.RetrieveDatabase(ctx, databaseID)
	// Create(Object{}, db.Properties)

	pagi, err := client.QueryDatabase(ctx, databaseID, &notionapi.QueryDatabaseOptions{})
	if err != nil {
		t.Fatal(err)
	}

	for _, page := range pagi.Results {
		obj := Object{}
		if err := Decode(page, &obj); err != nil {
			t.Fatal(err)
		}
		// d, _ := json.MarshalIndent(obj, "", "  ")
		// fmt.Println(string(d))

		if page.ID == "7827e04d-d13a-4a16-8274-4ec55bd85c56" {
			obj.Number += 1
			obj.RichText_String += " HOGE "

			opt, err := GetUpdatePageOptions(page, obj)
			if err != nil {
				t.Fatal(err)
			}

			d, _ := json.MarshalIndent(opt, "", "  ")
			fmt.Println(string(d))
		}
	}
}
