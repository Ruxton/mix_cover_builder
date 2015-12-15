package parsers

import (
  "github.com/ruxton/term"
  "bufio"
  "os"
)

func ParseSeratoTracklist(bufReader *bufio.Reader) {
  term.OutputError("Serato tracklist parsing unsupported")
  os.Exit(2)
}
