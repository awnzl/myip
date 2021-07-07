package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
)

var ErrIncorrectArgument = errors.New("cfg: incorrect argument")

const usageInfo = `Usage of MyIP: [-a|--all-providers][-t|--timeout=3]
  -a, --all-providers
    	Use all providers to obtain IP
  -t, --timeout=seconds
    	Timeout in seconds: --timeout=3 (default 5)`

type Config struct {
	AllProviders bool
	Timeout      time.Duration
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

		return Config{}, ErrIncorrectArgument
	}

	t, err := time.ParseDuration(fmt.Sprintf("%vs", timeout))
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return Config{}, err
	}

	return Config{
		AllProviders: allProviders,
		Timeout:      t,
	}, nil
}

func usage() {
	fmt.Println(usageInfo)
}
