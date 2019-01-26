package internal

import (
	"sync"
	"time"

	"github.com/dchenk/mazewire/pkg/hooks"
	hashiplug "github.com/hashicorp/go-plugin"
)

const pluginSetName = "mw_plugin"

// Activate is used by the host system to activate a plugin for the current server instance.
// This does not register the activation of a plugin in the system database.
func Activate(reg *Registration) {
	activePlugins.activate(reg)
}

func Deactivate(reg *Registration) {
	activePlugins.deactivate(reg)
}

type activePluginClient struct {
	hashiplug.Client
	grpcClient
	//diskName string
	//reg      *Registration
}

type activePluginsHost struct {
	mu sync.RWMutex

	// active contains each active plugin's disk name (ID with version) along with a client that can be
	// used to connect to the running plugin.
	active map[string]activePluginClient
}

// activate adds and activates a plugin and starts up its binary.
// The function first write the plugin to disk if necessary.
func (ap *activePluginsHost) activate(reg *Registration) {
	pluginDiskName := diskName(reg)

	ap.mu.Lock()
	defer ap.mu.Unlock()

	// Check if the plugin is already running active.
	if _, ok := ap.active[pluginDiskName]; ok {
		return
	}

	// TODO: Check if the plugin is on disk.

	cc := newClientConfig()
	cc.Plugins = hashiplug.PluginSet{
		pluginSetName: &pluginImplementation{
			// TODO
		},
	}
	hc := hashiplug.NewClient(cc)

	_, err := hc.Start()
	if err != nil {
		// TODO
	}

	for hookName, hh := range reg.Hooks {
		// Since we already checked that this plugin is not running active, there shouldn't be any
		// handlers for this plugin.
		// TODO: append to the slice
		HostSideHooks.Add(pluginDiskName, hooks.Hook(hookName), hh.Specs)
	}

}

// deactivate deactivates the plugin with the ID and version given in pluginID, removing all of
// the associated hook handlers and the connection client.
func (ap *activePluginsHost) deactivate(reg *Registration) {
	ap.mu.Lock()
	defer ap.mu.Unlock()
}

var activePlugins = activePluginsHost{
	active: make(map[string]activePluginClient, 8),
}

// newClientConfig returns the configuration used for gRPC plugins.
func newClientConfig() *hashiplug.ClientConfig {
	return &hashiplug.ClientConfig{
		HandshakeConfig:  handshakeConfig(),
		Managed:          true,
		StartTimeout:     time.Second * 30,
		AllowedProtocols: []hashiplug.Protocol{hashiplug.ProtocolGRPC},
	}
}
