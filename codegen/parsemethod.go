package codegen

import (
	"encoding/json"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func ParseMethod(url string) (*SSRProps, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	ssrPropsBytes := []byte(doc.Find(`#ssr-props`).AttrOr("data-initial-props", ""))
	ssrProps := SSRProps{}
	if err := json.Unmarshal(ssrPropsBytes, &ssrProps); err != nil {
		return nil, err
	}
	// ssrPropsBytes2, err := json.MarshalIndent(ssrProps, "", "  ")
	// if err != nil {
	// 	return nil, err
	// }

	// diff, err := gojsondiff.New().Compare(ssrPropsBytes, ssrPropsBytes2)
	// if err != nil {
	// 	return nil, err
	// }

	// if diff.Modified() {
	// 	for _, delta := range diff.Deltas() {
	// 		printDelta(delta)
	// 	}
	// 	return nil, fmt.Errorf("validation failed")
	// }

	return &ssrProps, nil
}

// func printDelta(delta gojsondiff.Delta, path ...string) {
// 	switch delta := delta.(type) {
// 	case *gojsondiff.Added:
// 		fmt.Printf("added: %v\n", delta.Position)
// 	case *gojsondiff.Deleted:
// 		value := fmt.Sprintf("%v", delta.Value)
// 		if len(value) > 50 {
// 			value = value[0:50-3] + "..."
// 		}
// 		fmt.Printf("%v %v `json:\"%v\"` // %v\n", getName(delta.Position.String()), reflect.TypeOf(delta.Value), delta.Position, value)
// 	case *gojsondiff.Object:
// 		path = append(path, delta.Position.String())
// 		fmt.Printf("%v\n", path)
// 		for _, delta := range delta.Deltas {
// 			printDelta(delta, path...)
// 		}
// 	case *gojsondiff.Modified:
// 		fmt.Printf("modified: %v %v->%v\n", delta.Position, delta.OldValue, delta.NewValue)
// 	case *gojsondiff.Array:
// 		path = append(path, delta.Position.String())
// 		fmt.Printf("%v\n", path)
// 		for _, delta := range delta.Deltas {
// 			printDelta(delta, path...)
// 		}
// 	default:
// 		fmt.Println(reflect.TypeOf(delta))
// 		panic(reflect.TypeOf(delta))
// 	}
// }
