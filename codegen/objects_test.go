package codegen

import (
	"testing"
)

func TestObjects(t *testing.T) {
	subTests := map[string]func() error{
		"Filter":        BuildFilter,
		"Property":      BuildProperty,
		"PropertyItem":  BuildPropertyItem,
		"PropertyValue": BuildPropertyValue,
		"Sort":          BuildSort,
	}

	for name, subTest := range subTests {
		subTest := subTest
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if err := subTest(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
