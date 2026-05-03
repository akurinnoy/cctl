package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadAwaySummary(t *testing.T) {
	projectsDir := "testdata/projects"
	recap := ReadAwaySummary("/tmp/test-project", "test-session-alive", projectsDir)
	if recap == "" {
		t.Fatal("expected non-empty away summary")
	}
	if recap != "Fixed the auth bug and pushed to main." {
		t.Errorf("unexpected recap: %q", recap)
	}
}

func TestReadRememberRecap(t *testing.T) {
	recap := ReadRememberRecap("testdata/remember")
	if recap == "" {
		t.Fatal("expected non-empty remember recap")
	}
}

func TestReadSnapshotDone(t *testing.T) {
	recap := ReadSnapshotDone("test-project", "testdata/snapshots")
	if recap == "" {
		t.Fatal("expected non-empty snapshot done")
	}
}

func TestReadSnapshotNext(t *testing.T) {
	next := ReadSnapshotNext("test-project", "testdata/snapshots")
	if next != "Should we run the integration tests?" {
		t.Errorf("unexpected next: %q", next)
	}
}

func TestGetFreshestRecap_PrefersNewer(t *testing.T) {
	rememberPath := filepath.Join("testdata", "remember", "now.md")
	now := time.Now()
	os.Chtimes(rememberPath, now, now)

	jsonlDir := "testdata/projects"
	jsonlPath := filepath.Join(jsonlDir, "-tmp-test-project", "test-session-alive.jsonl")
	past := now.Add(-10 * time.Minute)
	os.Chtimes(jsonlPath, past, past)

	recap := GetFreshestRecap("/tmp/test-project", "test-session-alive", "test-project",
		jsonlDir, "testdata/remember", "testdata/snapshots")
	if recap.Source != "remember" {
		t.Errorf("expected remember source (newer mtime), got %q", recap.Source)
	}
}
