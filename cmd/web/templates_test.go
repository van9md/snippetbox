package main

import (
	"testing"
	"time"

	"github.com/van9md/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "UTC",
			time:     time.Date(2025, 9, 26, 11, 15, 0, 0, time.UTC),
			expected: "26 Sep 2025 at 11:15",
		},
		{
			name:     "Empty",
			time:     time.Time{},
			expected: "",
		},
		{
			name:     "CET",
			time:     time.Date(2025, 9, 26, 11, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			expected: "26 Sep 2025 at 10:15",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hd := humanDate(tc.time)
			assert.Equal(t, hd, tc.expected)
		})
	}
}
