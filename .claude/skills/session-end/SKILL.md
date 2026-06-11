---
name: session-end
description: End-of-session checklist — verify items, update docs, check merge readiness
---

Run the CloudProject end-of-session checklist. Work through each item:

1. **Verification blockers** — list all unchecked `[ ]` items in
   `docs/LatestTask.md`. Any unchecked verification item is a hard blocker —
   branch does not merge.

2. **Hallucination check** — for any new Azure SDK/API usage this session,
   confirm it was verified per `docs/lessons.md` (Hallucination Mitigation).
   Flag anything still unverified.

3. **Line limit** — check all `.md` files in `docs/` for the 200-line rule.
   Run: `wc -l docs/*.md` — split any file over 200 lines immediately.

4. **Update `LatestTask.md`** — add what was done this session, update the
   verification checklist, update the "Next Session" items list.

5. **Update `lessons.md`** — if any corrections were made this session, add
   the pattern now.

6. **Git status** — show uncommitted changes (`git status` + `git diff
   --stat`). Remind: commit docs before ending session.

7. **Merge readiness verdict**:
   - All verification items checked? Yes/No
   - All `.md` files under 200 lines? Yes/No
   - Docs committed? Yes/No
   - → READY TO MERGE or BLOCKED (list reasons)

Do not commit or merge without Srikar's explicit confirmation.
