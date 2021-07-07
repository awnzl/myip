package ipfinder

import (
	"context"
	"sync"

	"github.com/awnzl/myip/internal/client"
)

type Finder struct {
	providers []client.IPClient
}

type ResponseWithError struct {
	Response client.Response
	Error	 error
}

func New(aProviders []client.IPClient) *Finder {
	return &Finder{providers: aProviders}
}

func (f *Finder) FindIp(ctx context.Context, useAllProviders bool) ([]client.Response, error) {
	c := f.requestProviders(ctx)

	if useAllProviders {
		return allResponses(c)
	}

	resp := <-c
	if resp.Error != nil {
		return nil, resp.Error
	}

	return []client.Response{resp.Response}, nil
}

func allResponses(c <-chan ResponseWithError) ([]client.Response, error) {
	var responses []client.Response

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

		for _, c := range f.providers {
			wg.Add(1)
			go func(c client.IPClient) {
				defer wg.Done()
				r, err := c.Get(ctx)
				out <- ResponseWithError{Response: r, Error: err}
			}(c)
		}

		go func() {
			wg.Wait()
			close(out)
		}()
	}()

	return out
}
