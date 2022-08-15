package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yudai/gojsondiff"
)

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
			return fmt.Errorf("bat status: %v, %v", res.Status, string(resBody))
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
			return debugErr{resBody, remarshaled, diff}
		}
	}

	return nil
}
