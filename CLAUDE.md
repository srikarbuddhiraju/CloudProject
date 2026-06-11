# Project Brief — Azure Landing Zone Conformance Dashboard (name TBD)

> This file is the canonical project context for Claude Code. It captures the full
> ideation and decision history from the founding conversation (June 2026).
> Read this before doing any work in this repo.

## 1. What this project is

A **self-hosted, read-only, fair-source dashboard for Azure Landing Zone governance**,
built for platform engineers, cloud architects, and consultants.

The user points the tool at their Azure tenant using a **Reader-only service principal**.
The tool performs a discovery sweep (management group hierarchy, subscriptions, policy
assignments, RBAC, network topology, logging configuration, cost data) and evaluates the
tenant against the **Azure Landing Zone (ALZ) / Cloud Adoption Framework reference
architecture**, producing:

1. A **live architecture map** of the tenant
2. A **conformance score with specific, opinionated findings**
   (e.g. "identity subscription missing", "no DDoS policy at platform scope",
   "spoke peered directly to spoke")
3. **Inventory and cost views** that hang off the architecture map, not separate tabs

Closest prior art: **AzGovViz** (beloved but ageing PowerShell script producing static
HTML). This project is the live, modern, opinionated, self-hosted successor.

**Positioning:** opinionated insights, not just graphs. The tool supplements consultants
and architects — it tells you what deviates and why it matters, not just what exists.

## 2. Product principles (non-negotiable)

- **Self-hosted & privacy-first.** No SaaS backend, no agents, no telemetry phoning home.
  Single container / single binary. `docker run` and done.
- **Read-only by design.** Reader role only. This is the trust story — never weaken it.
- **Dev-friendly.** CLI-first engine, config-as-code, JSON output, usable in CI pipelines.
- **Dirt-simple for small teams.** Zero-friction setup is a feature, not polish.
- **Opinionated.** Conformance findings over raw data dumps.

## 3. Architecture (locked)

Pipeline of four decoupled stages:

```
collectors  →  embedded store  →  rules engine  →  outputs
(per cloud)     (SQLite            (evaluates       (CLI report,
                 snapshots)         YAML rule        web UI,
                                    packs)           JSON)
```

- **Collectors:** talk to a cloud's APIs, emit normalized inventory. Azure first
  (Azure Resource Graph via KQL + ARM REST + Cost Management API). Future clouds
  (AWS, GCP) are new collectors, not rewrites.
- **Store:** embedded SQLite snapshots. Snapshot history enables **drift over time**
  ("what changed in my tenant since last month") as a near-free future feature.
- **Rules engine:** evaluates **rule packs defined as YAML data**, not code.
  Community contributors can write rules without touching Go.
- **Outputs:** CLI conformance report, JSON (for CI: fail a pipeline on conformance
  regression), and a thin web UI reading JSON from the engine.

**First pack: ALZ conformance** (stable, slow-moving APIs, publicly documented
reference architecture, near-zero maintenance burden).
**Second pack (later, post-traction): AI Landing Zone** — Azure OpenAI deployments,
model versions, token spend per team, quota utilization, private endpoint compliance.
Deliberately deferred: Azure's AI surface churns monthly; bad fit for a solo
maintainer's first wedge, but it is the strongest future marketing hook.
Future packs: security baseline, FinOps, AWS LZ / GCP foundation.

## 4. Stack (locked, with rationale)

- **Backend: Go.** Reasons: single static binary / ~20MB container distribution
  (the Gitea/Caddy/Terraform model the target audience trusts); famously low
  dependency churn (minimum long-term maintenance); boring, explicit, readable
  code — important because the maintainer reviews code with cloud-engineer
  instincts, not developer instincts.
- **Frontend: React + TypeScript**, kept as a thin layer over JSON served by the
  Go backend, embedded into the binary.
- Backend workload is narrow: authenticate SP → query Resource Graph (KQL) and
  ARM/Cost APIs → normalize → evaluate rules → serve JSON. No deep SDK surface needed.

## 5. v0.1 scope (locked)

- **Surface:** BOTH — CLI engine (conformance report + JSON output)
  AND a web view. This forces engine/UI decoupling from day one.
- **Visualization: FULL interactive architecture map IS in v0.1** (locked,
  delay accepted). Management group tree, subscriptions, hub-spoke network
  topology with peering lines, conformance status overlaid on the map.
  This is the hardest single piece of UI in the product and the primary
  demo/marketing asset. Use **React Flow** for the graph rendering — this is
  the one sanctioned exception to the minimal-dependency rule; do NOT
  hand-roll graph layout. Must handle large tenants (50+ subscriptions)
  without degrading into visual spaghetti.
- **Cost data: IN for v0.1** (one extra read-only Cost Management API call
  against the same subscription set).
- Azure only. ALZ conformance pack only.
- Read-only service principal setup documented as the very first user step.

## 6. Testing strategy (no access to enterprise environments)

This is read-only **metadata tooling** — no running workloads needed.

1. **Free personal Azure tenant** as the live integration environment.
   Management groups, policy assignments, RBAC, VNets, NSGs, tags cost ₹0.
   Deploy the official **ALZ accelerator (Bicep/Terraform)** into it for a
   textbook landing zone skeleton with zero compute.
2. **Fixture-based testing for ~90% of the suite.** Record real API responses
   once, sanitize, commit as JSON fixtures, replay in CI. Hand-craft fixtures
   for enterprise scenarios impossible to build for real (50 subscriptions,
   broken hub-spoke, policy conflicts).
3. **Community closes the gap** post-launch: early adopters run it read-only
   against real tenants and report breakage. The Reader-only SP makes this an
   easy ask.

### 6.1 Concrete steps (zero spend — Srikar has no budget for this project)

Everything below uses constructs that are free indefinitely in Azure
(management groups, policy definitions/assignments, RBAC role assignments,
empty VNets/NSGs, resource groups, tags). No paid compute, storage, or
networking is deployed at any point.

1. **Personal Azure tenant.** Sign up for an Azure free account with a
   personal Microsoft account — this makes Srikar Global Admin of his own
   Entra ID tenant. Note: signup requires a card for identity verification,
   but nothing is charged unless the subscription is manually upgraded to
   pay-as-you-go. Stay on the free trial / free-services tier — never add a
   payment method for upgrade.
2. **Deploy the ALZ accelerator in "platform-only" mode** (Bicep or
   Terraform, official Microsoft repo) — management group hierarchy
   (Platform/Landing Zones/Sandbox/Decommissioned), policy assignments at
   each scope, standard subscription placement. No networking/compute
   modules. This is the textbook reference structure the conformance pack
   checks against, built for free.
3. **Create the Reader-only SP** in this tenant exactly as an end user
   would (`az ad sp create-for-rbac`, Reader role at tenant root management
   group). Doubles as the live test of the onboarding docs.
4. **Record → sanitize → fixture.** Run the collector against this tenant
   once, capture raw Resource Graph/ARM/Cost Management JSON responses,
   strip tenant/subscription IDs and names, commit as the "golden" fixture
   set.
5. **Hand-craft "broken" and "at-scale" fixtures** (50 subscriptions,
   spoke-to-spoke peering, missing DDoS policy, policy conflicts) by
   mutating copies of the golden fixtures — same schema as real data, but
   covering scenarios impossible to build solo.

This sub-plan can change as the project takes shape — captured here as the
current best approach, not a commitment.

## 7. Licensing & repo decisions (locked)

- **License: PolyForm Small Business.** Free for individuals and small companies;
  larger commercial users need a paid license. This is **fair source, NOT open
  source (OSI)** — never describe it as "open source" in the README; say
  "source-available / fair source". Rationale: a license can always be relaxed
  later, almost never tightened once outside contributions exist.
- **DCO (or lightweight CLA) required before merging the FIRST external PR** —
  this preserves the legal freedom to relicense later. Do not skip this.
- **Public repo from day one. Build in the open.**
- **Name: TBD.** Deliberately parked until the product's shape is concrete.
  Use a working codename in the meantime.

## 8. About the maintainer (important for how Claude Code should work)

- Srikar is a **cloud platform engineer with deep Azure/platform expertise but
  no development fluency** (can read PowerShell at best). He builds AI-first.
- His scarce asset is **domain judgment**: which conformance checks matter, what
  a well-formed ALZ looks like, what a platform engineer needs to see first.
  Defer to him on domain decisions; carry the engineering.
- **Working norms for Claude Code:**
  - Explain code decisions in cloud-engineering terms, not developer jargon.
  - Prefer boring, explicit, readable code over clever abstractions.
  - Minimize dependencies aggressively — every dep is a maintenance liability
    he cannot personally debug.
  - Keep rule definitions in YAML so he can author conformance rules directly.
  - Small, reviewable increments; explain what each change does and why.
- **Rhythm:** a few focused hours daily. Plan increments accordingly.

## 9. Long-term roadmap (directional, not committed)

1. v0.1 — Azure ALZ conformance: CLI + web UI with full interactive
   architecture map + cost view. (Timeline expectation: months, not weeks,
   at a few-hours-daily rhythm — the map decision trades speed for impact.)
2. v0.2+ — drift-over-time views (snapshot diffing), CI mode hardening,
   exportable reports (Markdown/PDF for auditors)
3. AI LZ rule pack (post-traction)
4. Security baseline / FinOps packs; community rule contributions
5. AWS / GCP collectors → genuinely multi-cloud
6. Commercial path: paid licenses for larger orgs (PolyForm), possibly hosted
   or MSP multi-tenant offering — the self-hosted tool stays free for small users

## 10. Open questions (deliberately unresolved)

- Product name and visual identity
- Exact PolyForm size threshold and pricing
- Initial set of ALZ conformance rules for v0.1 (Srikar to drive selection —
  this is the domain-judgment core of the product)
- First-four-weeks build plan (to be drafted at repo creation)

## 11. Locked technical decisions (June 2026 follow-up)

- **Codename:** `CloudProject` (working folder/module name, until Section 7's
  product name is decided).
- **SQLite driver:** pure-Go (`modernc.org/sqlite`), not CGO-based
  (`mattn/go-sqlite3`). No C toolchain needed to build; cross-compilation stays
  trivial. Matches the minimal-dependency, single-binary principle.
- **Repo layout:** monorepo. `/backend` (Go) and `/frontend` (React/TS), with
  the frontend build embedded into the Go binary via `go:embed`. One repo, one
  version, one release artifact.
- **Auth (Azure):** use the Azure SDK's `DefaultAzureCredential` chain
  (covers env vars, `az login` session, managed identity, etc. automatically).
  For future multi-cloud collectors, define a small `Authenticator` interface
  per cloud so each collector owns its native credential chain — keeps the
  core engine cloud-agnostic.
- **Snapshot storage:** SQLite file defaults to `./<codename>.db` in the
  current working directory. Visible, easy to inspect/back up/`.gitignore`,
  works naturally in CI.

## 12. Working conventions & session workflow

Mirrors the conventions used in Srikar's other repos (RacePhysiX, Panchangam).

### Docs structure (`docs/`)
- `ConvoQA.md` — ongoing decisions and open questions from this point forward.
  (Founding decisions are Sections 1–11 of this CLAUDE.md — don't re-litigate them.)
- `LatestTask.md` — rolling session log, 200-line limit, trim oldest entries
  when exceeded.
- `lessons.md` — mistakes made and rules to avoid repeating, updated per the
  Self-Improvement Loop below.

### Session start checklist
Before doing anything else each session:
1. Read `docs/ConvoQA.md`
2. Read `docs/lessons.md`
3. Read `docs/LatestTask.md`

Then summarise: what was in progress last session, any open `[ ]` items,
relevant lessons — then ask what to work on.

### Self-improvement loop
- After ANY correction from Srikar, update `docs/lessons.md` immediately with
  the pattern and a rule that prevents repeating it.
- Review lessons at session start (above).

### Verification before done
- Never mark a task complete without proving it works: build passes, fixture
  tests pass, and (where relevant) a manual check against the real/fixture
  tenant from Section 6.1.
- Explain code decisions in cloud-engineering terms, per Section 8.

### Hallucination mitigation (important — Srikar can't review Go/SDK code himself)
Azure SDK for Go, ARM API versions, Resource Graph KQL syntax, and Cost
Management API shapes are exactly the kind of detail models misremember from
training data, and Srikar has no way to catch a wrong field/function name in
review.

- Before using any Azure SDK package, function, or struct field not already
  used elsewhere in this repo: run `go doc <package>` /
  `go doc <package>.<Type>` against the **locally downloaded module** (after
  `go get`) and confirm the real signature before writing the call. Never
  write an Azure SDK call from memory alone.
- For ARM API versions, Resource Graph table/column names, or Cost Management
  response shapes: verify against official Microsoft Learn docs, or against
  one real recorded fixture (Section 6.1) — never invent a field name.
- If verification genuinely isn't possible, say so explicitly, mark the code
  with a comment flagging it as unverified, and note it in `docs/LatestTask.md`
  so it gets a second look before merge.

### Git branching
- Feature branches: `feature/<short-description>` — never commit directly to
  `main`.
- Merge only when the session's verification items are all checked.

### File size
- 200-line limit on markdown docs in `docs/` — split into a new file when
  exceeded.
