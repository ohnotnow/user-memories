# user-memories

A small MCP server that gives Claude a global, cross-project memory store, backed by SQLite.

## What it does

Claude's built-in memory is scoped to a single project. Anything you tell it in one repo doesn't follow you to the next, which gets a bit tedious when you keep re-explaining the same preferences. This MCP server bolts a second, global layer on top: one SQLite file (under your OS config directory by default) that any Claude session can read from and write to.

It exposes five tools: `remember`, `search`, `list`, `delete` and `dream`. Claude can use them to keep hold of things worth remembering across projects, like the fact that you write British English, or that you prefer `uv` over `pip`, or that you really don't want another apology when it makes a mistake.

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

`-s user` registers it at user scope so every project gets it. Run `/mcp` inside Claude Code and you should see it listed with its five tools.

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
| `remember(content)`     | Store a new global memory.                                                                      |
| `search(query, limit?)` | Substring match against stored memories (case-insensitive for ASCII), newest first. Default limit 20. |
| `list(limit?)`          | List memories, newest first. Default limit 20.                                                  |
| `delete(id)`            | Remove a memory by id.                                                                          |
| `dream()`               | Return housekeeping instructions for Claude to tidy up the memory store (see [Dream mode](#dream-mode)). |

## CLI usage

The same binary that serves MCP over stdio also runs as a regular command-line tool. With no subcommand it stays in MCP mode (so `claude mcp add` keeps working); add a subcommand and it'll act on the store directly:

```bash
user-memories list                        # newest 20
user-memories list --limit 100
user-memories search glasgow
user-memories remember "uses British English spelling"
echo "piped content works too" | user-memories remember
user-memories delete 42
user-memories dream                       # prints the dream-mode instructions
user-memories help
```

All subcommands accept `--db PATH` if you want to point at a non-default store. `list` and `search` accept `--limit N`.

## Dream mode

Anthropic have been teasing a "dream" idea for Claude Code — pausing between sessions to look back over what's been stored and tidy it up. There's no new code needed for that here; it's really a prompt. The `dream` tool (MCP) and `user-memories dream` subcommand (CLI) both return a short set of instructions asking Claude to:

1. `list` the whole store,
2. look for duplicates, contradictions, stale entries, fragments that would be stronger as one memory, and things that really belong in a project-specific `CLAUDE.md`,
3. propose a plan back to you before deleting or merging anything.

From inside a Claude Code session you can kick it off with something like:

> Run the user-memories `dream` tool and then follow the instructions it returns.

Or from the terminal, if you'd rather pipe the prompt in yourself:

```bash
user-memories dream | pbcopy
```

## Getting claude to use it

As the MCP server will be 'deferred' claude only gets to know the name of the tool, not what it does.  So adding something like this to your global ~/.claude/CLAUDE.md can help:

```
## User memories

For cross-project lessons (working style, recurring preferences, things that apply regardless of project, interesting tidbit's about the user), use the user-memories MCP.
The built-in auto-memory at ~/.claude/projects/<dir>/memory/ is project-scoped — reserve it for facts specific to one codebase.

It offers :

- `remember(content)` -- Store a new memory
- `search(query, limit?)` -- Case insensitive search for memories
- `list(limit?)` -- List memories, newest first
- `delete(id)` -- Remove a memory
- `dream()` -- Fetch housekeeping instructions for tidying the store

Before calling remember, run a quick search for the topic — avoids writing a duplicate or a contradictory version of something already there.

Also search when:
- the user references prior context ("like last time", "as I mentioned", "remember that...")
- you're about to make a judgement call about their preferences in an area you haven't discussed this session (e.g. attribution, testing style, PR sizing)

Don't list at session start or search speculatively — the store will grow, and searching every turn is noise.

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
