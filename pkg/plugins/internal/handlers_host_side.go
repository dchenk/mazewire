package internal

import (
	"context"
	"sync"

	"github.com/dchenk/mazewire/pkg/filters"
	"github.com/dchenk/mazewire/pkg/hooks"
	"github.com/dchenk/mazewire/pkg/plugins/specs"
	"github.com/golang/protobuf/proto"
)

type hostSideHooks struct {
	mu      sync.RWMutex
	clients map[hooks.Hook][]hostHandler
}

type hostHandler struct {
	PluginClient
	diskName string
	specs    []*specs.Spec
}

// Add adds a handler of a named hook. The pluginID string must include the plugin version so that it's
// possible to remove handlers of specific versions of a plugin.
func (hhs *hostSideHooks) Add(pluginDiskName string, hook hooks.Hook, hookSpecs []*specs.Spec, client PluginClient) {
	hhs.mu.Lock()
	defer hhs.mu.Unlock()

	existing := hhs.clients[hook]

	// Make sure this plugin client hasn't already been added.
	for _, h := range existing {
		if h.diskName == pluginDiskName {
			return
		}
	}

	hhs.clients[hook] = append(existing, hostHandler{
		PluginClient: client,
		diskName:     pluginDiskName,
		specs:        hookSpecs,
	})
}

func (hhs *hostSideHooks) Do(hook hooks.Hook, payload proto.Message) ([]proto.Message, error) {
	data, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	reqData := HookRequest{
		HookName: string(hook),
		Data:     data,
	}

	// Make a copy of the slice of clients to avoid concurrent reads and writes to the slice
	// within the internal map.
	var handlerClients []hostHandler
	{
		hhs.mu.RLock()
		hh := hhs.clients[hook]
		handlerClients = make([]hostHandler, len(hh))
		for i := range hh {
			handlerClients[i] = hh[i]
		}
		hhs.mu.RUnlock()
	}

	responses := make([]proto.Message, 0, len(handlerClients))

	for _, h := range handlerClients {
		resp, err := h.PluginClient.DoHook(context.Background(), &reqData)
		if err != nil {
			return responses, err
		}
		r, err := hooks.ToMessageResponse(hook, resp.Data)
		if err != nil {
			return responses, err
		}
		responses = append(responses, r)
	}

	return responses, nil
}

func (hhs *hostSideHooks) Remove(pluginDiskName string) {
	// TODO
}

var HostSideHooks = hostSideHooks{
	clients: make(map[hooks.Hook][]hostHandler, 16),
}

type hostSideFilters struct {
	mu      sync.RWMutex
	clients map[filters.Filter][]hostHandler
}

// Add adds a handler of a named filter. The pluginDiskName string must include the plugin version so that it's
// possible to remove handlers of specific versions of a plugin.
func (hhf *hostSideFilters) Add(pluginDiskName string, filter filters.Filter, client PluginClient) {
	// TODO
}

func (hhf *hostSideFilters) Do(filter filters.Filter, payload proto.Message) (proto.Message, error) {
	data, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	reqData := FilterData{
		FilterName: string(filter),
		Data:       data,
	}

	// Make a copy of the slice of clients to avoid concurrent reads and writes to the slice
	// within the internal map.
	var handlerClients []hostHandler
	{
		hhf.mu.RLock()
		fh := hhf.clients[filter]
		handlerClients = make([]hostHandler, len(fh))
		for i := range fh {
			handlerClients[i] = fh[i]
		}
		hhf.mu.RUnlock()
	}

	for _, h := range handlerClients {
		resp, err := h.PluginClient.DoFilter(context.Background(), &reqData)
		if err != nil {
			return payload, err
		}
		reqData.Data = resp.Data
	}

	msg, err := filters.ToMessage(filter, reqData.Data)
	if err != nil {
		return payload, err
	}

	return msg, nil
}

func (hhf *hostSideFilters) Remove(pluginDiskName string) {
	// TODO
}

var HostSideFilters = hostSideFilters{
	clients: make(map[filters.Filter][]hostHandler, 16),
}
