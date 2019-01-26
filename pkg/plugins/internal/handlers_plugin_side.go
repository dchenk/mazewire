package internal

import (
	"github.com/dchenk/mazewire/pkg/filters"
	"github.com/dchenk/mazewire/pkg/hooks"
	"github.com/golang/protobuf/proto"
	hashiplug "github.com/hashicorp/go-plugin"
)

func Serve(reg *Registration, hooks map[hooks.Hook]hooks.Handler, filters map[filters.Filter]filters.Handler) {
	PluginSideHooks.add(hooks)
	PluginSideFilters.add(filters)
	hashiplug.Serve(&hashiplug.ServeConfig{
		HandshakeConfig: handshakeConfig(),
		Plugins: hashiplug.PluginSet{
			reg.Id: &pluginImplementation{
				srv: pluginServer{reg},
			},
		},
		GRPCServer: hashiplug.DefaultGRPCServer,
	})
}

// pluginSideHooks contains a map of hook name -> handler on the plugin side (unused in the main host app).
type pluginSideHooks struct {
	handlers map[string]hooks.Handler // key is a hooks.Hook
}

// add adds the hooks and their handlers to the package's internal data structure.
// Each plugin must call this function exactly once (within plugins.Serve).
// This function must not be called concurrently from different goroutines.
func (hps *pluginSideHooks) add(hh map[hooks.Hook]hooks.Handler) {
	for hookName, handler := range hh {
		hps.handlers[string(hookName)] = handler
	}
}

// Do executes the hook handler registered on the plugin side.
func (hps *pluginSideHooks) Do(hook string, payload []byte) (*HookResponse, error) {
	h := hps.handlers[hook]
	if h == nil {
		return &HookResponse{}, nil // This should be impossible.
	}
	msg, err := hooks.ToMessageRequest(hooks.Hook(hook), payload)
	if err != nil {
		return nil, err
	}
	resp, err := h.Handle(nil, msg)
	if err != nil {
		return nil, err
	}
	respBytes, err := proto.Marshal(resp)
	return &HookResponse{Data: respBytes}, err
}

var PluginSideHooks = pluginSideHooks{
	handlers: make(map[string]hooks.Handler),
}

// pluginSideFilters contains a map of filter name -> handler on the plugin side (unused in the main host app).
type pluginSideFilters struct {
	handlers map[string]filters.Handler // key is a filters.Filter
}

// add adds the filters and their handlers to the package's internal data structure.
// Each plugin must call this function exactly once (within plugins.Serve).
// This function must not be called concurrently from different goroutines.
func (fps *pluginSideFilters) add(ff map[filters.Filter]filters.Handler) {
	for filterName, handler := range ff {
		fps.handlers[string(filterName)] = handler
	}
}

func (fps *pluginSideFilters) do(filter string, payload []byte) (*FilterData, error) {
	h := fps.handlers[filter]
	if h == nil {
		return &FilterData{Data: payload}, nil // This should be impossible.
	}
	msg, err := filters.ToMessage(filters.Filter(filter), payload)
	if err != nil {
		return nil, err
	}
	resp, err := h.Handle(nil, msg)
	if err != nil {
		return nil, err
	}
	respBytes, err := proto.Marshal(resp)
	return &FilterData{FilterName: filter, Data: respBytes}, err
}

var PluginSideFilters = pluginSideFilters{
	handlers: make(map[string]filters.Handler, 8),
}
