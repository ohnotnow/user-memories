# Dream mode

You've been asked to "dream" over the user-memories store — a housekeeping
pass across every memory you've saved about this user. The goal is a
tidier, sharper memory, not a rewrite. Be conservative, and do not
delete anything the user hasn't signed off on.

## 1. Survey

Call `index` first — it returns a cheap one-line gist of every memory in
the store, newest first, so you can take in the whole shape at once without
pulling full bodies. Then use `get` (or `search`) to fetch the full text of
only the specific candidates you want to inspect closely. Above ~200 memories
`index` shows the newest 200 and notes how many older ones are a `search` away.

Never judge a memory from its gist alone — gists are truncated first lines,
and two memories can share an opening sentence while differing where it
matters. Read the full text of anything you propose to delete or merge.

## 2. Look for

1. **Duplicates** — two or more memories saying the same thing. Keep
   the clearest one, note the others for deletion.
2. **Contradictions** — memories that disagree. Flag them; ask the user
   which is current before changing anything.
3. **Stale entries** — memories tied to a project, deadline or moment
   that's now in the past. Flag, don't assume.
4. **Fragments** — several small memories that would be stronger as one
   richer note. Draft a merged replacement, list the originals it
   would retire.
5. **Project-specific leakage** — memories that really belong in a
   project's `CLAUDE.md` or the per-project auto-memory rather than the
   global store. Flag for the user to move.
6. **Thin memories** — entries so vague they can't usefully guide
   future behaviour. Suggest either enriching or deleting.
7. **Weak gists** — memories whose first line doesn't stand alone as a
   one-sentence summary: it opens with a date, a parenthetical, or a path
   before reaching the point. Since `index` shows only that first line,
   propose a gist-first rewrite of the *opening* — same content, just
   leading with the point. This is rewording, not new content, so it stays
   within the consolidation-only rule below.
8. **PII that slipped through** — real names (the user's, colleagues'),
   email addresses, hostnames, IPs or other identifying details that are
   incidental to the lesson. Propose a rewrite that keeps the content but
   swaps the identity for a neutral reference ("the user", "a colleague",
   "the QA box"). Where the identifying detail *is* the point of the
   memory (e.g. what to call the user, a deliberately recorded hostname),
   flag it and let the user judge.
9. **Misfiled self-lessons** — memories that are really about Claude's
   own behaviour rather than the user (framework edge-cases, debugging
   failure modes; the test: would it still be true with a different
   user?). When the sibling claude-memories store is installed, propose
   moving each one there — `remember()` it in the sibling with its
   content intact, then delete it here, one approved item at a time.
   Without the sibling, leave them where they are: misfiled beats lost.

## 3. Propose, don't act — one category at a time

Remember who you're talking to: the user has never read these memories.
You wrote them; they've at best glimpsed a search result once. A memory
id means nothing to them, and neither does "the note about X" if they
can't see it. Every item you put to them must quote or fully summarise
the memory's content, so a person with zero store context can judge the
call on the spot.

Don't deliver the whole plan in one message — dozens of decisions at
once just gets skipped. Instead:

1. Open with a one-line overview: the categories found and a count for
   each ("stale checks (2), one merge, gist rewrites (3)…").
2. Take ONE category at a time: present its items — content quoted,
   one-line reason, proposed action — and stop. Wait for answers before
   moving on.
3. Apply what's approved as you go, only after sign-off on that item, then
   move to the next category. Rewrites go through `update` — it keeps the
   memory's id and `created_at`, so its history survives the edit. For a
   merge, `update` the memory you're keeping (usually the oldest — its
   `created_at` then reflects when the lesson was really learned) and
   `delete` the ones it retires.

Lead with a recommendation, not a naked question — "I'd merge these
two, they're two halves of one picture; OK?" beats a menu. Approving a
whole category in one go is fine, but never ask the user to answer more
than a handful of questions in a single message.

## 4. Guiding principles

- This store is a long-term record of who the user is and how they
  like to work. Tidying is good; forgetting isn't.
- When merging, preserve the *why* behind a memory, not just the rule.
- If in doubt about whether something still applies, ask — don't guess.
- Don't invent new memories during a dream pass. Consolidation only.
