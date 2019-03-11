package users

import "testing"

func TestRoleAtLeast(t *testing.T) {
	cases := []struct {
		actual, min Role
		val         bool
	}{
		// True
		{Role_SUPER, Role_SUPER, true},
		{Role_SUPER, Role_ADMIN, true},
		{Role_SUPER, Role_AUTHOR, true},
		{Role_SUPER, Role_NONE, true},
		{Role_OWNER, Role_OWNER, true},
		{Role_OWNER, Role_AUTHOR, true},
		{Role_OWNER, Role_NONE, true},
		{Role_EDITOR, Role_EDITOR, true},
		{Role_EDITOR, Role_AUTHOR, true},
		{Role_EDITOR, Role_SUBSCRIBER, true},
		{Role_AUTHOR, Role_AUTHOR, true},
		{Role_AUTHOR, Role_SUBSCRIBER, true},
		{Role_AUTHOR, Role_NONE, true},
		{Role_SUBSCRIBER, Role_SUBSCRIBER, false},
		{Role_SUBSCRIBER, Role_NONE, false},
		{Role_NONE, Role_NONE, true},

		// False
		{Role_OWNER, Role_SUPER, false},
		{Role_ADMIN, Role_OWNER, false},
		{Role_ADMIN, Role_SUPER, false},
		{Role_EDITOR, Role_SUPER, false},
		{Role_EDITOR, Role_ADMIN, false},
		{Role_AUTHOR, Role_SUPER, false},
		{Role_AUTHOR, Role_OWNER, false},
		{Role_AUTHOR, Role_ADMIN, false},
		{Role_AUTHOR, Role_EDITOR, false},
		{Role_SUBSCRIBER, Role_ADMIN, false},
		{Role_SUBSCRIBER, Role_EDITOR, false},
		{Role_SUBSCRIBER, Role_AUTHOR, false},
		{Role_NONE, Role_ADMIN, false},
		{Role_NONE, Role_SUBSCRIBER, false},
		{-1, Role_SUBSCRIBER, false},
		{17, Role_SUBSCRIBER, false},
	}

	for i, tc := range cases {
		got := RoleAtLeast(tc.actual, tc.min)
		if got != tc.val {
			t.Errorf("(%d) got wrong value %v", i, got)
		}
	}
}

func TestRoleIsValid(t *testing.T) {
	cases := []struct {
		num  int64
		role Role
		ok   bool
	}{
		{0, Role_NONE, true},
		{1, Role_SUBSCRIBER, true},
		{3, Role_AUTHOR, true},
		{5, Role_EDITOR, true},
		{7, Role_ADMIN, true},
		{9, Role_OWNER, true},
		{11, Role_SUPER, true},

		{-5, Role_NONE, false},
		{-1, Role_NONE, false},
		{2, Role_NONE, false},
		{10, Role_NONE, false},
		{800, Role_NONE, false},
	}

	for _, tc := range cases {
		res, ok := ValidRoleInt64(tc.num)
		if res != tc.role {
			t.Errorf("got %v for role %d, expected %v", res, tc.num, tc.role)
		}
		if ok != tc.ok {
			t.Errorf("got %v for role %d, expected %v", ok, tc.num, tc.ok)
		}
	}
}
