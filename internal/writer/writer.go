package writer

import (
	"fmt"
	"strings"

	"github.com/awnzl/myip/internal/client"
)

type Writer struct {
}

func New() *Writer {
	return &Writer{}
}

func (p *Writer) Write(providers []client.Response) {
	maxLen := 0
	for _, p := range providers {
		if l := len(p.Provider); l > maxLen {
			maxLen = l
		}
	}

	fmt.Println("Provider", strings.Repeat(" ", maxLen - len("Provider")), "IP")

	for _, resp := range providers {
		fmt.Println(resp.Provider, strings.Repeat(" ", maxLen - len(resp.Provider)), resp.IP)
	}
}
