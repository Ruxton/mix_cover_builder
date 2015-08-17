package google

import (
	"encoding/json"
	"fmt"
	"github.com/ruxton/term"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var GOOGLE_URL = "https://ajax.googleapis.com/ajax/services/search/images?v=1.0&q=%s&rsz=8&start=%d&userip=192.168.1.205"
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

func GetPaginatedCoverFor(artist string, song string, start int) string {
	query := "site:bandcamp.com " + artist + " - " + song + " bandcamp cover art"

	cover_search_url := fmt.Sprintf(GOOGLE_URL, url.QueryEscape(query), start)

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

	if responseObj.Status != 200 {
		term.OutputError(fmt.Sprintf("Error fetching artwork from google - %d", responseObj.Status))
	} else {
		for _, result := range responseObj.Data.Results {
			if strings.Contains(result.OriginalContextUrl, "bandcamp.com") &&
				(strings.Contains(result.ContentNoFormatting, "cover art") ||
					strings.Contains(result.ContentNoFormatting, song) ||
					strings.Contains(result.TitleNoFormatting, song)) ||
				(strings.HasPrefix(result.ContentNoFormatting, song[0:15])) {
				songObj = &result
				break
			}
		}
	}
	if songObj.Url == "" {
		term.OutputMessage(term.Green + "." + term.Reset)
	}
	return songObj.Url
}

func GetCoverFor(artist string, song string) string {
	start := 0
	steps := 8
	MAX_STEP := 32
	url := ""
	for (url == "") && (start < MAX_STEP+1) {
		url = GetPaginatedCoverFor(artist, song, start)
		start = start + steps
		if url == "" {
			time.Sleep(1 * time.Second)
		}
	}

	return url
}
