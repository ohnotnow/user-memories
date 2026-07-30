package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
)

// sessionStartPreamble frames the index that `session-start` emits. The
// subcommand exists to be called from a Claude Code SessionStart hook, so the
// output is written for a cold model: what the store is, what to do with the
// gists, and the habits that don't survive a fresh context window.
//
//go:embed session-start.md
var sessionStartPreamble string

// sessionStartEmpty is the whole output when there is nothing to recall —
// one quiet line instead of ceremony wrapped around an empty index.
const sessionStartEmpty = "The user-memories store is empty — nothing recorded about this user yet. When a durable preference or working-style fact surfaces this session, store it with the user-memories MCP remember() tool.\n"

func cliSessionStart(ctx context.Context, store *Store) int {
	memories, err := store.List(ctx, indexCap)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	total, err := store.Count(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	printSessionStart(os.Stdout, buildIndex(memories), total)
	return 0
}

// printSessionStart writes the framed recce: the embedded preamble, then the
// same index lines the `index` subcommand prints.
func printSessionStart(w io.Writer, entries []IndexEntry, total int) {
	if len(entries) == 0 {
		fmt.Fprint(w, sessionStartEmpty)
		return
	}
	fmt.Fprint(w, sessionStartPreamble)
	fmt.Fprintln(w)
	printIndex(w, entries, total)
}
