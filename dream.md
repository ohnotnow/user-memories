# Dream mode

You've been asked to "dream" over the user-memories store — a housekeeping
pass across every memory you've saved about this user. The goal is a
tidier, sharper memory, not a rewrite. Be conservative, and do not
delete anything the user hasn't signed off on.

## 1. Survey

Call `list` with a generous limit (e.g. 200) so you can see the whole
store in one go. If you have more than 200, page through with follow-up
`list` calls adjusting the limit.

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

## 3. Propose, don't act

Report back as a plan, grouped by category above. For each item give:

- the memory id(s) involved
- a one-line reason
- the concrete action you'd take (delete / merge / rewrite / leave)

Only call `delete` or `remember` once the user has signed off on the
plan. If they approve the whole plan in one go, that's fine — but the
default is confirm-first.

## 4. Guiding principles

- This store is a long-term record of who the user is and how they
  like to work. Tidying is good; forgetting isn't.
- When merging, preserve the *why* behind a memory, not just the rule.
- If in doubt about whether something still applies, ask — don't guess.
- Don't invent new memories during a dream pass. Consolidation only.
