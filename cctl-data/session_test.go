package main

import "testing"

type stubProcessChecker struct {
	alivePIDs map[int]bool
}

func (s *stubProcessChecker) IsAlive(pid int) bool {
	return s.alivePIDs[pid]
}

func TestLoadRawSessions(t *testing.T) {
	sessions, err := LoadRawSessions("testdata/sessions")
	if err != nil {
		t.Fatalf("LoadRawSessions: %v", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestFilterAliveSessions(t *testing.T) {
	sessions, _ := LoadRawSessions("testdata/sessions")
	checker := &stubProcessChecker{alivePIDs: map[int]bool{99999: true}}
	alive := FilterAlive(sessions, checker)
	if len(alive) != 1 {
		t.Fatalf("expected 1 alive session, got %d", len(alive))
	}
	if alive[0].PID != 99999 {
		t.Errorf("expected PID 99999, got %d", alive[0].PID)
	}
}
