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
