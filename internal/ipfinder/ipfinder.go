package ipfinder

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const timedOut = "Timed Out"

type Finder struct {
	providers []string
}

type Provider struct {
	URL      string
	Response string
}

func New() *Finder {
	finder := &Finder{
		providers: []string{
			"https://icanhazip.com",
			"https://ifconfig.co",
			"https://ipecho.net/plain",
			"https://ifconfig.me",
			"https://checkip.amazonaws.com",
			//"https://whatismyip.com",
		},
	}

	return finder
}

func (f *Finder) FindIp(useAllProviders bool, timeout int) []Provider {
	if useAllProviders {
		return f.allProvidersResponse(timeout)
	}

	return f.anyProviderResponse(timeout)
}

func (f *Finder) allProvidersResponse(timeout int) []Provider {
	responseChan := make(chan Provider)

	go func() {
		wg := new(sync.WaitGroup)
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		defer close(responseChan)

		for _, url := range f.providers {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				requestIP(timeoutCtx, responseChan, url)
			}(url)
		}

		wg.Wait()
	}()

	var results []Provider

	for resp := range responseChan {
		results = append(results, resp)
	}

	return results
}

func (f *Finder) anyProviderResponse(timeout int) []Provider {
	responseChan := make(chan Provider)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	defer close(responseChan)

	for _, url := range f.providers {
		go requestIP(timeoutCtx, responseChan, url)
	}

	for {
		select {
		case <-timeoutCtx.Done():
			return []Provider{{URL: "All Providers", Response: timedOut}}
		case v := <-responseChan:
			return []Provider{v}
		}
	}
}

func requestIP(ctx context.Context, out chan<- Provider, url string) {
	select {
	case <-ctx.Done():
		out <- Provider{URL: url, Response: timedOut}
		return
	default:
		client := http.DefaultClient
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return
		}

		//todo: remove
		//rand.Seed(time.Now().UnixNano())
		//n := rand.Intn(10)
		//fmt.Println("sleep time", n, url)
		//time.Sleep(time.Duration(n) * time.Second)

		//req.Header.Set("User-Agent", "")
		//req.Header.Set("Accept", `*/*`)

		resp, err := client.Do(req)
		if errors.Is(err, context.DeadlineExceeded) {
			out <- Provider{URL: url, Response: timedOut}
			return
		}
		if err != nil {
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}

		out <- Provider{URL: url, Response: extractIP(resp.Body)}
	}
}

func extractIP(src io.Reader) string {
	bytes, _ := io.ReadAll(src)

	return strings.TrimSpace(string(bytes))
}
