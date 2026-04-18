package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var version = "dev"

func main() {
	defaultPath, defaultErr := defaultDBPath()
	dbPath := flag.String("db", defaultPath, "path to SQLite database")
	flag.Parse()

	if *dbPath == "" {
		fmt.Fprintf(os.Stderr, "could not determine default db path (%v); pass --db\n", defaultErr)
		os.Exit(1)
	}

	store, err := NewStore(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open store: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	server := mcp.NewServer(&mcp.Implementation{Name: "user-memories", Version: version}, nil)
	registerTools(server, store)

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "server: %v\n", err)
		os.Exit(1)
	}
}
