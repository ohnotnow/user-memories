package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RememberInput struct {
	Content string `json:"content" jsonschema:"the memory to store - a cross-project fact, preference, or note worth recalling in any future conversation"`
}

type RememberOutput struct {
	Memory Memory `json:"memory"`
}

type SearchInput struct {
	Query string `json:"query" jsonschema:"substring to match against memory content (case-insensitive for ASCII characters)"`
	Limit int    `json:"limit,omitempty" jsonschema:"maximum number of results, newest first (default 20)"`
}

type ListInput struct {
	Limit int `json:"limit,omitempty" jsonschema:"maximum number of results, newest first (default 20)"`
}

type MemoriesOutput struct {
	Memories []Memory `json:"memories"`
	Count    int      `json:"count"`
}

type DeleteInput struct {
	ID int64 `json:"id" jsonschema:"id of the memory to delete"`
}

type DeleteOutput struct {
	Deleted bool  `json:"deleted"`
	ID      int64 `json:"id"`
}

type DreamInput struct{}

type DreamOutput struct {
	Instructions string `json:"instructions"`
}

func registerTools(server *mcp.Server, store *Store) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "remember",
		Description: "Store a new global memory that should persist across every project. Use for cross-project facts, preferences, or notes about the user - not project-specific context.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RememberInput) (*mcp.CallToolResult, RememberOutput, error) {
		if in.Content == "" {
			return nil, RememberOutput{}, fmt.Errorf("content is required")
		}
		m, err := store.Add(ctx, in.Content)
		if err != nil {
			return nil, RememberOutput{}, err
		}
		return textResult(fmt.Sprintf("Stored memory %d.", m.ID)), RememberOutput{Memory: *m}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search",
		Description: "Search global memories by substring, newest first.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in SearchInput) (*mcp.CallToolResult, MemoriesOutput, error) {
		memories, err := store.Search(ctx, in.Query, in.Limit)
		if err != nil {
			return nil, MemoriesOutput{}, err
		}
		return textResult(summarise(memories)), MemoriesOutput{Memories: memories, Count: len(memories)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list",
		Description: "List global memories newest first.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ListInput) (*mcp.CallToolResult, MemoriesOutput, error) {
		memories, err := store.List(ctx, in.Limit)
		if err != nil {
			return nil, MemoriesOutput{}, err
		}
		return textResult(summarise(memories)), MemoriesOutput{Memories: memories, Count: len(memories)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete",
		Description: "Delete a global memory by id.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DeleteInput) (*mcp.CallToolResult, DeleteOutput, error) {
		ok, err := store.Delete(ctx, in.ID)
		if err != nil {
			return nil, DeleteOutput{}, err
		}
		if !ok {
			return textResult(fmt.Sprintf("No memory with id %d.", in.ID)), DeleteOutput{Deleted: false, ID: in.ID}, nil
		}
		return textResult(fmt.Sprintf("Deleted memory %d.", in.ID)), DeleteOutput{Deleted: true, ID: in.ID}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "dream",
		Description: "Return the 'dream mode' instructions — a housekeeping pass over the stored memories. Call this and then follow the returned instructions to tidy up duplicates, contradictions, and stale entries.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ DreamInput) (*mcp.CallToolResult, DreamOutput, error) {
		return textResult(dreamInstructions), DreamOutput{Instructions: dreamInstructions}, nil
	})
}

func summarise(memories []Memory) string {
	if len(memories) == 0 {
		return "No memories found."
	}
	out := fmt.Sprintf("Found %d memor%s:\n", len(memories), plural(len(memories)))
	for _, m := range memories {
		out += fmt.Sprintf("  [%d] %s  (%s)\n", m.ID, m.Content, m.CreatedAt)
	}
	return out
}

func plural(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}
