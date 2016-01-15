package parsers

import (
	"bufio"
	. "github.com/ruxton/mix_cover_builder/data"
	"github.com/ruxton/term"
	"io"
	"os"
	"strings"
)

func ParseVirtualDJTracklist(bufReader *bufio.Reader) []Track {
	var list []Track

	for line, _, err := bufReader.ReadLine(); err != io.EOF; line, _, err = bufReader.ReadLine() {
		data := strings.SplitN(string(line), " : ", 2)
		trackdata := strings.SplitN(data[1], " - ", 2)
		if len(trackdata) != 2 {
			term.OutputError("Error parsing track " + string(data[1]))
			term.OutputMessage("Please enter an artist for this track: ")
			artist, err := term.STD_IN.ReadString('\n')
			if err != nil {
				term.OutputError("Incorrect artist entry.")
				os.Exit(2)
			}
			term.OutputMessage("Please enter a name for this track: ")
			track, err := term.STD_IN.ReadString('\n')
			if err != nil {
				term.OutputError("Incorrect track name entry.")
				os.Exit(2)
			}

			trackdata = []string{artist, track}
		}

		thistrack := new(Track)
		thistrack.Artist = trackdata[0]
		thistrack.Song = trackdata[1]

		list = append(list, *thistrack)
	}

	return list
}
