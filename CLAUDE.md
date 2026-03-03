# TrailMates — CLAUDE.md

AI context file for Claude Code. Read this before making any changes.

## What This App Does
TrailMates is a full-stack web app for solo backpackers (18–35) planning multi-city Europe trips.
Users build trips, get an optimized route, browse city attractions, and connect with compatible travelers.

## Tech Stack
- **Backend:** Go 1.22, MySQL, JWT auth (`golang-jwt/jwt/v5`), port 8080
- **Frontend:** React 19 + Vite 7, Tailwind v4 + Shadcn UI, React Router v7, TanStack Query v5
- **Deployment:** GitHub Actions → AWS Lightsail (`ubuntu@3.130.164.2`), Nginx reverse proxy, systemd
- **Domain:** `macicapstone.duckdns.org`

## Key Patterns
- All API calls go through Vite proxy `/api` → `localhost:8080` (prefix stripped)
- JWT stored in `localStorage`; 401 responses redirect to `/login`
- Auth: full user table (email + password_hash), not single admin password
- Backend reads env vars for DB credentials and JWT secret
- Nginx serves frontend static files from `/var/www/app/frontend/dist`
- systemd service `app-backend` runs `/var/www/app/backend/app`

## File Structure
```
backend/
  main.go        # routes, CORS, server startup
  auth.go        # JWT middleware, login/register handlers
  handlers.go    # CRUD handlers for all resources
  database.go    # MySQL connection
  models.go      # structs + FormatSeconds helper

frontend/src/
  pages/         # one file per route
  components/    # shared UI + admin sub-components
  lib/           # api.js, queries.js, format.js, utils.js
  context/       # AuthContext.jsx

docs/
  project-proposal.md
  architecture.md
  database-schema.md
```

## Environment Variables (backend)
```
DB_USER, DB_PASSWORD, DB_HOST, DB_NAME
JWT_SECRET
```

## GitHub Secrets Required
`SSH_PRIVATE_KEY`, `SERVER_IP`, `SERVER_USER`,
`DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_NAME`, `JWT_SECRET`

## What Is NOT Being Built (MVP scope)
- In-app payments or affiliate booking links
- Real-time chat (connection requests only)
- Coverage outside Europe + Morocco weekend add-on
- Live airfare scraping (stored cost level estimates only)

## Development Notes
- `go mod tidy` removes unused deps — add source files referencing a dep BEFORE running tidy
- Shadcn init requires: `@import "tailwindcss"` in CSS, Tailwind vite plugin in config, jsconfig.json with `@/*` alias
- Windows bash: use forward slashes in paths, e.g. `cd "C:/path"`
