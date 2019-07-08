package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/romana/rlog"
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

// GetWikiPageContent returns the content and title of a Wikipedia as strings page identified by pageId
func GetWikiPageContent(pageID string) (string, string, error) {
	var content = ""
	var title = ""
	var err error
	var errorMessage string

	if _, err := strconv.ParseInt(pageID, 10, 64); err == nil {
		// urlFormat deviates from suggested URI in spec to elimnate HTML fron returned content
		const urlFormat = "https://en.wikipedia.org/w/api.php?action=query&prop=extracts&pageids=%s&format=json&exlimit=1&explaintext=true&exsectionformat=raw"
		var url = fmt.Sprintf(urlFormat, pageID)

		var resp, err = http.Get(url)

		if err != nil {
			rlog.Error("")
			rlog.Errorf("http.Get returned err : %s , check connection to Internet, DNS resolution and that Wikipedia is up", err)
		}

		defer resp.Body.Close()

		var extract pageExtract

		var dec = json.NewDecoder(resp.Body)
		if err := dec.Decode(&extract); err != nil {
			rlog.Errorf("json.NewDecoder(resp.Body).Decode(%#v) %q!\n", extract, err)
		}

		title = extract.Query.Pages[pageID].Title
		content = extract.Query.Pages[pageID].Text
		rlog.Debug(url)
	} else {
		errorMessage = fmt.Sprintf("wp.GetWikiPageContent ERROR pageID \"%s\" is not a number, REST API request NOT sent to Wikipedia", pageID)
		err = errors.New(errorMessage)
		rlog.Error(err)
		return "", "", err
	}
	return content, title, err
}
