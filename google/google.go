package google

import (
	"encoding/json"
	"fmt"
	"github.com/ruxton/term"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var GOOGLE_SEARCH_KEY string
var GOOGLE_SEARCH_CX string
var GOOGLE_URL = "https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&searchType=image&q=%s&start=%d&userIp=%s"
var HTTP_USER_AGENT = "Mix Cover Builder v" + versions.VERSION

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

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
	query := artist + " - " + song + " cover art"

	local_ip := GetLocalIP()
	cover_search_url := fmt.Sprintf(GOOGLE_URL, url.QueryEscape(query), GOOGLE_SEARCH_KEY, GOOGLE_SEARCH_CX, start, local_ip)
	term.OutputMessage(cover_search_url + "\n")
	request := BuildHttpRequest(cover_search_url, "GET")
	client := http.Client{}
	resp, doError := client.Do(request)
	if doError != nil {
		term.OutputError("Error fetching artwork for " + artist + " - " + song + ": " + doError.Error())
		return ""
	} else {
		defer resp.Body.Close()
	}

	var responseObj *Response = new(Response)

	err := json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		term.OutputError(err.Error())
		os.Exit(2)
	}

	var songObj *ResponseItem = new(ResponseItem)

	if resp.StatusCode != 200 {
		term.OutputError(fmt.Sprintf("Error fetching artwork from google - %d\n%v", resp.StatusCode, responseObj))
	} else {

		text_end := len(song)
		if text_end > 15 {
			text_end = 15
		}

		for _, result := range responseObj.Items {
			if (strings.Contains(result.Snippet, "cover") ||
				strings.Contains(result.Snippet, song) ||
				strings.Contains(result.Title, song)) ||
				(strings.HasPrefix(result.Snippet, song[0:text_end])) {
				songObj = &result
				break
			}
		}
	}
	if songObj.Link == "" {
		term.OutputMessage(term.Green + "." + term.Reset)
	}
	return songObj.Link
}

func GetCoverFor(artist string, song string) string {
	start := 1
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
