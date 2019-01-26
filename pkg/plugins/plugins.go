// Package plugins defines the Plugin interface and other types that work in the plugins system.
package plugins

import (
	"regexp"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/email"
	"github.com/dchenk/mazewire/pkg/filters"
	"github.com/dchenk/mazewire/pkg/hooks"
	"github.com/dchenk/mazewire/pkg/plugins/internal"
	"github.com/dchenk/mazewire/pkg/plugins/specs"
	"github.com/dchenk/mazewire/pkg/types/version"
	"github.com/golang/protobuf/proto"
)

// A Plugin is a plugin that responds to hooks that it has registered to handle.
type Plugin interface {
	// Identity returns the plugin's basic identifying metadata.
	// This function is called before the main plugin connection is established.
	Identity() (name string, id string, ver version.Version)

	// Specs returns the specifications within which the plugin works.
	//
	// For example, the plugin may indicate that it should be called asynchronously, or that it needs certain
	// database access capabilities, or that it depends on a particular version of the main app system or other
	// plugins.
	Specs() []specs.Spec

	// Hooks specifies all of the hooks which the plugin wants to handle with its hook handlers.
	Hooks() map[hooks.Hook]hooks.Handler

	// Filters specifies all of the filters which the plugins wants to handle with its filter handlers.
	Filters() map[filters.Filter]filters.Handler

	// CustomHooks returns the list of all of the custom hooks the plugin is registering.
	//
	// One plugin's custom hooks are not the same as another plugin's custom hooks with the same name because
	// the name of each of a plugin's custom hooks is prefixed with the plugin ID. So for one plugin to hook
	// into another plugin's registered custom hooks, the Plugin must list the other plugin's hook name with
	// the other plugin's ID followed by a dot followed by the hook name.
	//
	// So, for example, suppose a plugin with ID "alice" registers a custom hook named "jump" and another plugin
	// with ID "bob" wants to hook into that. The "bob" plugin needs to include its Handler with the hooks map
	// it returns in the Hooks function, the key of this custom hook being "alice.jump".
	CustomHooks() []hooks.Hook

	// CustomFilters returns the list of all of the custom filters the plugin is registering.
	// Namespacing across plugins works just as it does for hooks.
	CustomFilters() []filters.Filter
}

// Serve starts up a plugin's serving to the main host system. This function should be called within each plugin
// package's main function. This function blocks while the connection with the host system persists.
func Serve(p Plugin) {
	internal.Serve(makeRegistration(p), p.Hooks(), p.Filters())
}

// makeRegister generates a internal.Registration object given the implementation of a Plugin.
// The created object is then used to begin communication with the main app server.
func makeRegistration(p Plugin) *internal.Registration {
	name, id, ver := p.Identity()
	r := internal.Registration{
		Name: name,
		Id:   id,
		Ver:  &ver,
	}

	s := p.Specs()
	r.Specs = make([]*internal.Spec, len(s))
	for i := range s {
		r.Specs[i] = makeSpec(&s[i])
	}

	pHooks := p.Hooks()
	r.Hooks = make(map[string]*internal.Handler, len(pHooks))
	for k, v := range pHooks {
		s = v.Specs()
		hs := make([]*internal.Spec, len(s))
		for i := range s {
			hs[i] = makeSpec(&s[i])
		}
		r.Hooks[string(k)] = &internal.Handler{Specs: hs}
	}

	pFilters := p.Filters()
	r.Filters = make(map[string]*internal.Handler, len(pFilters))
	for k, v := range pFilters {
		s = v.Specs()
		fs := make([]*internal.Spec, len(s))
		for i := range s {
			fs[i] = makeSpec(&s[i])
		}
		r.Filters[string(k)] = &internal.Handler{Specs: fs}
	}

	ch := p.CustomHooks()
	r.CustomHooks = make([]string, len(ch))
	for i := range ch {
		r.CustomHooks[i] = string(ch[i])
	}

	cf := p.CustomFilters()
	r.CustomFilters = make([]string, len(cf))
	for i := range cf {
		r.CustomFilters[i] = string(cf[i])
	}

	return &r
}

func makeSpec(s *specs.Spec) *internal.Spec {
	return &internal.Spec{Type: uint32(s.Type), Value: s.Value.String()}
}

// DoHook calls all of the registered hook handles for the named hook with the given data.
// The slice returned indicates the number of hook handlers that were called.
//
// If an error is returned from any of the handlers, the messages already received are returned along with the
// error, and no more handlers are called.
func DoHook(hook hooks.Hook, payload proto.Message) ([]proto.Message, error) {
	return internal.HostSideHooks.Do(hook, payload)
}

// DoFilter calls all of the registered filter handlers of the named filter serially.
// If an error is returned from any of the handlers, the data that is available at that point is returned with
// the error without further filtering.
func DoFilter(filter filters.Filter, payload proto.Message) (proto.Message, error) {
	return internal.HostSideFilters.Do(filter, payload)
}

// Activate activates a plugin on the host side.
func Activate(p Plugin) {
	_, id, _ := p.Identity()
	internal.Activate(id, p.Specs(), p.Hooks(), p.Filters(), p.CustomHooks(), p.CustomFilters())
}

// Deactivate deactivates a plugin on the host side.
func Deactivate(pluginID string) {
	internal.Deactivate(pluginID)
}

// PluginCapabilities contains the capabilities a Plugin may use for its features. Only the user-authorized and
// explicitly requested capabilities are set here as non-nil.
type PluginCapabilities struct {
	data.SiteGetter
	data.SiteInserter
	data.SiteDeleter
	data.SiteManager

	data.BlobGetter
	data.BlobInserter
	data.BlobDeleter
	data.BlobManager

	data.ContentGetter
	data.ContentInserter
	data.ContentDeleter
	data.ContentManager

	data.UserGetter
	data.UserInserter
	data.UserDeleter
	data.UserManager

	data.UserMetaGetter
	data.UserMetaInserter
	data.UserMetaDeleter
	data.UserMetaManager

	data.OptionGetter
	data.OptionInserter
	data.OptionDeleter
	data.OptionManager

	email.Sender
}

// ValidID says whether id is a valid identifier for a plugin. An ID must be between 2 and 20 characters
// long, begin and end with a lowercase English character, and contain only lowercase English characters and
// underscores.
func ValidID(id string) bool {
	return validID.MatchString(id)
}

var validID = regexp.MustCompile("[a-z][a-z_]{0,18}[a-z]")
