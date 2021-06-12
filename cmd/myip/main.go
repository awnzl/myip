package main

import (
	"fmt"
)

func main() {
	cfg := parseConfig()

	fmt.Println(cfg.AllProviders, cfg.Timeout)
}
