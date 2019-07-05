package client

import (
	"os"
	"testing"

	"github.com/romana/rlog"
)

func TestGetWikiPageContentError(t *testing.T) {
	cases := []struct {
		testName                         string
		pageID                           string
		expectedErr                      string
		expectedApproximateContentLength int
		alloweLengthVariancePercentage   int
	}{
		{"Error is empty test", "21721040", "", 1000, 5},
		{"Error is set : invalid PageID", "2172104asdasd0", "wp.GetWikiPageContent ERROR pageID \"2172104asdasd0\" is not a number, REST API request NOT sent to Wikipedia", 0, 5},
	}
	rlog.SetOutput(os.Stderr)
	for _, c := range cases {
		r, err := GetWikiPageContent(c.pageID)
		if err != nil && err.Error() != c.expectedErr {
			t.Fatalf("%v : \nExpected err: %#v\nActual err  : %#v", c.testName, c.expectedErr, err)
		}

		rlog.Info(r)
	}
}
