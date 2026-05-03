package main

import (
	"testing"
	"time"
)

func TestRelativeAge(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		msEpoch  int64
		expected string
	}{
		{"just now", now.UnixMilli(), "<1m"},
		{"5 minutes ago", now.Add(-5 * time.Minute).UnixMilli(), "5m"},
		{"2 hours ago", now.Add(-2 * time.Hour).UnixMilli(), "2h"},
		{"3 days ago", now.Add(-3 * 24 * time.Hour).UnixMilli(), "3d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RelativeAge(tt.msEpoch)
			if got != tt.expected {
				t.Errorf("RelativeAge(%d) = %q, want %q", tt.msEpoch, got, tt.expected)
			}
		})
	}
}
