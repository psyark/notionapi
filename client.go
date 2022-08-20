package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/yudai/gojsondiff"
)

const APIVersion = "2022-06-28"

func NewClient(accessToken string) *Client {
	return &Client{accessToken: accessToken}
}

func NewDebugClient(accessToken string) *Client {
	return &Client{accessToken: accessToken, debug: true}
}

type Client struct {
	accessToken string
	debug       bool
}

func (c *Client) call(ctx context.Context, method string, path string, body interface{}, result interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, "https://api.notion.com"+path, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.accessToken)
	req.Header.Add("Notion-Version", APIVersion)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		e := Error{}
		if err := json.Unmarshal(resBody, &e); err != nil {
			return fmt.Errorf("bad status: %v, %v", res.Status, string(resBody))
		} else {
			return fmt.Errorf("%v (%v)", e.Code, e.Message)
		}
	}

	if err := json.Unmarshal(resBody, &result); err != nil {
		return err
	}

	if c.debug {
		remarshaled, _ := json.MarshalIndent(result, "", "  ")

		differ := gojsondiff.New()
		diff, err := differ.Compare(resBody, remarshaled)
		if err != nil {
			return err
		}

		if diff.Modified() {
			return validationError{resBody, remarshaled, diff}
		}
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
