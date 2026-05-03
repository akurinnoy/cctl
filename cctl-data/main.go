package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type OutputSession struct {
	PID            int    `json:"pid"`
	Alive          bool   `json:"alive"`
	CWD            string `json:"cwd"`
	Slug           string `json:"slug"`
	Status         string `json:"status"`
	Label          string `json:"label"`
	UpdatedAt      int64  `json:"updatedAt"`
	StartedAt      int64  `json:"startedAt"`
	Age            string `json:"age"`
	Name           string `json:"name"`
	SessionID      string `json:"sessionId"`
	Recap          string `json:"recap"`
	RecapSource    string `json:"recapSource"`
	NextAction     string `json:"nextAction"`
	NeedsAttention bool   `json:"needsAttention"`
}

func classifyStatus(raw string) string {
	switch raw {
	case "busy":
		return "BUSY"
	case "waiting":
		return "ASK"
	case "idle":
		return "IDLE"
	default:
		return "STALE"
	}
}

func sessionsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "sessions")
}

func projectsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "projects")
}

func snapshotsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "ok-session-context")
}

func visitedPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "ok-session-context", "visited.json")
}

func cmdSessions() error {
	raw, err := LoadRawSessions(sessionsDir())
	if err != nil {
		return err
	}

	checker := &OSProcessChecker{}
	visited := LoadVisited(visitedPath())

	var output []OutputSession
	for _, s := range raw {
		alive := checker.IsAlive(s.PID)
		if !alive {
			continue
		}

		slug := GetSlug(s.CWD)
		label := classifyStatus(s.Status)

		ts := s.UpdatedAt
		if label == "STALE" {
			if s.UpdatedAt > 0 {
				ts = s.UpdatedAt
			} else {
				ts = s.StartedAt
			}
		}

		rememberDir := filepath.Join(s.CWD, ".remember")
		recap := GetFreshestRecap(s.CWD, s.SessionID, slug,
			projectsDir(), rememberDir, snapshotsDir())

		nextAction := ReadSnapshotNext(slug, snapshotsDir())

		lastVisited := visited[s.SessionID]
		attention := NeedsAttention(s.Status, s.UpdatedAt, lastVisited)

		output = append(output, OutputSession{
			PID:            s.PID,
			Alive:          alive,
			CWD:            s.CWD,
			Slug:           slug,
			Status:         s.Status,
			Label:          label,
			UpdatedAt:      s.UpdatedAt,
			StartedAt:      s.StartedAt,
			Age:            RelativeAge(ts),
			Name:           s.Name,
			SessionID:      s.SessionID,
			Recap:          recap.Text,
			RecapSource:    recap.Source,
			NextAction:     nextAction,
			NeedsAttention: attention,
		})
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func cmdRecap(pidStr string) error {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Errorf("invalid PID: %s", pidStr)
	}

	raw, err := LoadRawSessions(sessionsDir())
	if err != nil {
		return err
	}

	for _, s := range raw {
		if s.PID != pid {
			continue
		}
		slug := GetSlug(s.CWD)
		rememberDir := filepath.Join(s.CWD, ".remember")
		recap := GetFreshestRecap(s.CWD, s.SessionID, slug,
			projectsDir(), rememberDir, snapshotsDir())
		if recap.Text != "" {
			fmt.Println(recap.Text)
		}
		return nil
	}
	return nil
}

func cmdVisited(sessionID string) error {
	return MarkVisited(visitedPath(), sessionID)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: cctl-data <sessions|recap|visited> [args]\n")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "sessions":
		err = cmdSessions()
	case "recap":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: cctl-data recap <PID>\n")
			os.Exit(1)
		}
		err = cmdRecap(os.Args[2])
	case "visited":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: cctl-data visited <sessionId>\n")
			os.Exit(1)
		}
		err = cmdVisited(os.Args[2])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
