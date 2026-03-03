# Architecture: TrailMates

## Overview

TrailMates is a full-stack web application with a Go REST API backend and a React single-page application frontend, deployed to a single AWS Lightsail instance.

```
Browser
  │
  │ HTTPS (443)
  ▼
Nginx (reverse proxy)
  ├── /* → /var/www/app/frontend/dist  (static files)
  └── /api/* → localhost:8080          (Go backend, prefix stripped)
       │
       ▼
  app-backend (systemd service)
       │
       ▼
  MySQL (localhost:3306)
```

---

## Backend

**Language:** Go 1.22
**Port:** 8080
**Router:** `net/http` standard library (method+path pattern matching)
**Auth:** JWT via `golang-jwt/jwt/v5`; tokens stored client-side in `localStorage`
**Database:** MySQL via `database/sql` + `go-sql-driver/mysql`

### File Layout
```
backend/
  main.go       # server startup, CORS middleware, route registration
  auth.go       # JWT generation, requireAuth middleware, /register, /login, /auth/verify
  handlers.go   # all resource CRUD handlers
  database.go   # ConnectDB(), connection pool setup
  models.go     # structs (User, Profile, Trip, City, Attraction, etc.)
```

### CORS
All routes wrapped in a `cors()` middleware that sets permissive headers for local development. In production, Nginx handles HTTPS and the backend is not exposed directly.

---

## Frontend

**Framework:** React 19 + Vite 7
**Styling:** Tailwind CSS v4 + Shadcn UI components
**Routing:** React Router v7
**Data fetching:** TanStack Query v5
**API calls:** `/api/*` proxied by Vite in dev; hits Nginx in production

### File Layout
```
frontend/src/
  pages/           # one file per route
  components/      # shared layout, navbar, footer, protected route
  context/         # AuthContext.jsx (JWT + user state)
  lib/
    api.js         # apiFetch() wrapper with auth header + 401 redirect
    queries.js     # TanStack Query query/mutation definitions
    format.js      # date and time formatting helpers
    utils.js       # Shadcn cn() utility
```

---

## API Design

### Auth
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/register` | — | Create account (email, password, display_name) |
| POST | `/login` | — | Return JWT |
| GET | `/auth/verify` | JWT | Confirm token is valid |

### Profile
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/profile` | JWT | Get own profile |
| PUT | `/profile` | JWT | Update bio, travel_style, interests, is_visible |

### Cities (public, seeded data)
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/cities` | — | List all cities |
| GET | `/cities/:id` | — | City detail + attractions |

### Attractions (public)
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/cities/:id/attractions` | — | Attractions for a city |
| GET | `/attractions/:id` | — | Single attraction |

### Trips
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/trips` | JWT | List user's trips |
| POST | `/trips` | JWT | Create trip |
| GET | `/trips/:id` | JWT | Trip detail with cities |
| PUT | `/trips/:id` | JWT | Update trip metadata |
| DELETE | `/trips/:id` | JWT | Delete trip |

### Trip Cities
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/trips/:id/cities` | JWT | Add city to trip |
| PUT | `/trips/:id/cities/:cityId` | JWT | Update dates or order |
| DELETE | `/trips/:id/cities/:cityId` | JWT | Remove city from trip |

### Route
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/trips/:id/route` | JWT | Optimized city order for this trip |

### Itinerary
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/trips/:id/itinerary` | JWT | List itinerary items |
| POST | `/trips/:id/itinerary` | JWT | Add attraction to itinerary |
| DELETE | `/trips/:id/itinerary/:itemId` | JWT | Remove item |

### Matching & Connections
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/trips/:id/matches` | JWT | Users with overlapping trips |
| GET | `/connections` | JWT | List own connections |
| POST | `/connections` | JWT | Send connection request |
| PUT | `/connections/:id` | JWT | Accept or decline request |

---

## Example Request / Response

### POST /register
**Request:**
```json
{
  "email": "alex@example.com",
  "password": "hunter2",
  "display_name": "Alex"
}
```
**Response 201:**
```json
{
  "token": "eyJhbGci...",
  "user": { "id": 1, "email": "alex@example.com", "display_name": "Alex" }
}
```

### GET /trips/:id/route
**Response 200:**
```json
{
  "trip_id": 3,
  "optimized_order": [
    { "city_id": 5, "name": "Lisbon", "country": "Portugal", "arrival_date": "2026-06-10", "departure_date": "2026-06-13", "order_index": 1 },
    { "city_id": 2, "name": "Madrid", "country": "Spain",    "arrival_date": "2026-06-13", "departure_date": "2026-06-16", "order_index": 2 },
    { "city_id": 8, "name": "Barcelona", "country": "Spain", "arrival_date": "2026-06-16", "departure_date": "2026-06-19", "order_index": 3 }
  ],
  "hops": [
    { "from": "Lisbon", "to": "Madrid",    "distance_km": 628,  "cost_level": "low" },
    { "from": "Madrid", "to": "Barcelona", "distance_km": 621,  "cost_level": "low" }
  ]
}
```

---

## Deployment Pipeline

1. Push to `main` branch
2. GitHub Actions runner (ubuntu-latest):
   - `npm install && npm run build` (frontend)
   - `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .` (backend)
   - SCP frontend dist + backend binary to Lightsail
   - Restart `app-backend` systemd service
   - Verify `https://macicapstone.duckdns.org/api/health`

---

## Route Optimizer Algorithm

The optimizer uses a **nearest-neighbor greedy algorithm** on geographic coordinates:

1. Start from the city with the earliest arrival date (or leftmost geographically if no dates set).
2. At each step, pick the unvisited city closest by straight-line distance (haversine).
3. Flag any hop over 1500 km as "expensive" in the response.
4. Return the ordered list with updated `order_index` values persisted to `trip_cities`.

This is intentionally simple for MVP. A future version could use TSP approximation or factor in transport costs.
