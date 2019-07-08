package client

import (
	"os"
	"strings"
	"testing"

	"github.com/romana/rlog"
)

// These tests use live Wikipedia to test funcionality developed.
// I would normally mock out a third party service say using Wirmock
// or implementa canned response myself depending on the range of
// calls made to the third part service.
func TestGetWikiPageContentError(t *testing.T) {
	cases := []struct {
		testName                         string
		pageID                           string
		expectedTitle                    string
		expectedContentSample            string
		expectedErr                      string
		expectedApproximateContentLength int
		alloweLengthVariancePercentage   int
	}{
		{"Error is empty test", "21721040", "Stack Overflow", "Jeff Atwood", "", 1000, 5},
		{"Error is empty test", "21721040", "Stack Overflow", "Jeff Atwood", "", 1000, 5},
		{"Error is set : invalid PageID", "2172104asdasd0", "", "", "wp.GetWikiPageContent ERROR pageID \"2172104asdasd0\" is not a number, REST API request NOT sent to Wikipedia", 0, 5},
	}

	rlog.SetOutput(os.Stderr)

	for _, c := range cases {
		actualContent, actualTitle, err := GetWikiPageContent(c.pageID)
		if err != nil && err.Error() != c.expectedErr {
			t.Fatalf("%v : \nExpected err: %#v\nActual err  : %#v", c.testName, c.expectedErr, err)
		}
		if actualTitle != c.expectedTitle {
			t.Fatalf("%v : \nExpected title: %#v\nActual err  : %#v", c.testName, c.expectedTitle, actualTitle)
		}
		if !strings.Contains(actualContent, c.expectedContentSample) {
			t.Logf("Warning : %v : \nExpected content: %#v not present in downloaded content\nActual err  : %#v\nCheck the page content and modify the text if necessary.", c.testName, c.expectedTitle, actualTitle)
		}
	}
}
