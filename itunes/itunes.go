package itunes

import (
	"encoding/json"
	"fmt"
	"github.com/ruxton/term"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var ITMS_URL = "http://ax.itunes.apple.com/WebObjects/MZStoreServices.woa/wa/wsSearch?term=%s&country=%s&entity=%s"
var HTTP_USER_AGENT = ""

func BuildHttpRequest(url string, request string) *http.Request {
	req, err := http.NewRequest(request, url, nil)
	if err != nil {
		fmt.Printf("Error - %s", err.Error())
	}

	if HTTP_USER_AGENT != "" {
		req.Header.Set("User-Agent", HTTP_USER_AGENT)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req
}

func GetCoverFor(artist string, song string) string {
	query := artist + " " + song

	cover_search_url := fmt.Sprintf(ITMS_URL, url.QueryEscape(query), "AU", "song")

	request := BuildHttpRequest(cover_search_url, "GET")
	client := http.Client{}
	resp, doError := client.Do(request)
	defer resp.Body.Close()

	if doError != nil {
		term.OutputError("Error fetching artwork for " + artist + " - " + song + ": " + doError.Error())
	}

	var responseObj *Response = new(Response)

	err := json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		term.OutputError(err.Error())
		os.Exit(2)
	}

	var songObj *Result = new(Result)
	var emptySongObj *Result = new(Result)

	if responseObj.Error != "" {
		fmt.Println(responseObj.Error)
	} else {
		for _, result := range responseObj.Results {
			if ((result.ArtistName == artist) && (result.TrackName == song)) ||
				strings.HasPrefix(result.TrackName, song) {
				songObj = &result
				break
			}
		}
	}
	if songObj != emptySongObj {
		return songObj.ArtworkUrl600()
	} else {
		return ""
	}
}
