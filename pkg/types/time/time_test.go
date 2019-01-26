package time

import (
	"strconv"
	"testing"
	"time"
)

func TestTime_ToTime(t *testing.T) {
	cases := []struct {
		t  Time
		tt time.Time
	}{
		{t: Time{Seconds: 0, Nanos: 0}, tt: time.Unix(0, 0)},
		{t: Time{Seconds: -150, Nanos: 10}, tt: time.Unix(-150, 10)},
		{t: Time{Seconds: -62135596800, Nanos: 0}, tt: time.Unix(-62135596800, 0)},
		{t: Time{Seconds: 253402300799, Nanos: 1e9 - 1}, tt: time.Unix(253402300799, 1e9-1)},
	}
	for i := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			c := &cases[i]
			got, err := c.t.ToTime()
			if err != nil {
				t.Fatalf("got an error; %v", err)
			}
			if !got.Equal(c.tt) {
				t.Errorf("got unequal times; wanted %s but got %s", c.tt.Format(time.RFC3339Nano), got.Format(time.RFC3339Nano))
			}
		})
	}
}
