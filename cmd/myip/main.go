package main

import (
	"fmt"
	"os"

	"github.com/awnzl/myip/internal/ipfinder"
)

var providers = []string{
	"https://icanhazip.com",
	"https://ifconfig.co",
	"https://ipecho.net/plain",
	"https://ifconfig.me",
	"https://checkip.amazonaws.com",
	//"https://whatismyip.com",
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		os.Exit(1)
	}

	finder := ipfinder.New(providers)
	resp, err := finder.FindIp(cfg.AllProviders, cfg.Timeout)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	fmt.Println(resp)
}
