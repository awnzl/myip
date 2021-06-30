package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/awnzl/myip/internal/client"
	"github.com/awnzl/myip/internal/ipfinder"
)

var textProviders = []string{
	"https://icanhazip.com",
	"https://ifconfig.co",
	"https://ipecho.net/plain",
	"https://ifconfig.me",
	"https://checkip.amazonaws.com",
	//"https://whatismyip.com",
}

var jsonProviders = []string{
	"https://ifconfig.co",
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		os.Exit(1)
	}

	var textClients []client.IPClient
	for _, url := range textProviders {
		textClients = append(textClients, client.NewTextClient(url))
	}

	finder := ipfinder.New(textClients)

	t, err := time.ParseDuration(fmt.Sprintf("%vs", cfg.Timeout))
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	resp, err := finder.FindIp(timeoutCtx, cfg.AllProviders)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	fmt.Println(resp)
}
