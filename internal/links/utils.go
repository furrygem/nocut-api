package links

import (
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"unicode/utf8"
)

// URLToID converts provided base62 string to hex
func URLToID(b62 string) (string, error) {
	var i big.Int
	_, ok := i.SetString(b62, 62)
	if !ok {
		return "", fmt.Errorf("Can't parse %s as base62", b62)
	}
	return i.Text(16), nil

}

// IDToURL converts provided string containing hex to base62
func IDToURL(id string) (string, error) {
	var i big.Int
	_, ok := i.SetString(id, 16)
	if !ok {
		return "", fmt.Errorf("Can't parse %s as hex", id)
	}
	return i.Text(62), nil
}

// URLHostIsUp checks if specified URL is up and returned status < 500
func URLHostIsUp(sourceURL string) (bool, error) {
	resp, err := http.Get(sourceURL)
	if err != nil {
		return false, err
	}
	if resp.StatusCode < 500 {
		return true, nil
	}
	return false, fmt.Errorf("Host has returned http status code %d", resp.StatusCode)
}

// URLIsValid checks if specified URL can be parsed and is absolute
func URLIsValid(sourceURL string) (bool, error) {
	u, err := url.Parse(sourceURL)
	if err != nil {
		return false, err
	}
	if !(u.IsAbs()) {
		return false, fmt.Errorf("URL %s is not absolute", sourceURL)
	}
	return true, nil
}

// CheckURLLength checks if specified URL length is not less than minLen
func CheckURLLength(sourceURL string, minLen int) (bool, error) {
	if utf8.RuneCountInString(sourceURL) < minLen {
		return false, fmt.Errorf("URL %s is too short", sourceURL)
	}
	return true, nil
}

// AddSlugToLink populates Link object with slug derived from ID
func AddSlugToLink(l *Link) error {
	slug, err := IDToURL(l.ID)
	if err != nil {
		return err
	}
	l.Slug = slug
	return nil
}

// AppendPrefixToURL Appends specified prefix to the beginning of the URL
func AppendPrefixToURL(prefix string, url string) string {
	r := path.Join("/", prefix, url)
	return r
}
