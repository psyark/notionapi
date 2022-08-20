package notionapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/yudai/gojsondiff"
)

type TestListener struct {
	body []byte
}

func (h *TestListener) OnReadBody(data []byte) error {
	h.body = data
	return nil
}
func (h *TestListener) OnUnmarshal(data interface{}) error {
	remarshaled, _ := json.MarshalIndent(data, "", "  ")

	diff, err := gojsondiff.New().Compare(h.body, remarshaled)
	if err != nil {
		return err
	}

	if diff.Modified() {
		return validationError{h.body, remarshaled, diff}
	}
	return nil
}

type diffMap map[string]interface{}

func (dm diffMap) add(delta gojsondiff.Delta, prefix string) {
	switch delta := delta.(type) {
	case *gojsondiff.Added:
		dm[prefix+delta.Position.String()] = struct {
			Value interface{} `json:"ADDED"`
		}{delta.Value}
	case *gojsondiff.Deleted:
		dm[prefix+delta.Position.String()] = struct {
			Value interface{} `json:"DELETED"`
		}{delta.Value}
	case *gojsondiff.Modified:
		dm[prefix+delta.Position.String()] = struct {
			OldValue interface{} `json:"MODIFIED_FROM"`
			NewValue interface{} `json:"MODIFIED_TO"`
		}{delta.OldValue, delta.NewValue}
	case *gojsondiff.Array:
		for _, d := range delta.Deltas {
			dm.add(d, prefix+delta.Position.String()+".")
		}
	case *gojsondiff.Object:
		for _, d := range delta.Deltas {
			dm.add(d, prefix+delta.Position.String()+".")
		}
	default:
		panic(reflect.TypeOf(delta))
	}
}

type validationError struct {
	source      []byte
	remarshaled []byte
	diff        gojsondiff.Diff
}

func (e validationError) Error() string {
	dm := diffMap{}

	// res, _ := formatter.NewDeltaFormatter().Format(e.diff)
	for _, delta := range e.diff.Deltas() {
		dm.add(delta, "")
	}

	diffStr, _ := json.MarshalIndent(dm, "", "  ")
	buffer := bytes.NewBuffer(nil)
	json.Indent(buffer, e.source, "", "  ")
	return fmt.Sprintf("validation failed. diff: %v\nwant: %v\ngot: %v", string(diffStr), buffer.String(), string(e.remarshaled))
}
