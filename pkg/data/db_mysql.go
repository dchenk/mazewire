// +build mysql

package data

import "github.com/dchenk/mazewire/pkg/data/mysql"

func init() {
	Conn = new(mysql.DB)
}
