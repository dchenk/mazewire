package internal

import (
	"github.com/hashicorp/go-plugin"
)

const (
	handshakeVersion        = 1
	handshakeMagicCookieKey = "MAZEWIRE_PLUGIN"
)

func handshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		ProtocolVersion:  handshakeVersion,
		MagicCookieKey:   handshakeMagicCookieKey,
		MagicCookieValue: pluginHandshakeCookieVal(),
	}
}

func pluginHandshakeCookieVal() string {
	return "mazewire_plugin"
}
