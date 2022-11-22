package mapping

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/psyark/notionapi"
)

type Hoge struct {
	Title_Interface       interface{} `notion:"title"`    // 名前
	Title_String          string      `notion:"title"`    // 名前
	RichText_Interface    interface{} `notion:"%40RTE"`   // テキスト
	RichText_String       string      `notion:"%40RTE"`   // テキスト
	Number_Interface      interface{} `notion:"p%7Bq%3E"` // 数値
	Number_Float64        float64     `notion:"p%7Bq%3E"` // 数値
	Number_Int            int         `notion:"p%7Bq%3E"` // 数値
	Checkbox_Interface    interface{} `notion:"%3DqUq"`   // チェックボックス
	Checkbox_Bool         bool        `notion:"%3DqUq"`   // チェックボックス
	Email                 interface{} `notion:"hclY"`     // メール
	URL_Inteface          interface{} `notion:"Udz%3F"`   // URL
	URL_String            string      `notion:"Udz%3F"`   // URL
	PhoneNumber_Interface interface{} `notion:"qjI%3B"`   // 電話
	PhoneNumber_String    string      `notion:"qjI%3B"`   // 電話
	// Date_Interface        interface{} `notion:"OL%3C%3F"` // 日付
	// Date_String           string      `notion:"OL%3C%3F"` // 日付
	// Date_Time             time.Time   `notion:"OL%3C%3F"` // 日付
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

	pagi, err := client.QueryDatabase(ctx, databaseID, &notionapi.QueryDatabaseOptions{})
	if err != nil {
		t.Fatal(err)
	}

	for _, page := range pagi.Results {
		hoge := Hoge{}
		if err := Decode(page, &hoge); err != nil {
			t.Fatal(err)
		}
		d, _ := json.MarshalIndent(hoge, "", "  ")
		fmt.Println(string(d))
	}
}
