package errors

import (
	"errors"
	"testing"
)

func TestGetFields(t *testing.T) {
	fields := Fields{
		"key1": "value1",
		"key2": "value2",
	}
	fieldsLength := len(fields)

	err := E("Errors", fields)
	fs := err.(*Error).GetFields()
	fsLength := len(fs)

	if fieldsLength != fsLength {
		t.Error("fields length is different")
		return
	}

	for key := range fs {
		if fields[key] != fs[key] {
			t.Error("value is incorrect")
		}
	}
}

func TestMatch(t *testing.T) {
	cases := []struct {
		err1        error
		err2        error
		expectMatch bool
	}{
		{
			err1:        E(errors.New("This is new error")),
			err2:        nil,
			expectMatch: false,
		},
		{
			err1:        E(errors.New("This is new error")),
			err2:        errors.New("This is new error"),
			expectMatch: true,
		},
		{
			err1:        E(errors.New("This is new error")),
			err2:        errors.New("Something is different"),
			expectMatch: false,
		},
	}

	for _, val := range cases {
		if match := Match(val.err1, val.err2); match != val.expectMatch {
			t.Errorf("TestMatch: Expecting %v but got %v", val.expectMatch, match)
		}
	}
}
