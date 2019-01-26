package main

import (
	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/types/plugins_list"
	"github.com/golang/protobuf/proto"
)

const pluginsListOption = "plugins_list"

// getPlugins returns the list of currently registered plugins in their Protocol Buffers format.
func getPlugins(siteID int64) (*plugins_list.PluginsList, error) {
	opts, err := data.Conn.OptionByKey(siteID, pluginsListOption)
	if err != nil {
		return nil, err
	}
	var list plugins_list.PluginsList
	err = proto.Unmarshal(opts.V, &list)
	return &list, err
}
