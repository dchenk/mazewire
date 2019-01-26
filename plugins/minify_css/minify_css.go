package main

import (
	"github.com/dchenk/mazewire/pkg/filters"
	"github.com/dchenk/mazewire/pkg/hooks"
	"github.com/dchenk/mazewire/pkg/plugins"
	"github.com/dchenk/mazewire/pkg/plugins/specs"
	"github.com/dchenk/mazewire/pkg/types/version"
	"github.com/golang/protobuf/proto"
)

const (
	name     = "Minify CSS"
	id       = "minify_css"
	verMajor = 0
	verMinor = 1
	verPatch = 1
)

func main() {
	plugins.Serve(&inst)
}

// inst is an instance if the plugin implementation.
var inst = minifyCSS{}

// minifyCSS implements the plugins.Plugin interface.
type minifyCSS struct {
}

func (m *minifyCSS) Identity() (string, string, version.Version) {
	return name, id, version.Version{Major: verMajor, Minor: verMinor, Patch: verPatch}
}

func (m *minifyCSS) Specs() []specs.Spec {
	return []specs.Spec{{Type: specs.PageRender}}
}

func (m *minifyCSS) Hooks() map[hooks.Hook]hooks.Handler {
	return map[hooks.Hook]hooks.Handler{hooks.UserCSS: new(handler)}
}

func (m *minifyCSS) Filters() map[filters.Filter]filters.Handler {
	return nil
}

func (m *minifyCSS) CustomHooks() []hooks.Hook { return nil }

func (m *minifyCSS) CustomFilters() []filters.Filter { return nil }

// A handler handles the UserCSS hook. The output is minified CSS code.
type handler struct{}

func (h *handler) Specs() []specs.Spec {
	return nil
}

func (h *handler) Handle(_ []interface{}, payload proto.Message) (proto.Message, error) {
	// TODO
	return nil, nil
}
