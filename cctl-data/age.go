package main

import (
	"fmt"
	"time"
)

func RelativeAge(msEpoch int64) string {
	ts := time.UnixMilli(msEpoch)
	diff := time.Since(ts)
	switch {
	case diff < time.Minute:
		return "<1m"
	case diff < time.Hour:
		return fmt.Sprintf("%dm", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%dh", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%dd", int(diff.Hours()/24))
	default:
		return ts.Format("Jan 02")
	}
}
