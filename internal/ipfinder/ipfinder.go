package ipfinder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

const timedOut = "Timed Out"

type Finder struct {
	providers []string
}

type Response struct {
	Provider string
	IP       string
}

type ResponseWithError struct {
	Response Response
	Error	 error
}

func New(aProviders []string) *Finder {
	return &Finder{providers: aProviders}
}

func (f *Finder) FindIp(ctx context.Context, useAllProviders bool) ([]Response, error) {
	c := f.requestProviders(ctx)

	if useAllProviders {
		return allResponses(c)
	}

	resp := <-c
	if resp.Error != nil {
		return nil, resp.Error
	}

	return []Response{resp.Response}, nil
}

func allResponses(c <-chan ResponseWithError) ([]Response, error) {
	var responses []Response

	for resp := range c {
		if resp.Error != nil {
			return nil, resp.Error
		}

		responses = append(responses, resp.Response)
	}

	return responses, nil
}

func (f *Finder) requestProviders(ctx context.Context) <-chan ResponseWithError {
	out := make(chan ResponseWithError)

	go func() {
		wg := sync.WaitGroup{}

		for _, url := range f.providers {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				r, err := requestIP(ctx, url)
				out <- ResponseWithError{Response: r, Error: err}
			}(url)
		}

		go func() {
			wg.Wait()
			close(out)
		}()
	}()

	return out
}

func requestIP(ctx context.Context, url string) (Response, error) {
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Response{}, err
	}

	//req.Header.Set("User-Agent", "")
	//req.Header.Set("Accept", `*/*`)

	resp, err := client.Do(req)
	if errors.Is(err, context.DeadlineExceeded) {
		return Response{Provider: url, IP: timedOut}, nil
	}
	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("ipfinder: response isn't useful, response code: %v", resp.StatusCode)
	}

	ip, err := extractIP(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{Provider: url, IP: ip}, nil
}

func extractIP(src io.Reader) (string, error) {
	bytes, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}
