# Security Orchestrator (minimal skeleton)

A minimal, container‑based skeleton for an authorized security‑testing orchestration platform.
Replace the stub implementations with the real logic (scope/RoE checks, LLM router, tool adapters,
approval workflow, reporter, etc.) to obtain a full‑featured product.

## Quick start

1. `docker compose build`
2. `docker compose up -d`
3. Create an admin user:  
   `docker exec -it orchestrator orchestrator cli admin create --username admin --email admin@example.com --password "Secret!2025" --role admin`
4. Initialise an engagement:  
   `docker exec -i orchestrator orchestrator cli init -n "Test" -d "Test engagement"`
5. Edit `authorizations/<uuid>/auth.yml` (scope, RoE, contacts, attach auth proof).
6. Run the workflow (passive → active → LLM triage → manual validation → report) – either via the CLI (`orchestrator cli …`) or via the UI at `http://localhost:8080`.
