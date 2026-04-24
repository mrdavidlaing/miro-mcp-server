package miro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// ConnectorPoint is a position on a connector free end. Miro's API sometimes
// returns x and y as JSON strings; encoding/json would otherwise fail.
type ConnectorPoint struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Origin string  `json:"origin,omitempty"`
}

// UnmarshalJSON decodes x/y as numbers or numeric strings.
func (p *ConnectorPoint) UnmarshalJSON(data []byte) error {
	var w struct {
		X      json.RawMessage `json:"x"`
		Y      json.RawMessage `json:"y"`
		Origin string          `json:"origin,omitempty"`
	}
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	var err error
	if p.X, err = decodeFlexibleFloat(w.X, "x"); err != nil {
		return err
	}
	if p.Y, err = decodeFlexibleFloat(w.Y, "y"); err != nil {
		return err
	}
	p.Origin = w.Origin
	return nil
}

func decodeFlexibleFloat(raw json.RawMessage, field string) (float64, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return 0, nil
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return 0, err
		}
		if s == "" {
			return 0, nil
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, fmt.Errorf("parse %s: %w", field, err)
		}
		return f, nil
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err != nil {
		return 0, fmt.Errorf("parse %s: %w", field, err)
	}
	return f, nil
}
