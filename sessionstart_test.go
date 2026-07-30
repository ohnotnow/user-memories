package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestSessionStartEmptyStoreIsOneQuietLine(t *testing.T) {
	var buf bytes.Buffer
	printSessionStart(&buf, nil, 0)
	out := buf.String()
	if !strings.Contains(out, "store is empty") {
		t.Errorf("empty store output = %q, want the quiet empty-store line", out)
	}
	if strings.Contains(out, "recce") {
		t.Errorf("empty store output should skip the preamble: %q", out)
	}
}

func TestSessionStartShowsPreambleThenIndex(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	if _, err := store.Add(ctx, "Synthetic gist line.\nfull detail below the fold"); err != nil {
		t.Fatalf("add: %v", err)
	}
	memories, err := store.List(ctx, indexCap)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	var buf bytes.Buffer
	printSessionStart(&buf, buildIndex(memories), 1)
	out := buf.String()

	preambleAt := strings.Index(out, "session-start recce")
	gistAt := strings.Index(out, "Synthetic gist line.")
	if preambleAt < 0 {
		t.Fatalf("output missing preamble heading: %q", out)
	}
	if gistAt < 0 {
		t.Fatalf("output missing index gist: %q", out)
	}
	if gistAt < preambleAt {
		t.Errorf("index should come after the preamble: %q", out)
	}
	if strings.Contains(out, "full detail") {
		t.Errorf("output should stop at each memory's first line: %q", out)
	}
	if !strings.Contains(out, "search") {
		t.Errorf("preamble should mention the search habit: %q", out)
	}
}

func TestSessionStartKeepsOverflowNote(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	if _, err := store.Add(ctx, "only entry"); err != nil {
		t.Fatalf("add: %v", err)
	}
	memories, err := store.List(ctx, indexCap)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	var buf bytes.Buffer
	printSessionStart(&buf, buildIndex(memories), 250)
	if !strings.Contains(buf.String(), "older memories not shown") {
		t.Errorf("expected overflow note, got %q", buf.String())
	}
}

func TestSessionStartStaysUnderHookStdoutCap(t *testing.T) {
	// Claude Code caps a SessionStart hook's stdout at 10,000 characters —
	// build an index big enough to bust that and check the trim holds.
	entries := make([]IndexEntry, 120)
	for i := range entries {
		entries[i] = IndexEntry{
			ID:        int64(i + 1),
			CreatedAt: "2026-01-01 00:00:00",
			Gist:      fmt.Sprintf("Synthetic filler lesson %03d — %s", i+1, strings.Repeat("pad ", 25)),
		}
	}

	var buf bytes.Buffer
	printSessionStart(&buf, entries, len(entries))
	out := buf.String()
	if len(out) > 10000 {
		t.Errorf("output is %d chars, want <= 10000", len(out))
	}
	if !strings.Contains(out, "Synthetic filler lesson 001") {
		t.Errorf("newest entry should survive the trim: %q", out[:200])
	}
	if !strings.Contains(out, "older memories not shown") {
		t.Errorf("expected overflow note after trim, got %q", out)
	}
}
