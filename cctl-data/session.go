package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"syscall"
)

type RawSession struct {
	PID       int    `json:"pid"`
	SessionID string `json:"sessionId"`
	CWD       string `json:"cwd"`
	StartedAt int64  `json:"startedAt"`
	UpdatedAt int64  `json:"updatedAt"`
	Status    string `json:"status"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
}

type ProcessChecker interface {
	IsAlive(pid int) bool
}

type OSProcessChecker struct{}

func (o *OSProcessChecker) IsAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

func LoadRawSessions(dir string) ([]RawSession, error) {
	pattern := filepath.Join(dir, "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var sessions []RawSession
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var s RawSession
		if err := json.Unmarshal(data, &s); err != nil {
			continue
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func FilterAlive(sessions []RawSession, checker ProcessChecker) []RawSession {
	var alive []RawSession
	for _, s := range sessions {
		if checker.IsAlive(s.PID) {
			alive = append(alive, s)
		}
	}
	return alive
}
