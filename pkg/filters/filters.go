package filters

import (
	"github.com/dchenk/mazewire/pkg/filters/payloads"
	"github.com/dchenk/mazewire/pkg/plugins/specs"
	"github.com/golang/protobuf/proto"
)

// A Filter is a string that is used to connect together events within an application.
type Filter string

type Handler interface {
	// Specs returns the specifications with the Handler must be called.
	// For example, the Handler may indicate that it should be called asynchronously as the main app
	// carries on. Or, as another example, a Spec may use a DependencySpec to specify that another
	// plugin's Handler must be called before this Handler.
	//
	// Capabilities should not be specified here, since they're specified only at the plugin level.
	Specs() []specs.Spec

	// Handle handles the registered hook. The concrete type of the returned Message is the same as
	// the concrete type of the payload Message.
	Handle(capabilities []interface{}, payload proto.Message) (proto.Message, error)
}

const (
	// UserCSS hooks get custom user CSS code before it is output on a page and are expected to
	// return the filtered code.
	UserCSS Filter = "user_css"

	// UserJS hooks get custom user JavaScript code before it is output on a page and are expected
	// to return the filtered code.
	UserJS Filter = "user_js"

	BeforeMediaUpload Filter = "before_media_upload"
)

func ToMessage(filter Filter, payload []byte) (proto.Message, error) {
	switch filter {
	case UserCSS:
		return nil, nil
	default:
		var msg payloads.Custom
		err := proto.Unmarshal(payload, &msg)
		return &msg, err
	}
}
