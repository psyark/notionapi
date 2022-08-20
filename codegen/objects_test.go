package codegen

import (
	"testing"
)

func TestObjects(t *testing.T) {
	subTests := map[string]func() error{
		"Block":         BuildBlock,
		"Database":      BuildDatabase,
		"FileOrEmoji":   BuildFileOrEmoji,
		"Page":          BuildPage,
		"Pagination":    BuildPagination,
		"Parent":        BuildParent,
		"Property":      BuildProperty,
		"PropertyItem":  BuildPropertyItem,
		"PropertyValue": BuildPropertyValue,
		"RichText":      BuildRichText,
		"User":          BuildUser,
	}

	for name, subTest := range subTests {
		t.Run(name, func(t *testing.T) {
			if err := subTest(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
