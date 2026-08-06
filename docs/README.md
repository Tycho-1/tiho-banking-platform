# Documentation — Tiho Banking Platform

## This repo (`docs/`)

| Document | Contents |
|----------|----------|
| [architecture.md](architecture.md) | Services, JWT auth, money flows, Kind/GKE deploy options |
| [product-catalog-requirements.md](product-catalog-requirements.md) | Product-catalog (Go) — design spec, acceptance criteria, v1 scope |
| [upstream.md](upstream.md) | Bank of Anthos attribution, version pin, how to sync `src/` |
| [TODO.md](TODO.md) | Done + follow-ups (OTEL/Prometheus, Flux, public platform repo, …) |

Deploy detail (all overlays): **[deploy/README.md](../deploy/README.md)**.

## Per-service notes (upstream, kept in `src/`)

Each microservice folder has its own `README.md` with build and API notes from Google. Start here when reading code:

| Service | README |
|---------|--------|
| Frontend | [src/frontend/README.md](../src/frontend/README.md) |
| User service | [src/accounts/userservice/README.md](../src/accounts/userservice/README.md) |
| Contacts | [src/accounts/contacts/README.md](../src/accounts/contacts/README.md) |
| Balance reader | [src/ledger/balancereader/README.md](../src/ledger/balancereader/README.md) |
| Ledger writer | [src/ledger/ledgerwriter/README.md](../src/ledger/ledgerwriter/README.md) |
| Transaction history | [src/ledger/transactionhistory/README.md](../src/ledger/transactionhistory/README.md) |
| Load generator | [src/loadgenerator/README.md](../src/loadgenerator/README.md) |
| Product catalog (Go, Tiho) | [src/products/README.md](../src/products/README.md) |

## Upstream docs — what to copy vs skip

Do **not** copy the full upstream `docs/` folder. Most of it targets GKE, Cloud Build, Workload Identity, and Cloud Deploy.

| Upstream doc | Use for Tiho? |
|--------------|---------------|
| [development.md](https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/main/docs/development.md) | **Reference only** — Skaffold/GKE dev loop; rewrite for your CI later |
| [environments.md](https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/main/docs/environments.md) | **Reference** — `ENABLE_METRICS=false` / `ENABLE_TRACING=false` for non-GKE |
| [troubleshooting.md](https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/main/docs/troubleshooting.md) | **Reference** — common pod failures |
| [ci-cd-pipeline.md](https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/main/docs/ci-cd-pipeline.md) | **Skip** — GCP fleet pipeline |
| [workload-identity.md](https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/main/docs/workload-identity.md) | **Skip for Kind** — GKE only |
| Architecture screenshots (`docs/img/`) | **Copied** — login, transactions, architecture diagram in root [README.md](../README.md) (Apache-2.0, from upstream) |

## Platform (cluster bootstrap)

This app expects a Kubernetes cluster (usually **Kind**) with overlay-specific add-ons. Cluster install docs are **not** in this public app repo yet (local/private Kind bootstrap is fine). Prerequisites are listed per overlay in [deploy/README.md](../deploy/README.md).
