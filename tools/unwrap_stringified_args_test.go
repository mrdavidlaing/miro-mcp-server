package tools

import (
	"encoding/json"
	"testing"
)

func TestUnwrapStringifiedInRawJSON_FieldsString(t *testing.T) {
	raw := []byte(`{"board_id":"b1","title":"t1","status":"connected","fields":"[{\"value\":\"A\",\"iconShape\":\"round\"}]"}`)
	out, err := unwrapStringifiedInRawJSON(raw)
	if err != nil {
		t.Fatalf("unwrap: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	f, ok := m["fields"].([]any)
	if !ok || len(f) != 1 {
		t.Fatalf("fields: got %#v, want 1-element array", m["fields"])
	}
	o, _ := f[0].(map[string]any)
	if o["value"] != "A" {
		t.Errorf("value = %v, want A", o["value"])
	}
}

func TestUnwrapStringifiedInRawJSON_ItemsString(t *testing.T) {
	raw := []byte(`{"board_id":"b1","items":"[{\"type\":\"text\",\"x\":0}]"}`)
	out, err := unwrapStringifiedInRawJSON(raw)
	if err != nil {
		t.Fatalf("unwrap: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	it, ok := m["items"].([]any)
	if !ok || len(it) != 1 {
		t.Fatalf("items: got %#v", m["items"])
	}
}

func TestUnwrapStringifiedInRawJSON_ArrayUnchanged(t *testing.T) {
	want := `{"board_id":"b1","fields":[{"value":"A"}]}`
	out, err := unwrapStringifiedInRawJSON([]byte(want))
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != want {
		t.Errorf("got %s, want %s", out, want)
	}
}
