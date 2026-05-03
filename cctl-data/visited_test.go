package main

import (
	"path/filepath"
	"testing"
)

func TestLoadVisited_Empty(t *testing.T) {
	dir := t.TempDir()
	visited := LoadVisited(filepath.Join(dir, "visited.json"))
	if len(visited) != 0 {
		t.Errorf("expected empty map, got %d entries", len(visited))
	}
}

func TestMarkVisited_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "visited.json")

	MarkVisited(path, "session-123")
	visited := LoadVisited(path)

	ts, ok := visited["session-123"]
	if !ok {
		t.Fatal("session-123 not found in visited")
	}
	if ts <= 0 {
		t.Errorf("expected positive timestamp, got %d", ts)
	}
}

func TestNeedsAttention(t *testing.T) {
	tests := []struct {
		name      string
		status    string
		updatedAt int64
		visited   int64
		want      bool
	}{
		{"busy never needs attention", "busy", 1777550000000, 0, false},
		{"idle unseen", "idle", 1777550000000, 0, true},
		{"idle seen", "idle", 1777550000000, 1777551000, false},
		{"idle updated after seen", "idle", 1777552000000, 1777551000, true},
		{"waiting unseen", "waiting", 1777550000000, 0, true},
		{"stale never needs attention", "unknown", 1777550000000, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NeedsAttention(tt.status, tt.updatedAt, tt.visited)
			if got != tt.want {
				t.Errorf("NeedsAttention(%q, %d, %d) = %v, want %v",
					tt.status, tt.updatedAt, tt.visited, got, tt.want)
			}
		})
	}
}
