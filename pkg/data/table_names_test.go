package data

import (
	"strconv"
	"testing"
)

func TestBlobsTable(t *testing.T) {
	cases := []struct {
		siteID int64
		table  string
	}{
		{5, "blobs5"},
		{0, "blobs"},
		{43588, "blobs43588"},
	}
	for i, c := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			got := BlobsTable(c.siteID)
			if got != c.table {
				t.Errorf("got blobs table name %q but expected %q", got, c.table)
			}
		})
	}
}
