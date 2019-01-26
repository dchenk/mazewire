// Package sample has sample Tree structures useful for testing.
package sample

import "github.com/dchenk/mazewire/pkg/room"

// SampleTree1 is a simple sample tree structure.
var SampleTree1 = room.Tree{
	room.Section{
		Type: "standard",
		Rows: []room.Row{
			{
				Type: "fullwidth",
				Modules: [][]room.Module{
					{
						{Type: "text"},
						{Data: []byte("something here")},
					},
				},
			},
			{
				Type: "halfhalf",
				Modules: [][]room.Module{
					{
						{Type: "stuff"},
						{Data: []byte("something here")},
					},
					{
						{Type: "other"},
						{Data: []byte("something here")},
					},
				},
			},
			{
				Type: "twothirdsthird",
				Modules: [][]room.Module{
					{
						{Type: "more"},
						{Data: []byte("something here")},
					},
					{
						{Type: "stuff"},
						{Data: []byte("something here")},
					},
					{
						{Type: "mod-type"},
						{Data: []byte("something here")},
					},
				},
			},
		},
	},
	room.Section{
		Type: "fullwidth",
		Rows: []room.Row{},
	},
}
