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

func TestNullOrNilBool(t *testing.T) {
	cases := []struct {
		input  sql.NullBool
		expect interface{}
	}{
		{sql.NullBool{Valid: false}, nil},
		{sql.NullBool{Bool: true, Valid: true}, true},
		{sql.NullBool{Bool: false, Valid: true}, false},
	}

	for _, c := range cases {
		got := nullOrNilBool(c.input)
		if got != c.expect {
			t.Errorf("nullOrNilBool(%#v) = %#v; want %#v", c.input, got, c.expect)
		}
	}
}
