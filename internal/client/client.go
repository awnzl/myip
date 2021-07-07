package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const timedOut = "Timed Out"

var ErrorIncorrectResponse = errors.New("client: response isn't useful")

type Response struct {
	Provider string
	IP       string
}

type IPClient interface {
	Get(ctx context.Context) (Response, error)
}

type TextClient struct {
	url string
}

func NewTextClient(u string) *TextClient {
	return &TextClient{url: u}
}

func (c *TextClient) Get(ctx context.Context) (Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return Response{}, err
	}

	bytes, err := responseData(req, c.url)
	if errors.Is(err, context.DeadlineExceeded) {
		return Response{Provider: c.url, IP: timedOut}, nil
	}
	if err != nil {
		return Response{}, err
	}

	return Response{Provider: c.url, IP: strings.TrimSpace(string(bytes))}, nil
}

type JSONClient struct {
	url string
}

type JSONResponse struct {
	IP string `json:"ip"`
}

func NewJSONClient(u string) *JSONClient {
	return &JSONClient{url: u}
}

func (c *JSONClient) Get(ctx context.Context) (Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Accept", "application/json")

	bytes, err := responseData(req, c.url)
	if errors.Is(err, context.DeadlineExceeded) {
		return Response{Provider: c.url, IP: timedOut}, nil
	}
	if err != nil {
		return Response{}, err
	}

	var resp JSONResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return Response{}, err
	}

	return Response{Provider: c.url, IP: resp.IP}, nil
}

func responseData(req *http.Request, url string) ([]byte, error) {
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w, response code: %v", ErrorIncorrectResponse, resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
