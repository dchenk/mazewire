package main

import (
	"context"
	"net/http"

	"github.com/dchenk/mazewire/pkg/data"
)

// getUserMessages retrieves the stored notification messages for the user.
func getUserMessages(ctx context.Context, userID int64) ([]data.UserMessage, error) { // TODO
	//q := datastore.NewQuery(userMessageKind).Filter("uid =", userID)
	//it := datastoreClient.Run(ctx, q)
	//ums := make([]userMessage, 0, 2)
	//for {
	//	var um userMessage
	//	_, err := it.Next(&um)
	//	if err == iterator.Done {
	//		break
	//	}
	//	if err != nil {
	//		return ums, err
	//	}
	//	ums = append(ums, um)
	//}
	//return ums, nil
	return nil, nil
}

func saveUserMessage(r *http.Request, userID int64, key string, msg string) { // TODO
	m := &data.UserMessage{
		UserId:  userID,
		K:       key,
		Message: msg,
	}
	_ = m
	//_, err := datastoreClient.Put(context.Background(), datastore.IncompleteKey(userMessageKind, nil), m)
	//if err != nil {
	//	log.Err(r, fmt.Sprintf("could not save user message: k = %q; v = %q", key, msg), err)
	//}
}

// getSiteMessages retrieves the stored notification messages for the user.
func getSiteMessages(siteID int64, role string) ([]data.SiteMessage, error) { // TODO
	//q := datastore.NewQuery(siteMessageKind).Filter("sid =", siteID).Filter("role =", role)
	//it := datastoreClient.Run(ctx, q)
	//sms := make([]siteMessage, 0, 2)
	//for {
	//	var sm siteMessage
	//	_, err := it.Next(&sm)
	//	if err == iterator.Done {
	//		break
	//	}
	//	if err != nil {
	//		return sms, err
	//	}
	//	sms = append(sms, sm)
	//}
	//return sms, nil
	return nil, nil
}

// saveSiteMessage saves a message to be shown to users who have at least the role specified on the site.
func saveSiteMessage(r *http.Request, siteID int64, role string, key string, msg string) { // TODO
	m := &data.SiteMessage{
		SiteId:  siteID,
		Role:    role,
		K:       key,
		Message: msg,
	}
	_ = m
	//_, err := datastoreClient.Put(r.Context(), datastore.IncompleteKey(siteMessageKind, nil), m)
	//if err != nil {
	//	log.Err(r, fmt.Sprintf("could not save site message: k = %q; v = %q", key, msg), err)
	//}
}
