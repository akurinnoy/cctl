package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

func LoadVisited(path string) map[string]int64 {
	data, err := os.ReadFile(path)
	if err != nil {
		return make(map[string]int64)
	}
	var visited map[string]int64
	if json.Unmarshal(data, &visited) != nil {
		return make(map[string]int64)
	}
	return visited
}

func MarkVisited(path, sessionID string) error {
	visited := LoadVisited(path)
	visited[sessionID] = time.Now().Unix()

	data, err := json.MarshalIndent(visited, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func NeedsAttention(status string, updatedAtMs int64, lastVisited int64) bool {
	switch status {
	case "busy":
		return false
	case "idle", "waiting":
		updatedSec := updatedAtMs / 1000
		return updatedSec > lastVisited
	default:
		return false
	}
}
