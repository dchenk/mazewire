package users

import "testing"

func TestRoleAtLeast(t *testing.T) {

	cases := []struct {
		actual, min string
		val         bool
	}{
		{RoleSuper, RoleAdmin, true},
		{RoleSuper, RoleAuthor, true},
		{RoleSuper, RoleNone, true},
		{"something_weird", RoleSubscriber, false},
		{RoleOwner, RoleOwner, true},
		{RoleOwner, RoleAuthor, true},
		{RoleNone, RoleNone, true},
	}

	for i, tc := range cases {

		got := RoleAtLeast(tc.actual, tc.min)
		if got != tc.val {
			t.Errorf("(index %d) got wrong value %v", i, got)
		}

	}

}

func TestRoleIsValid(t *testing.T) {

	valid := []string{RoleOwner, RoleAdmin, RoleEditor, RoleAuthor, RoleSubscriber, RoleNone}
	invalid := []string{RoleSuper, "something_else"}

	for _, r := range valid {
		if !RoleIsValid(r) {
			t.Errorf("saying that role %q is not valid", r)
		}
	}

	for _, r := range invalid {
		if RoleIsValid(r) {
			t.Errorf("saying that role %q is valid", r)
		}
	}

}
