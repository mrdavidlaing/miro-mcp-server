package miro

import (
	"encoding/json"
	"testing"
)

func TestConnector_StringPositionOnEndpoint(t *testing.T) {
	const payload = `{
		"id": "conn-1",
		"type": "connector",
		"startItem": {
			"item": "a1",
			"position": {"x": "10.5", "y": "-20", "origin": "center"}
		},
		"endItem": { "item": "b1" }
	}`
	var c Connector
	if err := json.Unmarshal([]byte(payload), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.StartItem.Position == nil {
		t.Fatal("expected start position")
	}
	if c.StartItem.Position.X != 10.5 || c.StartItem.Position.Y != -20.0 {
		t.Errorf("position = %v,%v want 10.5,-20", c.StartItem.Position.X, c.StartItem.Position.Y)
	}
}

func TestConnector_NumericPositionUnchanged(t *testing.T) {
	const payload = `{
		"id": "c2",
		"type": "connector",
		"startItem": { "position": { "x": 1, "y": 2 } },
		"endItem": { "item": "b1" }
	}`
	var c Connector
	if err := json.Unmarshal([]byte(payload), &c); err != nil {
		t.Fatal(err)
	}
	if c.StartItem.Position == nil {
		t.Fatal("missing position")
	}
	if c.StartItem.Position.X != 1 || c.StartItem.Position.Y != 2 {
		t.Errorf("got %g %g", c.StartItem.Position.X, c.StartItem.Position.Y)
	}
}
