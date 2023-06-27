package nethttpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type netHttpClient struct {
	client *http.Client
}

func NewNetHttpClient() Http {
	return &netHttpClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *netHttpClient) Dispatch(
	ctx context.Context,
	response any,
	request any,
	customHeaders map[string]string,
) (*Response, error) {
	if req, ok := request.(Request); ok {
		var reqBody io.Reader

		switch req.Method {
		case http.MethodPost, http.MethodPut, http.MethodDelete:
			reqBody = bytes.NewBuffer([]byte(req.Body))
		default:
			reqBody = nil
		}

		req, err := http.NewRequestWithContext(ctx, req.Method, req.Url, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP request: %w", err)
		}

		if len(customHeaders) != 0 {
			for key, value := range customHeaders {
				req.Header.Add(key, value)
			}
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("failed to close response body: %v", err)
			}
		}()

		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode >= http.StatusBadRequest {
			return nil, fmt.Errorf("server returned %d status code: %s", resp.StatusCode, string(resBody))
		}

		if response != nil {
			if err := json.Unmarshal(resBody, &response); err != nil {
				return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
			}
		}

		return &Response{
			Body:   string(resBody),
			Status: resp.StatusCode,
		}, nil
	}

	return nil, errors.New("request model is incorrect")
}
