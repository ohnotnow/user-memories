# user-memories

A small MCP server that gives Claude a global, cross-project memory store, backed by SQLite. It records what Claude learns about *you*; its sibling [claude-memories](https://github.com/ohnotnow/claude-memories) does the same job for what Claude learns about *itself*.

## What it does

Claude's built-in memory is scoped to a single project. Anything you tell it in one repo doesn't follow you to the next, which gets a bit tedious when you keep re-explaining the same preferences. This MCP server bolts a second, global layer on top: one SQLite file (under your OS config directory by default) that any Claude session can read from and write to.

It exposes eight tools: `remember`, `get`, `search`, `list`, `index`, `update`, `delete` and `dream`. Claude can use them to keep hold of things worth remembering across projects, like the fact that you write British English, or that you prefer `uv` over `pip`, or that you really don't want another apology when it makes a mistake. `index` is the session-start recce — a cheap one-line-per-memory overview so a fresh session can see what it already knows before it starts guessing (see [Session-start recce](#session-start-recce)).

The same binary doubles as a regular CLI, so you can list, search, add or delete memories straight from your terminal without going through Claude.

## Prerequisites

- Go 1.25 or newer (for building from source)
- [Claude Code](https://docs.claude.com/en/docs/claude-code), or any other MCP-capable client

## Getting started

### Install

The quickest option is `go install`:

```bash
go install github.com/ohnotnow/user-memories@latest
```

That drops the binary at `$(go env GOPATH)/bin/user-memories`, which is usually `~/go/bin/user-memories`.

### Use a prebuilt binary

If you'd rather not build it yourself, grab one for your platform from the [releases page](https://github.com/ohnotnow/user-memories/releases). Binaries are named `user-memories-<os>-<arch>`, so pick the one that matches your machine.

On macOS or Linux, make it executable and stash it somewhere on your PATH:

```bash
chmod +x user-memories-darwin-arm64
mv user-memories-darwin-arm64 /usr/local/bin/user-memories
```

The macOS binary isn't signed, so Gatekeeper will block it the first time you try to run it. Right-click the file in Finder, choose Open, and it'll stop complaining from then on.

On Windows, rename `user-memories-windows-amd64.exe` to something friendlier like `user-memories.exe` and drop it somewhere on your PATH.

### Register with Claude Code

```bash
claude mcp add -s user user-memories ~/go/bin/user-memories
```

Swap `~/go/bin/user-memories` for the actual path if you downloaded the binary instead.

`-s user` registers it at user scope so every project gets it. Run `/mcp` inside Claude Code and you should see it listed with its eight tools.

### Database location

The SQLite file lives in your OS's standard config directory (whatever Go's `os.UserConfigDir()` returns):

| OS      | Path                                                   |
| ------- | ------------------------------------------------------ |
| macOS   | `~/Library/Application Support/user-memories/memories.db` |
| Linux   | `~/.config/user-memories/memories.db`                  |
| Windows | `%AppData%\user-memories\memories.db`                  |

Pass `--db /path/to/custom.db` if you'd like it somewhere else.

## Tools

| Tool                    | Description                                                                                     |
| ----------------------- | ----------------------------------------------------------------------------------------------- |
| `remember(content)`     | Store a new global memory. Start it with a one-sentence gist — the first line is what shows in `index()`. |
| `get(id)`               | Fetch one memory in full by id — the follow-up read when `index()` shows a promising gist.      |
| `search(query, limit?)` | All-words match: every word in the query must appear in a memory, in any order (case-insensitive for ASCII). Newest first, default limit 20. |
| `list(limit?)`          | List memories in full, newest first. Default limit 20.                                          |
| `index()`               | One line (id, date, gist) per memory across the whole store, newest first — the cheap session-start recce. |
| `update(id, content)`   | Rewrite a memory in place, keeping its id and `created_at` (`updated_at` records the rewrite).  |
| `delete(id)`            | Remove a memory by id.                                                                          |
| `dream()`               | Return housekeeping instructions for Claude to tidy up the memory store (see [Dream mode](#dream-mode)). |

## CLI usage

The same binary that serves MCP over stdio also runs as a regular command-line tool. With no subcommand it stays in MCP mode (so `claude mcp add` keeps working); add a subcommand and it'll act on the store directly:

```bash
user-memories list                        # newest 20
user-memories list --limit 100
user-memories search glasgow
user-memories search moat gitlab          # all words must appear, any order
user-memories index                       # one line per memory, whole store
user-memories session-start               # the index framed for a fresh session (see Session-start hook)
user-memories show 42                     # one memory in full
user-memories remember "uses British English spelling"
echo "piped content works too" | user-memories remember
user-memories delete 42
user-memories dream                       # prints the dream-mode instructions
user-memories help
```

All subcommands accept `--db PATH` if you want to point at a non-default store. `list` and `search` accept `--limit N`.

## Session-start recce

The store is pull-based: memories only surface when Claude decides to read them. Left to chance that reliably fails — a fresh session doesn't know what it doesn't know, so it re-derives things you told it weeks ago. The fix is a cheap "have a nosey" at the start of each session.

`index()` returns one line per memory across the whole store — `id`, date, and a gist (the memory's first line, truncated to ~120 characters, counted in runes so it never mangles emoji or accents). At roughly 2–3k tokens for 100 memories it's cheap enough to call once per session: skim the gists, then `get()` or `search()` for the full text of anything relevant to the day's work. Above ~200 memories `index()` returns the newest 200 and notes how many older ones are only a `search()` away.

The server ships this workflow to clients through the MCP `instructions` field, so harnesses that surface server instructions (Claude Code does) pick it up automatically. For anything that doesn't, the block in [Getting claude to use it](#getting-claude-to-use-it) puts the same guidance in your own instructions file. It also pays to write memories gist-first — a one-sentence opening line — since that first line is all `index()` shows.

## Session-start hook

Everything above still depends on Claude *choosing* to look, and in practice that's a coin-toss weighted the wrong way — "call `index()` first" is exactly the kind of instruction that loses out to whatever you actually asked for. A [Claude Code SessionStart hook](https://code.claude.com/docs/en/hooks) takes the choice away: the harness runs a command when a session starts and adds its stdout to Claude's context before the first prompt. That's deterministic for you (your preferences are in view before Claude's first reply, in every project, without re-explaining) and for future Claude (the index is simply there, the way `CLAUDE.md` is simply there).

The `session-start` subcommand exists to be that command. It prints the same one-line-per-memory index as `index`, framed by a short preamble telling a cold model what the store is and which habits matter: `get()` the entries relevant to today's work, `search()` before judgement calls about your preferences, `search()` before `remember()`. When the store is empty it prints one quiet line instead of ceremony wrapped around nothing. And it keeps itself under Claude Code's 10,000-character cap on hook stdout (past which output is diverted to a file and never reaches context), trimming oldest-first into the same "older memories not shown — search reaches them" note that `index` uses. Whenever trimming happens, the output ends with a one-line nudge telling Claude to suggest a [dream](#dream-mode) session at a natural stop-point — so an overgrown store prompts its own consolidation instead of silently shedding its oldest memories.

Wire it up in `~/.claude/settings.json` (user scope, so every project gets it):

```json
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "startup|resume|clear",
        "hooks": [
          {
            "type": "command",
            "command": "/absolute/path/to/user-memories session-start"
          }
        ]
      }
    ]
  }
}
```

Use the binary's absolute path — hooks don't reliably inherit your login shell's `PATH`. The matcher above fires on new sessions, resumed ones, and `/clear`; omit the `matcher` field entirely and it also fires after compaction and on forked sessions. If you run the sibling [claude-memories](https://github.com/ohnotnow/claude-memories) too, give it its own entry in the same `hooks` array — each store frames its own index.

The cost is the index itself landing in every session — roughly 2–3k tokens per 100 memories, needed or not. That's the price of never re-explaining yourself; if it stings, [dream mode](#dream-mode) keeps the store lean. You can also run the same command by hand mid-session: type `! user-memories session-start` at the prompt and the output lands straight in context.

## Dream mode

Anthropic have been teasing a "dream" idea for Claude Code — pausing between sessions to look back over what's been stored and tidy it up. There's no new code needed for that here; it's really a prompt. The `dream` tool (MCP) and `user-memories dream` subcommand (CLI) both return a short set of instructions asking Claude to:

1. `index` the whole store for a cheap one-line overview, then pull the full text of anything worth a closer look,
2. look for duplicates, contradictions, stale entries, fragments that would be stronger as one memory, memories whose first line doesn't work as an `index` gist, and things that really belong in a project-specific `CLAUDE.md`,
3. walk you through the plan one category at a time — quoting each memory's content, since you've likely never read them — and only delete, merge or rewrite once you've signed off on that item.

From inside a Claude Code session you can kick it off with something like:

> Run the user-memories `dream` tool and then follow the instructions it returns.

Or from the terminal, if you'd rather pipe the prompt in yourself:

```bash
user-memories dream | pbcopy
```

## Getting claude to use it

The server advertises this workflow itself through the MCP `instructions` field, which Claude Code folds into the session automatically. But the tools are 'deferred' (Claude sees only a tool's name until it looks closer), and not every harness surfaces server instructions — so it's worth putting the workflow in your global `~/.claude/CLAUDE.md` too. (With the [SessionStart hook](#session-start-hook) wired up, the recce itself already happens mechanically — the block below then mostly earns its keep through the mid-session habits.)

```
## User memories

For cross-project lessons (working style, recurring preferences, things that apply regardless of project, interesting tidbits about the user), use the user-memories MCP.
The built-in auto-memory at ~/.claude/projects/<dir>/memory/ is project-scoped — reserve it for facts specific to one codebase.
Lessons about your own behaviour (framework rakes, your recurring failure modes — still true with a different user) belong in the sibling claude-memories MCP if it's installed.

It offers:

- `remember(content)` -- Store a new memory (start it with a one-sentence gist — the first line is what shows in `index`)
- `get(id)` -- Fetch one memory in full by id
- `search(query, limit?)` -- Case-insensitive all-words search (every word must appear, any order)
- `list(limit?)` -- List memories in full, newest first
- `index()` -- One line per memory across the whole store, newest first
- `update(id, content)` -- Rewrite a memory in place (keeps its id and created_at)
- `delete(id)` -- Remove a memory
- `dream()` -- Fetch housekeeping instructions for tidying the store

At the start of a session, call `index()` once to skim every memory's gist, then `get()` or `search()` for the full text of anything relevant to the work at hand. Don't `list()` the whole store speculatively — the index is the cheap overview.

Before calling remember, run a quick search for the topic — avoids writing a duplicate or a contradictory version of something already there. If a memory is wrong or outdated, prefer `update()` over delete + re-remember: it keeps the memory's id and created_at.

Also search when:
- the user references prior context ("like last time", "as I mentioned", "remember that...")
- you're about to make a judgement call about their preferences in an area you haven't discussed this session (e.g. attribution, testing style, PR sizing)

```

## Running tests

```bash
go test ./...
```

Tests run against an in-memory SQLite database, so there's no setup to do.

## Releases

Pushing a tag matching `v*.*.*` (for example `v0.1.0`) kicks off the release workflow at `.github/workflows/release.yml`. It builds binaries for Linux, macOS and Windows across amd64/arm64, generates SHA256 checksums, and attaches the lot to a GitHub release with auto-generated notes.

```bash
git tag v0.1.0
git push origin v0.1.0
```

## Contributing

```bash
git clone git@github.com:ohnotnow/user-memories.git
cd user-memories
go test ./...
```

Then edit, test, open a PR. The project is deliberately tiny, so small changes are very welcome. Please don't send me a PR that turns it into a platform.

## Licence

MIT. See [LICENSE](LICENSE).
