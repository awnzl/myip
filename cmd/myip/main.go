package main

import (
	"fmt"

	"github.com/awnzl/myip/internal/ipfinder"
)

func main() {
	cfg := parseConfig()

	finder := ipfinder.New()
	fmt.Println(finder.FindIp(cfg.AllProviders, cfg.Timeout))
}
