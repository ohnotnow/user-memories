package main

import (
	"bytes"
	"context"
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
