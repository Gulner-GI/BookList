package handlers

import (
	"database/sql"
	"testing"
)

func TestNullOrNil(t *testing.T) {
	cases := []struct {
		input  sql.NullString
		expect interface{}
	}{
		{sql.NullString{Valid: false}, nil},
		{sql.NullString{String: "foo", Valid: true}, "foo"},
	}

	for _, c := range cases {
		got := nullOrNil(c.input)
		if got != c.expect {
			t.Errorf("nullOrNil(%#v) = %#v; want %#v", c.input, got, c.expect)
		}
	}
}
