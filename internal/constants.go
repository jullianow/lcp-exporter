package internal

import (
	"fmt"
)

const (
	prefix = "lcp_api"
)

func Name(c string) func(string) string {
	return func(s string) string {
		return fmt.Sprintf("%s_%s_%s", prefix, c, s)
	}
}
