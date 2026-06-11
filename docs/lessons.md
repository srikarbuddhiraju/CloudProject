# Claude Lessons — CloudProject

Running log of mistakes and rules to avoid repeating.
Updated after every correction per CLAUDE.md Section 12 (Self-Improvement Loop).

---

## Hallucination Mitigation

### Never write Azure SDK calls from memory alone
- **Rule**: Before using any Azure SDK for Go package, function, or struct
  field not already used elsewhere in this repo, run `go doc <package>` /
  `go doc <package>.<Type>` against the locally downloaded module and confirm
  the real signature first. For ARM API versions, Resource Graph KQL
  table/column names, or Cost Management response shapes, verify against
  Microsoft Learn docs or a real recorded fixture (Section 6.1).
- **Why**: Srikar has deep Azure domain knowledge but no Go/SDK fluency — he
  cannot catch a wrong field or function name in review. The SDK surface and
  API versions also evolve, so training-data recall is unreliable.
- **How to apply**: Any time a new SDK type/function is introduced, show the
  `go doc` output (or doc link/fixture excerpt) alongside the code that uses
  it. If verification isn't possible, mark the code as unverified and log it
  in `docs/LatestTask.md`.

---

## Git Branching

### Create the feature branch FIRST, before writing any code
- **Rule**: Run `git checkout -b feature/<short-description>` before creating
  or editing any files for that piece of work — not after the work is done
  and committed.
- **Why**: Session "let's begin" (2026-06-11) — the Azure auth check
  (`internal/azure`, `internal/cli/auth.go`) was implemented and committed
  directly on `main`, violating Section 12's git branching rule. Caught
  before push, fixed by moving the commit to a branch retroactively
  (`git branch <name>`, `git reset --hard HEAD~1`, `git checkout <name>`) —
  but only worked because nothing had been pushed yet. Don't rely on that.
- **How to apply**: The very first action of any new piece of work (after
  reading session-start docs) is checking out the feature branch.

## Token Efficiency

### Do not re-run verification commands that already succeeded
- **Rule**: If a command (build, vet, test) already passed in this session
  and nothing relevant has changed since, don't run it again "just to be
  sure."
- **Why**: Same session — re-ran `go vet ./... && go build ./... && go mod
  tidy` immediately after it had already passed (the only failure was a stale
  binary, which was already fixed and verified separately). Caused
  unnecessary tool-approval friction for Srikar.
- **How to apply**: Track what's already been verified in-session. Only
  re-verify after a relevant file changes.

### Don't over-apply the hallucination-mitigation rule to stable, well-known facts
- **Rule**: The "verify against go doc / docs / fixtures" rule (above) is for
  SDK function/struct/field names that drift between versions — not for
  long-stable, universally-documented constants (e.g. the ARM OAuth scope
  `https://management.azure.com/.default`).
- **Why**: Same session — was about to grep SDK source to "verify" a constant
  that's been stable for years and is documented everywhere. Use judgment:
  verify what's actually likely to have changed or be misremembered.

---

## Engineering Discipline

### No untested assumptions
- **Rule**: Any assumption that drives a code change (API field name, KQL
  syntax, policy definition ID, ALZ structure detail) must be verified before
  the change is written — not after.
- **Why**: Carried over from RacePhysiX (Session 16) — an unverified
  assumption became the foundation of a large change and was wrong, wasting
  the whole implementation.
- **How to apply**: If verification takes 5 minutes and the implementation
  takes hours, do the 5 minutes first.
