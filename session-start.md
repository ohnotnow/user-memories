# User memories — session-start recce

The index below is the user-memories store: what previous sessions learned about this user and how they like to work, surfaced automatically so you don't re-derive it. Skim the gists now, then use the user-memories MCP tools:

- `get(id)` — pull the full text of anything relevant to today's work; a gist is a truncated first line, not the memory.
- `search(query)` — before a judgement call about the user's preferences in an area you haven't discussed this session (attribution, testing style, PR sizing…), and when they reference prior context ("like last time", "as I mentioned").
- Before storing something new with `remember()`, `search()` the topic first — update an existing memory rather than duplicating it.
