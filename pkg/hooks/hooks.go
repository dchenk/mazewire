package hooks

import (
	"github.com/dchenk/mazewire/pkg/hooks/payloads"
	"github.com/dchenk/mazewire/pkg/plugins/specs"
	"github.com/golang/protobuf/proto"
)

// A Hook is a string that is used to connect together events within an application.
type Hook string

// A Handler is a type that can respond a particular hook event.
type Handler interface {
	// Specs returns the specifications with the Handler must be called.
	// For example, the Handler may indicate that it should be called asynchronously as the main app
	// carries on. Or, as another example, a Spec may use a DependencySpec to specify that another
	// plugin's Handler must be called before this Handler.
	//
	// Capabilities should not be specified here, since they're specified only at the plugin level.
	Specs() []specs.Spec

	// Handle handles the registered hook.
	Handle(capabilities []interface{}, payload proto.Message) (proto.Message, error)
}

const (
	BeforePluginActivate   Hook = "before_plugin_activate"
	PluginActivated        Hook = "plugin_activated"
	BeforePluginDeactivate Hook = "before_plugin_deactivate"
	BeforePluginRemove     Hook = "before_plugin_remove"

	BeforeUserCreate Hook = "before_user_create"
	UserCreated      Hook = "user_created"

	BeforeSiteCreate Hook = "before_site_create"
	SiteCreated      Hook = "site_created"

	BeforeEmailSend Hook = "before_email_send"

	// Document-modifying hooks for pages.
	BeforePageHead Hook = "before_page_head"
	PageHead       Hook = "page_head"
	BeforePageBody Hook = "before_page_body"
	PageBody       Hook = "page_body"
	AfterPageBody  Hook = "after_page_body"
	PageFoot       Hook = "page_foot"

	// Document-modifying hooks for blog posts.
	BeforePostHead Hook = "before_post_head"
	PostHead       Hook = "post_head"
	BeforePostBody Hook = "before_post_body"
	PostBody       Hook = "post_body"
	AfterPostBody  Hook = "after_post_body"
	PostFoot       Hook = "post_foot"

	// AdminMenu hooks get an array of the menu links in the Admin UI and are expected to
	// return the filtered array.
	//AdminMenu Hook = "admin_menu"
	//AdminHead Hook = "admin_head"
	//AdminFoot Hook = "admin_foot"

	ContentCreate Hook = "content_create"

	BeforePagePublish Hook = "before_page_publish"
	BeforePostPublish Hook = "before_post_publish"

	BeforePageDelete Hook = "before_page_delete"
	BeforePostDelete Hook = "before_post_delete"

	BeforeMediaUpload Hook = "before_media_upload"
	MediaUpload       Hook = "media_upload"
	MediaDelete       Hook = "media_delete"
)

// ToMessageRequest returns the concrete proto.Message implementation for the request payload for the named hook,
// the data already unmarshalled.
func ToMessageRequest(hookName Hook, rawData []byte) (proto.Message, error) {
	switch hookName {
	default:
		var msg payloads.Custom
		err := proto.Unmarshal(rawData, &msg)
		return &msg, err
	}
}

// ToMessageResponse returns the concrete proto.Message implementation for the response payload for the named hook,
// the data already unmarshalled.
func ToMessageResponse(hookName Hook, rawData []byte) (proto.Message, error) {
	switch hookName {
	//case UserCSS:
	//	var msg payloads.UserCSS
	//	err := proto.Unmarshal(rawData, &msg)
	//	return &msg, err
	default:
		var msg payloads.Custom
		err := proto.Unmarshal(rawData, &msg)
		return &msg, err
	}
}
