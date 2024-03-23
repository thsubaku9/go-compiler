package repl

import (
	"bufio"
	"io"
)

func Start(in io.Reader, out io.Writer) {
	_ = bufio.NewScanner(in)

}
