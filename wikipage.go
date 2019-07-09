package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type pages struct {
	Pageid float64
	Ns     float64
	Title  string
	Text   string `json:"extract"`
}
type query struct {
	Pages map[string]pages
}

type pageExtract struct {
	Batchcomplete string
	Warnings      interface{}
	Query         query
}

// WPQueryURIFmtString contains a format string that requires a list of Page IDs (only one Page ID is used here)
// See https://www.mediawiki.org/wiki/API:Query for more details
const WPQueryURIFmtString = "https://en.wikipedia.org/w/api.php?action=query&prop=extracts&pageids=%s&format=json&exlimit=1&explaintext=true&exsectionformat=raw"

// GetWikiPageContent returns the content and title of a Wikipedia as strings page identified by pageId
func GetWikiPageContent(pageID string) (string, string, error) {
	var content = ""
	var title = ""
	var err error

	if _, err := strconv.ParseInt(pageID, 10, 64); err == nil {
		// urlFormat deviates from suggested URI in spec to elimnate HTML fron returned content
		var url = fmt.Sprintf(WPQueryURIFmtString, pageID)

		var resp, err = http.Get(url)

		if err != nil {
			log.Errorf("http.Get returned err : %s , check connection to Internet, DNS resolution and that Wikipedia is up", err)
		} else {
			defer resp.Body.Close()

			var extract pageExtract

			var dec = json.NewDecoder(resp.Body)
			if err := dec.Decode(&extract); err != nil {
				log.Errorf("json.NewDecoder(resp.Body).Decode(%#v) %q!\n", extract, err)
			}

			title = extract.Query.Pages[pageID].Title
			content = extract.Query.Pages[pageID].Text
			log.Debug(url)
		}
	} else {
		log.WithFields(log.Fields{"page_id": pageID}).Errorf("wp.GetWikiPageContent : pageID is not a valid number, REST API request NOT sent to Wikipedia")
		return "", "", errors.New("page_id is not a valid number")
	}
	return content, title, err
}
