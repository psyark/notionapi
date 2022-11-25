package mapping

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/psyark/notionapi"
)

type Object struct {
	Title_Raw       notionapi.RichTextArray `notion:"title"`  // 名前
	Title_String    string                  `notion:"title"`  // 名前
	RichText_String string                  `notion:"%40RTE"` // テキスト
	// RichText_Raw    notionapi.RichTextArray `notion:"%40RTE"` // テキスト
	Email       string  `notion:"hclY"`     // メール
	URL         string  `notion:"Udz%3F"`   // URL
	PhoneNumber string  `notion:"qjI%3B"`   // 電話
	Number      float64 `notion:"p%7Bq%3E"` // 数値
	Checkbox    bool    `notion:"%3DqUq"`   // チェックボックス
	// Date_Time   time.Time            `notion:"OL%3C%3F"` // 日付
	Date_Raw *notionapi.DateValue `notion:"OL%3C%3F"` // 日付
	// LastEditedTime_Interface interface{} `notion:"CHbM"`     // 最終更新日時
	// CreatedTime_Interface    interface{} `notion:"~gd%5C"`   // 作成日時
	// Select           interface{} `notion:"rMGi"`     // セレクト
	// Multi_select     interface{} `notion:"%3Ewkp"`   // マルチセレクト
	// Status           interface{} `notion:"GCe%3C"`   // ステータス
	// Files            interface{} `notion:"sAR_"`     // ファイル&メディア
	Relation1 []uuid.UUID `notion:"wpAL"` // リレーション1
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

func ExampleDecodePage() {
	ctx := context.Background()

	// db, _ := client.RetrieveDatabase(ctx, databaseID)
	// Create(Object{}, db.Properties)

	pagi, err := client.QueryDatabase(ctx, databaseID, &notionapi.QueryDatabaseOptions{})
	if err != nil {
		panic(err)
	}

	for _, page := range pagi.Results {
		obj := Object{}
		if err := DecodePage(&obj, page); err != nil {
			panic(err)
		}

		d, _ := json.Marshal(obj)
		fmt.Println(string(d))

		// if page.ID.String() == "7827e04d-d13a-4a16-8274-4ec55bd85c56" {
		// 	obj.Number += 1
		// 	obj.RichText_String += " HOGE "
		// 	obj.Date_Raw.Start = "2050-01-02"

		// 	opt, err := UpdatePageFrom(obj, page)
		// 	if err != nil {
		// 		t.Fatal(err)
		// 	}

		// 	d, _ := json.MarshalIndent(opt, "", "  ")
		// 	fmt.Println(string(d))
		// }
	}

	// Output:
	// {"Title_Raw":[],"Title_String":"","RichText_String":"text","Email":"","URL":"","PhoneNumber":"","Number":0,"Checkbox":false,"Date_Raw":null,"Relation1":[]}
	// {"Title_Raw":[{"annotations":{"bold":false,"code":false,"color":"default","italic":false,"strikethrough":false,"underline":false},"href":null,"plain_text":"Item 1","text":{"content":"Item 1","link":null},"type":"text"}],"Title_String":"Item 1","RichText_String":"The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox","Email":"me@example.com","URL":"http://example.com","PhoneNumber":"123-456","Number":123,"Checkbox":true,"Date_Raw":{"start":"2022-08-09","end":null,"time_zone":null},"Relation1":["8b1cf30a-f939-4c1e-a09b-b63bcc890569"]}
	// {"Title_Raw":[{"annotations":{"bold":false,"code":false,"color":"default","italic":false,"strikethrough":false,"underline":false},"href":null,"plain_text":"Item 3","text":{"content":"Item 3","link":null},"type":"text"}],"Title_String":"Item 3","RichText_String":"","Email":"","URL":"","PhoneNumber":"","Number":1.2143432785732895e+27,"Checkbox":false,"Date_Raw":{"start":"2022-08-09T00:00:00.000+09:00","end":null,"time_zone":null},"Relation1":["ba53d412-b627-4e3d-9e2e-425e20988010"]}
	// {"Title_Raw":[{"annotations":{"bold":false,"code":false,"color":"default","italic":false,"strikethrough":false,"underline":false},"href":null,"plain_text":"Item 2","text":{"content":"Item 2","link":null},"type":"text"}],"Title_String":"Item 2","RichText_String":"Text Page Link  Web Link Bold Italic Underline Strike Code Formula Red @Keiichi Yoshikawa 2022-09-28 ","Email":"","URL":"","PhoneNumber":"","Number":45.67,"Checkbox":false,"Date_Raw":{"start":"2022-08-09","end":"2022-08-11","time_zone":null},"Relation1":[]}
}

func ExampleUpdatePageFrom() {
	ctx := context.Background()
	page, err := client.RetrievePage(ctx, "7827e04d-d13a-4a16-8274-4ec55bd85c56")
	if err != nil {
		panic(err)
	}

	obj := Object{}
	if err := DecodePage(&obj, *page); err != nil {
		panic(err)
	}

	obj.Number += 1
	obj.RichText_String += " HOGE "
	obj.Date_Raw.Start = "2050-01-02"

	opt, err := UpdatePageFrom(obj, *page)
	if err != nil {
		panic(err)
	}

	d, _ := json.MarshalIndent(opt, "", "  ")
	fmt.Println(string(d))
	// Output:
	// {
	//   "properties": {
	//     "%40RTE": {
	//       "rich_text": [
	//         {
	//           "href": null,
	//           "plain_text": "",
	//           "text": {
	//             "content": "The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox The quick brown fox HOGE ",
	//             "link": null
	//           },
	//           "type": "text"
	//         }
	//       ],
	//       "type": "rich_text"
	//     },
	//     "OL%3C%3F": {
	//       "date": {
	//         "end": null,
	//         "start": "2050-01-02",
	//         "time_zone": null
	//       },
	//       "type": "date"
	//     },
	//     "p%7Bq%3E": {
	//       "number": 124,
	//       "type": "number"
	//     }
	//   }
	// }
}

func ExampleCreate() {
	ctx := context.Background()

	db, err := client.RetrieveDatabase(ctx, databaseID)
	if err != nil {
		panic(err)
	}

	obj := Object{
		Title_String:    "TTL",
		RichText_String: "RICH",
		Number:          3.14159265,
		Date_Raw:        &notionapi.DateValue{Start: "2022-11-25"},
		Relation1:       []uuid.UUID{uuid.MustParse("8b668578-6cc6-47ec-b614-dbd9b3ee5113")},
	}

	opt, err := CreatePageFrom(obj, db)
	if err != nil {
		panic(err)
	}

	d, _ := json.MarshalIndent(opt, "", "  ")
	fmt.Println(string(d))

	// Output:
	// {
	//   "parent": {
	//     "database_id": "8b668578-6cc6-47ec-b614-dbd9b3ee5113",
	//     "type": "database_id"
	//   },
	//   "properties": {
	//     "%40RTE": {
	//       "rich_text": [
	//         {
	//           "href": null,
	//           "plain_text": "",
	//           "text": {
	//             "content": "RICH",
	//             "link": null
	//           },
	//           "type": "text"
	//         }
	//       ],
	//       "type": "rich_text"
	//     },
	//     "OL%3C%3F": {
	//       "date": {
	//         "end": null,
	//         "start": "2022-11-25",
	//         "time_zone": null
	//       },
	//       "type": "date"
	//     },
	//     "p%7Bq%3E": {
	//       "number": 3.14159265,
	//       "type": "number"
	//     },
	//     "title": {
	//       "title": [
	//         {
	//           "href": null,
	//           "plain_text": "",
	//           "text": {
	//             "content": "TTL",
	//             "link": null
	//           },
	//           "type": "text"
	//         }
	//       ],
	//       "type": "title"
	//     },
	//     "wpAL": {
	//       "relation": [
	//         {
	//           "id": "8b668578-6cc6-47ec-b614-dbd9b3ee5113"
	//         }
	//       ],
	//       "type": "relation"
	//     }
	//   }
	// }
}
