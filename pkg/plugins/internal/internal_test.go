package internal

import (
	"strconv"
	"testing"

	"github.com/dchenk/mazewire/pkg/types/version"
)

func TestDiskName(t *testing.T) {
	cases := []struct {
		p    Registration
		name string
	}{
		{
			p: Registration{
				Name: "irrelevant",
				Id:   "the_id",
				Ver:  &version.Version{Major: 3, Minor: 1, Patch: 5},
			},
			name: "the_id3_1_5",
		},
	}
	for i := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			got := diskName(&cases[i].p)
			if got != cases[i].name {
				t.Errorf("got %q but expected %q", got, cases[i].name)
			}
		})
	}
}
