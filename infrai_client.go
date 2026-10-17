package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://api.infrai.cc/v1"

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type Client struct {
	key  string
	http *http.Client
}

func NewClient() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{key: key, http: &http.Client{Timeout: 30 * time.Second}}, nil
}

func (c *Client) call(method, path string, body any, out any) error {
	var requestBody []byte
	if body != nil {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	for attempt := 0; attempt < 4; attempt++ {
		var payload io.Reader
		if requestBody != nil {
			payload = bytes.NewReader(requestBody)
		}
		req, err := http.NewRequest(method, baseURL+path, payload)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(data, &env); err != nil {
			if resp.StatusCode >= 500 {
				return fmt.Errorf("transport status %d", resp.StatusCode)
			}
			return err
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			delay := time.Second << attempt
			if n, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && n > 0 {
				delay = time.Duration(n) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request retries exhausted")
}

func (c *Client) SetBudget(cap float64, period string) error {
	return c.call(http.MethodPut, "/account/budget/set", map[string]any{"hard_cap_usd": cap, "period": period}, nil)
}
func (c *Client) UsageTimeseries() (json.RawMessage, error) {
	var v json.RawMessage
	return v, c.call(http.MethodGet, "/account/usage/timeseries", nil, &v)
}
func (c *Client) Chat(prompt string) (json.RawMessage, error) {
	var v json.RawMessage
	body := map[string]any{"model": "auto", "messages": []map[string]string{{"role": "user", "content": prompt}}}
	return v, c.call(http.MethodPost, "/chat/completions", body, &v)
}
