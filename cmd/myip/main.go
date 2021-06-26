package main

import (
	"fmt"
	"os"

	"github.com/awnzl/myip/internal/ipfinder"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		os.Exit(1)
	}

	finder := ipfinder.New()
	fmt.Println(finder.FindIp(cfg.AllProviders, cfg.Timeout))

}
