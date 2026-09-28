# LIS Khanza Mapper

Standalone Go web app to map **LIS tests** (`lis_tests`) to **SIMRS Khanza** lab templates (`template_laboratorium.id_template`). Uses the same MySQL database as SIMRS.


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
| `MEDQLAB_WEBHOOK_API_KEY` | For webhook | `X-API-Key` for `POST /api/v1/medqlab/hasil` |
| `MEDQLAB_BRIDGING_NIP` | For webhook | NIP stored on `periksa_lab.nip` |
| `APP_LISTEN` | No | Default `:8080` (use `:8080` in Docker) |
| `APP_PORT` | No | Host port published by Compose (default `8080`) |
| `MYSQL_PORT` | No | Host port for bundled MariaDB (default `3307`) |

## Health

- `GET /healthz` — no auth
- `GET /readyz` — DB ping, no auth

## API

JSON under `/api/v1` (Basic auth). See spec for endpoints: LIS tests, SIMRS templates/panels, bulk mappings.

### MedQLab hasil webhook (automated bridging)

MedQLab pushes validated results here (no SIMRS button click). Auth is **API key**, not Basic auth.

```http
POST /api/v1/medqlab/hasil
Content-Type: application/json
X-API-Key: <MEDQLAB_WEBHOOK_API_KEY>
```

Body: MedQLab getResult-shaped JSON (`response` + `metaData`) — see project `MEDQLAB.md`.

On success the service:

1. Resolves MedQLab truncated `noOrder` + `visitNumber` (`no_rawat`) → full SIMRS `permintaan_lab.noorder`
   Example: MedQLab `2602250001` + visit `2026/02/25/000001` → SIMRS `PK202602250001`
   (`SUBSTRING(noorder,5,10) = noOrder AND no_rawat = visitNumber`)
2. Flattens / sorts leaf examinations
3. Maps `testId` via `lis_mapping_tests` (scoped to order panels)
4. Writes `periksa_lab`, `detail_periksa_lab`, `saran_kesan_lab`
5. Does **not** post accounting journal

Requires `MEDQLAB_WEBHOOK_API_KEY` and `MEDQLAB_BRIDGING_NIP` in `.env`. Audit rows go to `lis_hasil_inbox`.

## Tables

Created on app startup (if missing):

- `lis_tests` — LIS `testId` in `lis_test_id` (unique)
- `lis_mapping_tests` — links `lis_tests.id` → `id_template`, with `kd_jenis_prw` from `template_laboratorium`
- `lis_hasil_inbox` — webhook audit (payload + post status)

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
