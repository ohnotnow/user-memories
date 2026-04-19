package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var version = "dev"

func main() {
	args := os.Args[1:]

	// Intercept help requests before anything else so `--help` / `-h` always
	// show the full CLI usage, not just the MCP flagset's terse output.
	if wantsHelp(args) {
		printHelp(os.Stdout)
		return
	}

	cmd, rest := splitSubcommand(args)

	switch cmd {
	case "", "mcp":
		// No subcommand (or explicit `mcp`) → run as an MCP stdio server.
		// Keeps backward compatibility for `claude mcp add user-memories <path>`.
		os.Exit(runMCP(rest))
	case "list", "search", "remember", "delete", "dream":
		os.Exit(runCLI(cmd, rest))
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", cmd)
		printHelp(os.Stderr)
		os.Exit(2)
	}
}

func wantsHelp(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" || a == "help" {
			return true
		}
	}
	return false
}

// splitSubcommand peels the first non-flag argument off as a subcommand.
// Anything else (flags and their values) stays with the remainder so
// global-style flags like `--db /path` placed either before or after the
// subcommand still work.
func splitSubcommand(args []string) (string, []string) {
	// Flags that consume the next arg as their value. Keeps us from
	// mis-identifying "/some/path" as a subcommand in `--db /some/path`.
	valueFlags := map[string]bool{"db": true, "limit": true}

	rest := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "-") {
			rest = append(rest, a)
			name, _, hasValue := splitFlag(a)
			if !hasValue && valueFlags[name] && i+1 < len(args) {
				rest = append(rest, args[i+1])
				i++
			}
			continue
		}
		rest = append(rest, args[i+1:]...)
		return a, rest
	}
	return "", rest
}

func runMCP(args []string) int {
	fs := flag.NewFlagSet("mcp", flag.ExitOnError)
	defaultPath, defaultErr := defaultDBPath()
	dbPath := fs.String("db", defaultPath, "path to SQLite database")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *dbPath == "" {
		fmt.Fprintf(os.Stderr, "could not determine default db path (%v); pass --db\n", defaultErr)
		return 1
	}

	store, err := NewStore(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open store: %v\n", err)
		return 1
	}
	defer store.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	server := mcp.NewServer(&mcp.Implementation{Name: "user-memories", Version: version}, nil)
	registerTools(server, store)

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "server: %v\n", err)
		return 1
	}
	return 0
}

func printHelp(w *os.File) {
	fmt.Fprintf(w, `user-memories — global memory store for Claude, usable as MCP or CLI.

Usage:
  user-memories [--db PATH]                  run as an MCP stdio server (default)
  user-memories mcp [--db PATH]              same, explicit
  user-memories list [--db PATH] [--limit N] newest-first list
  user-memories search [--limit N] QUERY     substring search, newest first
  user-memories remember CONTENT             store a new memory (or pipe via stdin)
  user-memories delete ID                    delete by id
  user-memories dream                        print the dream-mode instructions
  user-memories help                         show this message

Flags:
  --db PATH    path to the SQLite file (default: OS config dir)
  --limit N    cap results (default 20)
`)
}
