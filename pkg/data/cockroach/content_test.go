package cockroach

import "testing"

func TestStatusList(t *testing.T) {

	cases := []struct {
		statuses []string
		list     string
	}{
		{[]string{"a"}, " AND status IN ('a')"},
		{[]string{"a", "b"}, " AND status IN ('a','b')"},
	}

	for i, pair := range cases {
		got := statusList(pair.statuses)
		if got != pair.list {
			t.Errorf("(index %d) did not construct list correctly; got: %s", i, got)
		}
	}

}

func TestGetOffsetOrNot(t *testing.T) {

	cases := []struct {
		n   uint64
		str string
	}{
		{0, ""},
		{1, " OFFSET 1"},
		{23, " OFFSET 23"},
	}

	for i, pair := range cases {
		got := getOffsetOrNot(pair.n)
		if got != pair.str {
			t.Errorf("(index %d) did not construct offset correctly; got: %s", i, got)
		}
	}

}
