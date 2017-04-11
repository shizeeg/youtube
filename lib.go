package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	// URI is an API endpoint
	URI = "https://www.googleapis.com/youtube/v3/videos?id=%s&part=contentDetails&key=%s"
)

// ErrKey YouTube API key error
var ErrKey = errors.New("get a key here: https://console.developers.google.com/apis/api/youtube.googleapis.com")

// GetDuration returns duration from YouTube video id
// you should pass a valid YouTube Data API Key. Can get one here:
// https://console.developers.google.com/apis/api/youtube.googleapis.com"
func GetDuration(key, id string) (duration string, err error) {
	resp, err := http.Get(fmt.Sprintf(URI, id, key))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var yt ContentDetails
	if err = json.NewDecoder(resp.Body).Decode(&yt); err != nil {
		return
	}
	for _, item := range yt.Items {
		return ParseDuration(item.ContentDetails.Duration), nil
	}
	return
}

// ParseDuration converts duration in YouTube format to 15:04:05
func ParseDuration(dur string) string {
	layouts := []string{
		"PT15H4M5S",
		"PT4M5S",
		"PT5S",
	}
	for _, l := range layouts {
		t, err := time.Parse(l, dur)
		if err == nil {
			tmp := strings.TrimLeft(t.Format("15:04:05"), "00:")
			if strings.Index(tmp, ":") == -1 { // < 1 min
				return tmp + "s"
			}
			return tmp
		}
	}
	return ""
}

// IDs extracts ID parts from YouTube links found in a string
func IDs(msg string) (ids []string) {
	re := regexp.MustCompile(`(v=|youtu\.be\/)(?P<id>[0-9A-Za-z_-]{11})`)
	for _, match := range re.FindAllString(msg, -1) {
		//fmt.Println(match, "found at index", i)
		ids = append(ids, re.ReplaceAllString(match, "${id}"))
	}
	return ids
}

// ContentDetails is a response struct for YouTube API call
type ContentDetails struct {
	Etag  string `json:"etag"`
	Items []struct {
		ContentDetails struct {
			Caption         string `json:"caption"`
			Definition      string `json:"definition"`
			Dimension       string `json:"dimension"`
			Duration        string `json:"duration"`
			LicensedContent bool   `json:"licensedContent"`
			Projection      string `json:"projection"`
		} `json:"contentDetails"`
		Etag string `json:"etag"`
		ID   string `json:"id"`
		Kind string `json:"kind"`
	} `json:"items"`
	Kind     string `json:"kind"`
	PageInfo struct {
		ResultsPerPage int64 `json:"resultsPerPage"`
		TotalResults   int64 `json:"totalResults"`
	} `json:"pageInfo"`
}
