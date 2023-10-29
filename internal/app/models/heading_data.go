package models

import (
	"fmt"
	"encoding/json"
)

type HeadingData struct {
	Data []Node `json:"data"`
}

func (r *HeadingData) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, r)
	case string:
		return json.Unmarshal([]byte(src), r)
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, *r)
	}
}
