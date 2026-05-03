package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type RecapResult struct {
	Text   string
	Source string // "away_summary", "remember", "snapshot"
}

func projectKey(cwd string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return re.ReplaceAllString(cwd, "-")
}

func ReadAwaySummary(cwd, sessionID, projectsDir string) string {
	key := projectKey(cwd)
	jsonlPath := filepath.Join(projectsDir, key, sessionID+".jsonl")

	f, err := os.Open(jsonlPath)
	if err != nil {
		return ""
	}
	defer f.Close()

	var lastContent string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, `"away_summary"`) {
			continue
		}
		var entry struct {
			Content string `json:"content"`
			Subtype string `json:"subtype"`
		}
		if json.Unmarshal([]byte(line), &entry) == nil && entry.Subtype == "away_summary" {
			lastContent = entry.Content
		}
	}

	lastContent = strings.TrimSuffix(lastContent, " (disable recaps in /config)")
	return lastContent
}

func ReadRememberRecap(cwdRememberDir string) string {
	nowFile := filepath.Join(cwdRememberDir, "now.md")
	sourceFile := ""

	if info, err := os.Stat(nowFile); err == nil && info.Size() > 0 {
		sourceFile = nowFile
	} else {
		pattern := filepath.Join(cwdRememberDir, "today-*.md")
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			sort.Slice(matches, func(i, j int) bool {
				fi, _ := os.Stat(matches[i])
				fj, _ := os.Stat(matches[j])
				return fi.ModTime().After(fj.ModTime())
			})
			if info, err := os.Stat(matches[0]); err == nil && info.Size() > 0 {
				sourceFile = matches[0]
			}
		}
	}

	if sourceFile == "" {
		return ""
	}

	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return ""
	}

	// Extract body of last ## section
	lines := strings.Split(string(data), "\n")
	var body []string
	found := false
	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			found = true
			body = nil
			continue
		}
		if found {
			body = append(body, line)
		}
	}

	result := strings.Join(body, " ")
	result = collapseSpaces(result)
	return result
}

func ReadSnapshotDone(slug, snapshotsDir string) string {
	return readSnapshotSection(slug, snapshotsDir, "Done")
}

func ReadSnapshotNext(slug, snapshotsDir string) string {
	text := readSnapshotSection(slug, snapshotsDir, "Next")
	if text == "" || strings.HasPrefix(text, "No explicit next") || strings.HasPrefix(text, "No goal captured") {
		return ""
	}
	if !strings.HasSuffix(text, "?") {
		return ""
	}
	return text
}

func readSnapshotSection(slug, snapshotsDir, section string) string {
	file := filepath.Join(snapshotsDir, slug, "latest.md")
	data, err := os.ReadFile(file)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	var body []string
	inSection := false
	for _, line := range lines {
		if line == "## "+section {
			inSection = true
			continue
		}
		if inSection && strings.HasPrefix(line, "## ") {
			break
		}
		if inSection && strings.TrimSpace(line) != "" {
			body = append(body, line)
			if len(body) >= 3 {
				break
			}
		}
	}

	result := strings.Join(body, " ")
	result = collapseSpaces(result)

	if strings.HasPrefix(result, "No goal captured") || strings.HasPrefix(result, "No summary available") {
		return ""
	}
	return result
}

func fileMtime(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.ModTime().Unix()
}

func GetFreshestRecap(cwd, sessionID, slug, projectsDir, rememberDir, snapshotsDir string) RecapResult {
	key := projectKey(cwd)
	jsonlPath := filepath.Join(projectsDir, key, sessionID+".jsonl")

	var awayLine string
	var awayTS int64
	if _, err := os.Stat(jsonlPath); err == nil {
		awayLine = ReadAwaySummary(cwd, sessionID, projectsDir)
		if awayLine != "" {
			awayTS = fileMtime(jsonlPath)
		}
	}

	var rememberLine string
	var rememberTS int64
	rememberFile := findRememberFile(rememberDir)
	if rememberFile != "" {
		rememberLine = ReadRememberRecap(rememberDir)
		if rememberLine != "" {
			rememberTS = fileMtime(rememberFile)
		}
	}

	if awayLine != "" && rememberLine != "" {
		if rememberTS > awayTS {
			return RecapResult{Text: rememberLine, Source: "remember"}
		}
		return RecapResult{Text: awayLine, Source: "away_summary"}
	}
	if awayLine != "" {
		return RecapResult{Text: awayLine, Source: "away_summary"}
	}
	if rememberLine != "" {
		return RecapResult{Text: rememberLine, Source: "remember"}
	}

	snapDone := ReadSnapshotDone(slug, snapshotsDir)
	if snapDone != "" {
		return RecapResult{Text: snapDone, Source: "snapshot"}
	}

	return RecapResult{}
}

func findRememberFile(dir string) string {
	nowFile := filepath.Join(dir, "now.md")
	if info, err := os.Stat(nowFile); err == nil && info.Size() > 0 {
		return nowFile
	}
	pattern := filepath.Join(dir, "today-*.md")
	matches, _ := filepath.Glob(pattern)
	if len(matches) > 0 {
		sort.Slice(matches, func(i, j int) bool {
			fi, _ := os.Stat(matches[i])
			fj, _ := os.Stat(matches[j])
			return fi.ModTime().After(fj.ModTime())
		})
		if info, err := os.Stat(matches[0]); err == nil && info.Size() > 0 {
			return matches[0]
		}
	}
	return ""
}

func collapseSpaces(s string) string {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}
