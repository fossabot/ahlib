// +build !windows

package xcolor

import (
	"io"
)

func checkTerminal(w io.Writer) bool {
	return false
}
