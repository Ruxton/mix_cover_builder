package main

import (
	"bufio"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/disintegration/imaging"
  "github.com/ruxton/tracklist_parsers/parsers"
  . "github.com/ruxton/tracklist_parsers/data"
	. "github.com/ruxton/mix_cover_builder/confirm"
	"github.com/ruxton/mix_cover_builder/google"
	"github.com/ruxton/mix_cover_builder/itunes"
	"github.com/ruxton/mix_cover_builder/versions"
	"github.com/ruxton/term"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	flag "launchpad.net/gnuflag"
	"net/http"
	"os"
	"sort"
	"strings"
	"unicode"
)

var DATABASE_FILENAME = "cover_urls.boltdb"
var DATABASE_BUCKET = []byte("CoverImages")
var DEFAULT_OUTPUT_FILENAME = "output.jpg"

var aboutFlag = flag.Bool("about", false, "About the application")
var trackListFlag = flag.String("tracklist", "", "A file containing a tracklist")
var trackListTypeFlag = flag.String("tracklist-type", "basic", "A file containing a tracklist (basic,virtualdj,serato)")
var outputFileFlag = flag.String("output", "", "The file to output the generate cover to")
var overlayFlag = flag.String("overlay", "", "A file to overlay on top of the generated image")

func main() {

	flag.Parse(true)
	showWelcomeMessage()
	if *aboutFlag == true {
		showAboutMessage()
		os.Exit(0)
	}

	var tracklist []Track
	tracklistFileName := ""

	if *trackListFlag != "" {
		tracklistFileName = *trackListFlag
		tracklist = parseTracklist(trackListFlag)
	} else {
		if len(os.Args) > 1 {
			if os.Args[1][0:2] != "--" {
				tracklistFileName = string(os.Args[1])
				tracklist = parseTracklist(&tracklistFileName)
			} else {
				term.OutputError("Please provide a tracklist to parse")
				os.Exit(1)
			}
		} else {
			term.OutputError("Please provide a tracklist to parse")
			os.Exit(1)
		}
	}

	if len(tracklist) < 9 {
		term.OutputError("Not enough tracks in tracklist")
		os.Exit(1)
	}

	overlay_filename := ""
	if *overlayFlag != "" {
		overlay_filename = *overlayFlag
	}

	output_filename := strings.Replace(tracklistFileName, ".txt", ".jpg", -1)

	if *outputFileFlag != "" {
		output_filename = *outputFileFlag
	} else {
		if len(os.Args) > 2 {
			if os.Args[2][0:2] != "--" && os.Args[2] != "" {
				output_filename = string(os.Args[2])
			}
		}
	}

	artists := countArtists(tracklist)
	sortedArtists := sortMapByValue(artists)

	selection := sortedArtists[:9]

	db, err := bolt.Open(DATABASE_FILENAME, 0600, nil)
	if err != nil {
		term.OutputError(err.Error())
		os.Exit(2)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(DATABASE_BUCKET)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	for i, _ := range selection {
		j := 0
		for !GetCover(db, i, &selection[i]) {
			replacement := sortedArtists[9+i+j]
			term.OutputError(fmt.Sprintf("Replacing with %s", replacement.Artist))
			selection[i] = replacement
			j++
		}
	}

	createImage(selection, output_filename, overlay_filename)
}

func FetchImage(artistTrack *ArtistTrack) {
	term.OutputMessage(term.Yellow + "Fetching image.." + term.Reset)
	term.OutputMessage(term.Yellow + "." + term.Reset)
	client := http.Client{}

	req, err := http.NewRequest("GET", artistTrack.Tracks[0].Cover, nil)
	if err != nil {
		fmt.Printf("Error - %s", err.Error())
	}
	resp, doError := client.Do(req)
	if doError != nil {
		term.OutputError("Error - " + doError.Error())
	}

	outputImage, _, err := image.Decode(resp.Body)
	artistTrack.Tracks[0].CoverImage = outputImage

	term.OutputMessage(term.Yellow + term.Bold + "DONE!\n" + term.Reset)
}

func GetCover(db *bolt.DB, index int, artistTrack *ArtistTrack) bool {
	artist := artistTrack.Artist
	track := artistTrack.Tracks[0].Song

	artistTrack.Tracks[0].Cover = getCoverUrl(db, artist, track)

	if artistTrack.Tracks[0].Cover != "" {
		FetchImage(artistTrack)
		term.OutputImageUrl(artistTrack.Tracks[0].Cover, "Image.jpg")
		term.OutputMessage("Is this the correct image? (y/n) ")
		if AskForConfirmation() {
			return true
		} else {
			artistTrack.Tracks[0].Cover = getCoverUrlFromInput()
			if artistTrack.Tracks[0].Cover == "" {
				return false
			} else {
				FetchImage(artistTrack)
				term.OutputImageUrl(artistTrack.Tracks[0].Cover, "Image.jpg")
				term.OutputMessage("Is this the correct image? (y/n) ")
				if AskForConfirmation() {
					InsertData(db, artist, track, artistTrack.Tracks[0].Cover)
					return true
				} else {
					return false
				}
			}
		}
	} else {
		return false
	}
}

func GetData(db *bolt.DB, artist string, track string) string {
	dbCover := ""
	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(DATABASE_BUCKET)

		val := bucket.Get([]byte(artist + " - " + track))
		dbCover = string(val)
		return nil
	})

	return dbCover
}

func InsertData(db *bolt.DB, artist string, track string, cover string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(DATABASE_BUCKET)
		err := b.Put([]byte(artist+" - "+track), []byte(cover))
		return err
	})
}

func createImage(selection ArtistTrackList, output_filename string, overlay string) {
	if len(selection) != 9 {
		term.OutputError("Unable to find enough images to make cover")
	} else {
		term.OutputMessage(term.Yellow + "Generating image.." + term.Reset)
		xt := 3
		yt := 3
		width := xt * 160
		height := yt * 160

		dst := image.NewRGBA(image.Rect(0, 0, width, height))
		count := 0
		for y := 0; y < yt; y++ {
			for x := 0; x < xt; x++ {
				term.OutputMessage(term.Yellow + "." + term.Reset)
				src := selection[count].Tracks[0].CoverImage
				draw.Draw(dst, image.Rect(x*160, y*160, (x+1)*160, (y+1)*160), imaging.Thumbnail(src, 160, 160, imaging.CatmullRom), image.Pt(0, 0), draw.Src)
				count = count + 1
			}
		}

		if overlay != "" {
			term.OutputMessage(term.Yellow + "." + term.Reset)
			maskImage(dst, overlay)
			term.OutputMessage(term.Yellow + "." + term.Reset)
		}

		f, err := os.Create(output_filename)
		if err != nil {
			term.OutputError("\ncant save picture: " + err.Error())
		}
		defer f.Close()

		term.OutputMessage(term.Yellow + "." + term.Reset)
		jpeg.Encode(f, dst, nil)
		term.OutputMessage(term.Yellow + term.Bold + "DONE!\n" + term.Reset)
	}
}

func maskImage(dst draw.Image, filename string) {
	infile, err := os.Open(filename)
	defer infile.Close()

	if err != nil {
		term.OutputError(fmt.Sprintf("Unable to open overlay image - %s", err.Error()))
	} else {
		overlay, _, err := image.Decode(infile)
		if err != nil {
			term.OutputError(fmt.Sprintf("Unable to decode image file - %s", err.Error()))
		} else {
			mask := image.NewUniform(color.Alpha{128})
			draw.DrawMask(dst, dst.Bounds(), imaging.Thumbnail(overlay, 480, 480, imaging.CatmullRom), image.Pt(0, 0), mask, image.Pt(0, 0), draw.Over)
		}
	}
}

func getCoverUrl(db *bolt.DB, artist string, track string) string {
	term.OutputMessage(fmt.Sprintf(term.Green+"Searching for %s - %s.."+term.Reset, artist, track))

	cover_url := GetData(db, artist, track)
	if cover_url != "" {
		term.OutputMessage(term.Green + term.Bold + "DONE!\n" + term.Reset)
	} else {
		term.OutputMessage(term.Green + "." + term.Reset)
		cover_url = itunes.GetCoverFor(artist, track)
		if cover_url == "" {
			term.OutputMessage(term.Green + "." + term.Reset)
			cover_url = google.GetCoverFor(artist, track)
		}
		if cover_url == "" {
			term.OutputMessage(term.Green + ".\n" + term.Reset)
			term.OutputError(fmt.Sprintf("Unable to find a cover for %s - %s", artist, track))
			cover_url = getCoverUrlFromInput()
		}
		if cover_url != "" {
			term.OutputMessage(term.Green + term.Bold + "DONE!\n" + term.Reset)
			InsertData(db, artist, track, cover_url)
		}
	}

	return cover_url
}

func getCoverUrlFromInput() string {
	term.OutputMessage("Enter a URL the cover image: ")
	cover_url, err := term.STD_IN.ReadString('\n')
	if err != nil {
		term.OutputError("Error accepting input.")
		os.Exit(2)
	}
	return strings.TrimRightFunc(cover_url, unicode.IsSpace)
}

func showWelcomeMessage() {
	term.OutputMessage(term.Green + "Cover Builder v" + versions.VERSION + term.Reset + "\n\n")
}

func showAboutMessage() {
	term.OutputMessage(fmt.Sprintf("Build Number: %s\n", versions.MINVERSION))
	term.OutputMessage("Created by: Greg Tangey (http://ignite.digitalignition.net/)\n")
	term.OutputMessage("Website: http://www.rhythmandpoetry.net/\n")
}

func parseTracklist(tracklist *string) []Track {

	fin, err := os.Open(*tracklist)
	if err != nil {
		term.OutputError(fmt.Sprintf("The file %s does not exist!\n", *tracklist))
		os.Exit(1)
	}
	defer fin.Close()

	bufReader := bufio.NewReader(fin)

	var list []Track

	if *trackListTypeFlag == "basic" {
		list = parsers.ParseBasicTracklist(bufReader)
	} else if *trackListTypeFlag == "virtualdj" {
		list = parsers.ParseVirtualDJTracklist(bufReader)
	} else if *trackListTypeFlag == "serato" {
		parsers.ParseSeratoTracklist(bufReader)
	}

	return list
}

func sortMapByValue(m map[string][]Track) ArtistTrackList {
	p := make(ArtistTrackList, len(m))
	i := 0

	for k, v := range m {
		p[i] = ArtistTrack{k, v}
		i++
	}

	sort.Sort(sort.Reverse(p))

	return p
}

func countArtists(tracks []Track) map[string][]Track {
	//make a map
	var artists map[string][]Track = make(map[string][]Track)

	for _, track := range tracks {
		artists[track.Artist] = append(artists[track.Artist], track)
	}

	return artists
}
