package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const timedOut = "Timed Out"

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
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return Response{}, err
	}

	resp, err := client.Do(req)
	if errors.Is(err, context.DeadlineExceeded) {
		return Response{Provider: c.url, IP: timedOut}, nil
	}
	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("client: response isn't useful, response code: %v", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{Provider: c.url, IP: strings.TrimSpace(string(bytes))}, nil
}
