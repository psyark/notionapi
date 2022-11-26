package mapping

import (
	"reflect"
	"testing"

	"github.com/psyark/notionapi"
)

type Sample struct {
	Title string                `notion:"title"`
	Icon  notionapi.FileOrEmoji `notion:",icon"`
	Cover *notionapi.File       `notion:",cover"`
}

func TestTag(t *testing.T) {
	cases := []struct {
		Type reflect.Type
		Want *tagInfo
	}{
		{reflect.TypeOf(struct {
			Title string `json:"title"`
		}{}), nil},
		{reflect.TypeOf(struct {
			Title string `notion:"title"`
		}{}), &tagInfo{name: "title"}},
		{reflect.TypeOf(struct {
			Icon notionapi.FileOrEmoji `notion:",icon"`
		}{}), &tagInfo{icon: true}},
		{reflect.TypeOf(struct {
			Cover *notionapi.File `notion:",cover"`
		}{}), &tagInfo{cover: true}},
	}

	for _, c := range cases {
		got := parseTag(c.Type.Field(0))
		if !reflect.DeepEqual(got, c.Want) {
			t.Errorf("got:%v want: %v", got, c.Want)
		}
	}
}
