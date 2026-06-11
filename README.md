# CloudProject (name TBD)

A self-hosted, read-only, fair-source dashboard for **Azure Landing Zone
governance**, built for platform engineers, cloud architects, and
consultants.

Point it at your Azure tenant with a **Reader-only service principal**. It
discovers your management group hierarchy, subscriptions, policy
assignments, RBAC, network topology, logging configuration, and cost data,
then evaluates it against the **Azure Landing Zone (ALZ) / Cloud Adoption
Framework** reference architecture — producing a live architecture map, a
conformance score with specific findings, and inventory/cost views.

Closest prior art: [AzGovViz](https://github.com/JulianHayward/Azure-MG-Sub-Governance-Reporting)
— this is the live, modern, self-hosted successor.

## Status

Pre-build / planning. See `CLAUDE.md` for the full project brief, locked
technical decisions, and roadmap.

## Principles

- Self-hosted, privacy-first — no SaaS backend, no telemetry
- Read-only by design — Reader role only, never more
- `docker run` and done
- Opinionated conformance findings, not just raw data dumps

## License

[PolyForm Small Business License 1.0.0](LICENSE) — source-available
("fair source"), free for individuals and small companies. Larger
commercial users need a paid license.
