package itunes

import (
	"strings"
)

type Response struct {
	Error       string   `json:"errormessage,omitempty"`
	ResultCount int      `json:"resultCount,omitempty"`
	Results     []Result `json:"results,omitempty"`
}

type Result struct {
	wrapperType            string  `json:"wrapperType,omitempty"`
	Kind                   string  `json:"kind,omitempty"`
	ArtistId               int64   `json:"artistId,omitempty"`
	CollectionId           int64   `json:"collectionId,omitempty"`
	TrackId                int64   `json:"trackId,omitempty"`
	ArtistName             string  `json:"artistName,omitempty"`
	CollectionName         string  `json:"collectionName,omitempty"`
	TrackName              string  `json:"trackName,omitempty"`
	CollectionCensoredName string  `json:"collectionCensoredName,omitempty"`
	TrackCensoredName      string  `json:"trackCensoredName,omitempty"`
	ArtistViewUrl          string  `json:"artistViewUrl,omitempty"`
	CollectionViewUrl      string  `json:"collectionViewUrl,omitempty"`
	TrackViewUrl           string  `json:"trackViewUrl,omitempty"`
	PreviewUrl             string  `json:"previewUrl,omitempty"`
	ArtworkUrl30           string  `json:"artworkUrl30,omitempty"`
	ArtworkUrl60           string  `json:"artworkUrl60,omitempty"`
	ArtworkUrl100          string  `json:"artworkUrl100,omitempty"`
	CollectionPrice        float64 `json:"collectionPrice,omitempty"`
	TrackPrice             float64 `json:"trackPrice,omitempty"`
	ReleaseDate            string  `json:"releaseDate,omitempty"`
	CollectionExplicitness string  `json:"collectionExplicitness,omitempty"`
	TrackExplicitness      string  `json:"trackExplicitness,omitempty"`
	DiscCount              int64   `json:"discCount,omitempty"`
	DiscNumber             int64   `json:"discNumber,omitempty"`
	TrackCount             int64   `json:"trackCount,omitempty"`
	TrackNumber            int64   `json:"trackNumber,omitempty"`
	TrackTimeMillis        int64   `json:"trackTimeMillis,omitempty"`
	Country                string  `json:"country,omitempty"`
	Currency               string  `json:"currency,omitempty"`
	PrimareyGenreName      string  `json:"primaryGenreName,omitempty"`
	ContentAdvisoryRating  string  `json:"contentAdvisoryRating,omitempty"`
	RadioStationUrl        string  `json:"radioStationUrl,omitempty"`
	IsStreamable           bool    `json:"isStreamable,omitempty"`
}

func (r *Result) ArtworkUrl600() string {
	return strings.Replace(r.ArtworkUrl100, "100x100", "600x600", -1)
}
