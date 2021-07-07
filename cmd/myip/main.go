package main

import (
	"context"
	"os"

	"github.com/awnzl/myip/internal/client"
	"github.com/awnzl/myip/internal/ipfinder"
	"github.com/awnzl/myip/internal/writer"
)

var textProviders = []string{
	"https://icanhazip.com",
	"https://ifconfig.co",
	"https://ipecho.net/plain",
	"https://ifconfig.me",
	"https://checkip.amazonaws.com",
}

var jsonProviders = []string{
	"https://ifconfig.co",
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		os.Exit(1)
	}

	var clients []client.IPClient
	for _, url := range textProviders {
		clients = append(clients, client.NewTextClient(url))
	}
	for _, url := range jsonProviders {
		clients = append(clients, client.NewJSONClient(url))
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	finder := ipfinder.New(clients)
	resp, err := finder.FindIp(timeoutCtx, cfg.AllProviders)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	writer.New().Write(resp)
}
