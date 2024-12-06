package json

import (
	"encoding/json"
)

type Student struct {
	Name  string
	Grade int
	Score int
}

func marshal(s Student) ([]byte, error) {
	return json.Marshal(s)
}

func unmarshal(data []byte) (Student, error) {
	s := Student{}
	err := json.Unmarshal(data, &s)
	return s, err
}
