package json

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	s := Student{
		Name:  "Steve",
		Grade: 33,
		Score: 56,
	}

	data, err := marshal(s)
	if err != nil {
		t.Error(err)
	}

	actual, err := unmarshal(data)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actual, s) {
		t.Errorf("expected: %v\ngot: %v\n", s, actual)
	}
}
