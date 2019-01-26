package email

import (
	"strconv"
	"testing"
)

func TestEmail_MarshalJSON(t *testing.T) {

	cases := []struct {
		em     Email
		result string
	}{
		{Email{}, `{"to":{"name":"","address":""},"from":{"name":"","address":""},"subject":"","body":""}`},
	}

	for i, tc := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			out, err := tc.em.MarshalJSON()
			if err != nil {
				t.Fatalf("got an error; %v", err)
			}
			str := string(out)
			if str != tc.result {
				t.Errorf("got incorrect result: %s", str)
			}
		})
	}

}

func TestParty_MarshalJSON(t *testing.T) {
	t.FailNow() // TODO
}

func TestCopies_MarshalJSON(t *testing.T) {

	cases := []struct {
		c      Copies
		result string
	}{
		{Copies{}, `{"cc":[],"bcc":[]}`},
		{
			c: Copies{
				CC:  []Party{{"Abc Def", "abcd@abcd.efg"}},
				BCC: []Party{{"Ghi Jkl", "ghi@jlk.mn.op"}, {"Qrs", "tu_v.wx@yz.ab"}},
			},
			result: `{"cc":[{"name":"Abc Def","address":"abcd@abcd.efg"}],"bcc":[{"name":"Ghi Jkl","address":"ghi@jlk.mn.op"},{"name":"Qrs","address":"tu_v.wx@yz.ab"}]}`,
		},
	}

	for i, tc := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			out, err := tc.c.MarshalJSON()
			if err != nil {
				t.Fatalf("got an error; %v", err)
			}
			str := string(out)
			if str != tc.result {
				t.Errorf("got incorrect result: %s", str)
			}
		})
	}

	t.FailNow()

}

func TestContent_MarshalJSON(t *testing.T) {
	t.FailNow() // TODO
}
