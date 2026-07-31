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

// sessionStartDreamNudge closes the loop when the index no longer fits the
// recce: Claude suggests a dream session, the store consolidates, and the
// index fits again. Timing is spelled out so a cold model doesn't open the
// session by nagging.
const sessionStartDreamNudge = "This store has outgrown the session-start recce. At a natural stop-point — not mid-task — suggest the user runs a dream session to consolidate it.\n"

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

// sessionStartCharBudget keeps the whole output under Claude Code's 10,000-
// character cap on SessionStart hook stdout — past that the harness diverts
// the output to a file and the index never reaches the session's context.
// The margin below 10,000 leaves room for the overflow note.
const sessionStartCharBudget = 9500

// printSessionStart writes the framed recce: the embedded preamble, then the
// same index lines the `index` subcommand prints, trimmed (oldest first) to
// stay inside sessionStartCharBudget. Trimmed entries fold into the existing
// "older memories not shown — search reaches them" note.
func printSessionStart(w io.Writer, entries []IndexEntry, total int) {
	if len(entries) == 0 {
		fmt.Fprint(w, sessionStartEmpty)
		return
	}
	used := len(sessionStartPreamble) + 1 // +1 for the blank separator line
	kept := 0
	for _, e := range entries {
		line := indexLine(e)
		if used+len(line) > sessionStartCharBudget {
			break
		}
		used += len(line)
		kept++
	}
	if kept == 0 {
		kept = 1 // a lone index line can't bust the budget; an empty index would lie
	}
	fmt.Fprint(w, sessionStartPreamble)
	fmt.Fprintln(w)
	printIndex(w, entries[:kept], total)
	if total > kept {
		fmt.Fprint(w, sessionStartDreamNudge)
	}
}
