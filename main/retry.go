package main

import (
	"net/http"

	"github.com/dchenk/mazewire/pkg/try"
)

//var retryApiActions = map[string]func(*http.Request, *data.Site, *data.User, *APIResponse){
//(createBlobsTableRetrier{}).Key(): retryCreateBlobsTable,
//}

const retryKind = "Retry"

// saveRetry saves to the datastore the settings for a retrier to schedule it or allow the user to
// manually retry.
func saveRetry(r *http.Request, rt try.Trier) error { // TODO
	//_, err := datastoreClient.Put(r.Context(), datastore.IncompleteKey(retryKind, nil), rt)
	// TODO: ping the Service that will make the request
	return nil
}
