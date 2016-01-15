package parsers

import (
	"bufio"
	"github.com/ruxton/term"
	"os"
)

func ParseSeratoTracklist(bufReader *bufio.Reader) {
	term.OutputError("Serato tracklist parsing unsupported")
	os.Exit(2)
}
