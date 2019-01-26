package main

import "time"

// PageVersions is a timestamp-sortable slice representing version of a page (or post).
// The sort.Sort implementation sorts the slice in reverse-chronological order.
type PageVersions []PageVersion

// A PageVersion is a saved version of a page or post. The data may represent a working draft
// or a compiled page.
type PageVersion struct {
	Id          int64                         `msgp:"id"`     // ID of the record in the Blobs table.
	Data        []byte                        `msgp:"data"`   // The body of the page
	Static      bool                          `msgp:"static"` // Whether Data contains the final, built static HTML to display.
	CSS         string                        `msgp:"css"`    // The user's custom code if this is a draft, otherwise all compiled CSS.
	Styles      map[string]string             `msgp:"styles"` // room-type key-value style pairs
	FeaturedImg string                        `msgp:"img"`
	Scripts     map[string]PageResourceConfig `msgp:"scripts"`
	Stylesheets map[string]PageResourceConfig `msgp:"stylesheets"`
	Meta        map[string]string             `msgp:"meta"` // various additional settings
	Timestamp   time.Time                     `msgp:"ts"`
}

// Len returns the number of elements in the slice.
func (pv PageVersions) Len() int {
	return len(pv)
}

// Less compares the stringified timestamps of the two elements. The slice is supposed to
// be ordered reverse-chronologically. So the returned bool says if time at j is less than
// time at i.
func (pv PageVersions) Less(i, j int) bool {
	return pv[j].Timestamp.Before(pv[i].Timestamp)
}

// Swap swaps the elements at i and j.
func (pv PageVersions) Swap(i, j int) {
	pv[i], pv[j] = pv[j], pv[i]
}

// PageResourceConfig is a configuration of a static resource that should be included with a page.
type PageResourceConfig struct {
	Src  string `msgp:"src"`  // A complete URL to the resource's location.
	Head bool   `msgp:"head"` // Whether the resource should be linked to in the head of the page.
}
