# Using the user-memories store

These are global memories about the user and how they like to work. They carry across every project — unlike project-scoped notes, which belong in that project's own files.

At the start of a session, call `index()` once. It returns a cheap one-line gist of every memory, newest first. Skim it, then pull the full text of anything relevant to the work at hand — `get(id)` for a specific line from the index, `search(query)` for a topic. Search matches every word in the query (in any order), so a few distinctive words work better than an exact phrase. Don't `list()` the whole store speculatively — the index is the cheap overview; get and search are the targeted reads.

During a session:

- When a new topic goes live, `search()` for it before assuming you know nothing about the user's preferences there — a relevant memory may already exist.
- Before `remember()`, `search()` for the topic first, to avoid storing a duplicate or a contradiction of something already there.
- When you write a memory, start it with a one-sentence gist: its first line is what appears in `index()`.
- If a memory turns out to be wrong or outdated, `update(id, content)` it in place — that keeps its id and `created_at`, so its history stays honest. Reserve `delete()` for memories that no longer apply at all.
- Before you `update()` or `delete()` anything, `get()` its full text first — gists are truncated first lines, and two memories can share an opening sentence while differing where it matters. `update()` replaces the whole body, so fold in anything from the old text still worth keeping.
- Ids survive updates, but another session can still delete or consolidate a memory. If `get()`, `update()` or `delete()` says an id doesn't exist, re-run `index()` or `search()` to find the current copy rather than assuming a fault.

Use these for cross-project facts, preferences, and working style — not for context specific to a single codebase. Lessons about your own behaviour (framework rakes, your recurring failure modes — anything that would still be true with a different user) belong in the sibling claude-memories store when it is installed; without it, keeping them here beats losing them.
