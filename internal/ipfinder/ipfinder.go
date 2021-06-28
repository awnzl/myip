package ipfinder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
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

func (f *Finder) FindIp(useAllProviders bool, timeout int) ([]Response, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var responses []Response

	c := f.requestProviders(timeoutCtx)
	r := <-c
	if r.Error != nil {
		return nil, r.Error
	}

	responses = append(responses, r.Response)

	if useAllProviders {
		for r := range c {
			if r.Error != nil {
				return nil, r.Error
			}
			responses = append(responses, r.Response)
		}
	}

	return responses, nil
}

func (f *Finder) requestProviders(timeoutCtx context.Context) <-chan ResponseWithError {
	out := make(chan ResponseWithError)

	go func() {
		wg := sync.WaitGroup{}

		for _, url := range f.providers {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				r, err := requestIP(timeoutCtx, url)
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

	//todo: remove test block
	// test block start
	done := make(chan struct{})
	go func() {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(10)
		fmt.Println("sleep time", n, url)
		time.Sleep(time.Duration(n) * time.Second)

		close(done)
	}()

	select {
	case <-done:
		break
	case <-ctx.Done():
		break
	}
	// test block end

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
		// TODO: is this correct error?
		return Response{}, http.ErrBodyNotAllowed
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
