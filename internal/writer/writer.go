package writer

import (
	"fmt"
	"io"
	"strings"

	"github.com/awnzl/myip/internal/client"
)

type Writer struct {
	writer   io.Writer
}

func New(aWriter io.Writer) *Writer {
	return &Writer{
		writer: aWriter,
	}
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
