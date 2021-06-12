package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var ErrParseArgs = errors.New("cfg: incorrect argument")

const usageInfo = `Usage of MyIP: [-a|--all-providers][-t|--timeout=3]
  -a, --all-providers
    	Use all providers to obtain IP
  -t, --timeout
    	Timeout in seconds: --timeout=3 (default 5)`

type Config struct {
	AllProviders bool
	Timeout      int
}

func parseConfig() Config {
	var allProviders, shortAllProviders bool
	var timeout, shortTimeout int

	flag.IntVar(&timeout, "timeout", 5, "Timeout in seconds: --timeout=3")
	flag.BoolVar(&allProviders, "all-providers", false, "Bool flag: --all-providers")
	flag.IntVar(&shortTimeout, "t", 5, "Timeout in seconds: --timeout=3")
	flag.BoolVar(&shortAllProviders, "a", false, "Bool flag: --all-providers")

	flag.Usage = usage
	flag.Parse()

	if shortTimeout != 5 {
		timeout = shortTimeout
	}

	if timeout < 1 {
		fmt.Printf("invalid value %v for flag -t\n", timeout)
		usage()
		os.Exit(0)
	}

	if shortAllProviders {
		allProviders = shortAllProviders
	}

	return Config{
		AllProviders: allProviders,
		Timeout:      timeout,
	}
}

func usage() {
	fmt.Println(usageInfo)
}
