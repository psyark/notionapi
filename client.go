package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const APIVersion = "2022-06-28"

type EventListener interface {
	OnReadBody([]byte) error
	OnUnmarshal(interface{}) error
}

func NewClient(accessToken string) *Client {
	return &Client{accessToken: accessToken}
}

func (c *Client) SetListener(listener EventListener) *Client {
	c.listener = listener
	return c
}

type Client struct {
	accessToken string
	listener    EventListener
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

	switch method {
	case http.MethodPost, http.MethodPatch:
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if c.listener != nil {
		if err := c.listener.OnReadBody(resBody); err != nil {
			return err
		}
	}

	if res.StatusCode != http.StatusOK {
		errBody := Error{}
		if err := json.Unmarshal(resBody, &errBody); err != nil {
			return fmt.Errorf("bad status: %v, %v", res.Status, string(resBody))
		} else {
			return errBody
		}
	}

	if err := json.Unmarshal(resBody, &result); err != nil {
		return err
	}

	if c.listener != nil {
		if err := c.listener.OnUnmarshal(result); err != nil {
			return err
		}
	}

	return nil
}
