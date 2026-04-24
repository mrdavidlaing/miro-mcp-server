package miro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ConnectorPoint is a position on a connector free end. Miro's API sometimes
// returns x and y as JSON strings; encoding/json would otherwise fail. Anchor
// positions on attached items use percentage strings (e.g. "50%") for relative
// placement along the edge; we normalize those to 0..1 and keep the wire form
// in RawX/RawY when set.
type ConnectorPoint struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Origin string  `json:"origin,omitempty"`
	// RawX, RawY are the original string forms when Miro sent a non-numeric
	// value such as a percentage. Empty when x/y were plain numbers.
	RawX string `json:"x_raw,omitempty"`
	RawY string `json:"y_raw,omitempty"`
}

// UnmarshalJSON decodes x/y as numbers, numeric strings, or percentage strings
// like "50%" (relative anchor, stored as 0.5 in X or Y and the full token in RawX/RawY).
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
	if p.X, p.RawX, err = decodeConnectorCoordinate(w.X, "x"); err != nil {
		return err
	}
	if p.Y, p.RawY, err = decodeConnectorCoordinate(w.Y, "y"); err != nil {
		return err
	}
	p.Origin = w.Origin
	return nil
}

func decodeConnectorCoordinate(raw json.RawMessage, field string) (value float64, rawString string, err error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return 0, "", nil
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return 0, "", err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, "", nil
		}
		// Miro: relative connection point along an item side, e.g. "50%" = midpoint.
		if after, ok := strings.CutSuffix(s, "%"); ok {
			after = strings.TrimSpace(after)
			n, perr := strconv.ParseFloat(after, 64)
			if perr != nil {
				return 0, "", fmt.Errorf("parse %s: %w", field, perr)
			}
			return n / 100.0, s, nil
		}
		f, ferr := strconv.ParseFloat(s, 64)
		if ferr != nil {
			return 0, "", fmt.Errorf("parse %s: %w", field, ferr)
		}
		return f, "", nil
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err != nil {
		return 0, "", fmt.Errorf("parse %s: %w", field, err)
	}
	return f, "", nil
}
