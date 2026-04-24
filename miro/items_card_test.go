package miro

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetItem_CardDescriptionAndAssignee(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards/b1/items/card1" {
			t.Errorf("path %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":   "card1",
			"type": "card",
			"data": map[string]any{
				"title":       "T",
				"description": "body text",
				"dueDate":     "2025-12-01",
				"assignee":    map[string]any{"id": "u1", "name": "Alex"},
			},
			"position": map[string]any{"x": 0.0, "y": 0.0},
		})
	}))
	defer server.Close()
	c := newTestClientWithServer(server.URL)
	res, err := c.GetItem(context.Background(), GetItemArgs{BoardID: "b1", ItemID: "card1"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Description != "body text" {
		t.Errorf("Description=%q", res.Description)
	}
	if res.DueDate != "2025-12-01" {
		t.Errorf("DueDate=%q", res.DueDate)
	}
	if res.Assignee != "Alex" {
		t.Errorf("Assignee=%q", res.Assignee)
	}
}

func TestSearchBoard_MatchesDescription(t *testing.T) {
	pages := [][]map[string]any{
		{
			{
				"id":   "1",
				"type": "card",
				"data": map[string]any{
					"title":       "X",
					"description": "uniqueprobe-abc description here",
				},
			},
		},
	}
	n := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n >= len(pages) {
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []any{}, "cursor": ""})
			return
		}
		body := map[string]any{
			"data":   pages[n],
			"cursor": "",
		}
		n++
		_ = json.NewEncoder(w).Encode(body)
	}))
	defer server.Close()
	c := newTestClientWithServer(server.URL)
	out, err := c.SearchBoard(context.Background(), SearchBoardArgs{
		BoardID: "b1",
		Query:   "uniqueprobe",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Count != 1 || len(out.Matches) != 1 {
		t.Fatalf("count=%d matches=%d", out.Count, len(out.Matches))
	}
	if out.Matches[0].Type != "card" {
		t.Errorf("type=%s", out.Matches[0].Type)
	}
}
