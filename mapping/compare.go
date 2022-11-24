package mapping

import (
	"bytes"
	"encoding/json"
)

func compareInJSON(a, b interface{}) (bool, error) {
	ab, err := json.Marshal(a)
	if err != nil {
		return false, err
	}

	bb, err := json.Marshal(b)
	if err != nil {
		return false, err
	}

	return bytes.Equal(ab, bb), nil
}
