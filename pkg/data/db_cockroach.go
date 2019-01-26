// +build cockroach

package data

import "github.com/dchenk/mazewire/pkg/data/cockroach"

func init() {
	Conn = new(cockroach.DB)
}
