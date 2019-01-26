// +build postgres

package data

import "github.com/dchenk/mazewire/pkg/data/postgres"

func init() {
	Conn = new(postgres.DB)
}
