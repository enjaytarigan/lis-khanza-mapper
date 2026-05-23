# LIS Khanza Mapper

Standalone Go web app to map **LIS tests** (`lis_tests`) to **SIMRS Khanza** lab templates (`template_laboratorium.id_template`). Uses the same MySQL database as SIMRS.

Specification: [MAPPING_APP_SPEC.md](https://github.com/your-org/SIMRS-Khanza/blob/main/MAPPING_APP_SPEC.md) (in SIMRS-Khanza repo).

## Quick start (Docker)

```bash
cp .env.example .env
# Edit DATABASE_DSN (SIMRS MySQL), AUTH_USERNAME, AUTH_PASSWORD
./scripts/deploy.sh
```

Open `http://localhost:8080` and sign in with HTTP Basic auth.

**Migrations** (`lis_tests`, `lis_mapping_tests`) run automatically every time the app starts.

## Local development (Docker)

Bundled MariaDB + app (good for trying the mapper without touching host MySQL):

```bash
cp .env.example .env
# Edit AUTH_USERNAME / AUTH_PASSWORD
./scripts/run-local.sh
```

- App: `http://localhost:8080`
- MariaDB: `localhost:3307` (user `mapper`, password `mapper`, database `sik`)

Use your host SIMRS database instead:

```bash
./scripts/run-local.sh --external-db
```

Set `DATABASE_DSN` in `.env` to reach MySQL on the host, e.g. `...@tcp(host.docker.internal:3306)/sik?...`.

## Docker Compose files

| File | Purpose |
|------|---------|
| `docker-compose.yml` | App service; connects via `DATABASE_DSN` in `.env` |
| `docker-compose.local.yml` | Adds MariaDB and default DSN for local stack |

Production / deploy:

```bash
docker compose up -d --build
```

Local stack:

```bash
docker compose -f docker-compose.yml -f docker-compose.local.yml up --build
```

## Configuration

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_DSN` | Yes | MySQL DSN |
| `AUTH_USERNAME` | Yes | HTTP Basic user |
| `AUTH_PASSWORD` | Yes | HTTP Basic password |
| `APP_LISTEN` | No | Default `:8080` (use `:8080` in Docker) |
| `APP_PORT` | No | Host port published by Compose (default `8080`) |
| `MYSQL_PORT` | No | Host port for bundled MariaDB (default `3307`) |

## Health

- `GET /healthz` — no auth
- `GET /readyz` — DB ping, no auth

## API

JSON under `/api/v1` (Basic auth). See spec for endpoints: LIS tests, SIMRS templates/panels, bulk mappings.

## Tables

Created on app startup (if missing):

- `lis_tests` — LIS `testId` in `lis_test_id` (unique)
- `lis_mapping_tests` — links `lis_tests.id` → `id_template`, with `kd_jenis_prw` from `template_laboratorium`

SIMRS tables `template_laboratorium` and `jns_perawatan_lab` are read-only.

## Integration (SIMRS Java)

Example SIMRS integration lookup:

```sql
SELECT m.id_template FROM lis_mapping_tests m
INNER JOIN lis_tests t ON t.id = m.lis_tests_pk
WHERE t.lis_test_id = ? AND t.status='aktif' AND m.status='aktif'
  AND m.kd_jenis_prw = ?
LIMIT 1;
```
