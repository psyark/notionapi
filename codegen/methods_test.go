package codegen

import (
	"testing"
)

func TestDatabaseMethods(t *testing.T) {
	type Method struct {
		DocURL  string
		Returns string
	}
	methods := []Method{
		{"https://developers.notion.com/reference/post-database-query", "Pagination"},
		{"https://developers.notion.com/reference/create-a-database", "Database"},
		{"https://developers.notion.com/reference/update-a-database", "Database"},
		{"https://developers.notion.com/reference/retrieve-a-database", "Database"},
		{"https://developers.notion.com/reference/retrieve-a-page", "Page"},
		{"https://developers.notion.com/reference/post-page", "Page"},
		{"https://developers.notion.com/reference/patch-page", "Page"},
		// {"https://developers.notion.com/reference/retrieve-a-page-property", "?"},
		// {"https://developers.notion.com/reference/retrieve-a-block", "?"},
	}

	builder := NewBuilder()

	for _, method := range methods {
		ssrProps, err := ParseMethod(method.DocURL)
		if err != nil {
			t.Fatal(err)
		}

		builder.Add(MethodCoder{DocURL: method.DocURL, Props: *ssrProps, Returns: method.Returns})
	}

	if err := builder.Build("../client.methods.go"); err != nil {
		t.Fatal(err)
	}
}
