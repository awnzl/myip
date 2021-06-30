package main

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
)

var ErrParseArgs = errors.New("cfg: incorrect argument")

const usageInfo = `Usage of MyIP: [-a|--all-providers][-t|--timeout=3]
  -a, --all-providers
    	Use all providers to obtain IP
  -t, --timeout=seconds
    	Timeout in seconds: --timeout=3 (default 5)`

type Config struct {
	AllProviders bool
	Timeout      float64
}

func parseConfig() (Config, error) {
	var allProviders bool
	var timeout float64

	pflag.Float64VarP(&timeout, "timeout", "t", 5, "")
	pflag.BoolVarP(&allProviders, "all-providers", "a", false, "")

	pflag.Usage = usage
	pflag.Parse()

	if timeout < 0.1 {
		fmt.Printf("invalid value %v for flag -t\n", timeout)
		usage()

		return Config{}, ErrParseArgs
	}

	return Config{
		AllProviders: allProviders,
		Timeout:      timeout,
	}, nil
}

func usage() {
	fmt.Println(usageInfo)
}
