# TrailMates

Plan smarter Europe trips and connect with compatible solo travelers.

## What It Does

TrailMates helps backpackers build multi-city Europe itineraries, reorder stops to cut backtracking, browse top attractions city by city, and find other travelers with overlapping plans.

## Live App

**[macicapstone.duckdns.org](https://macicapstone.duckdns.org)**

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22, MySQL |
| Frontend | React 19, Vite 7, Tailwind v4, Shadcn UI |
| Auth | JWT (per-user accounts) |
| Deployment | AWS Lightsail, GitHub Actions, Nginx |

## Pages

| Route | Description |
|-------|-------------|
| `/` | Landing page — value prop + sign up CTA |
| `/cities` | Browse supported European cities |
| `/cities/:id` | City detail — top attractions |
| `/signup` `/login` | Auth |
| `/app` | Dashboard — your saved trips |
| `/app/trips/new` | Create a new trip |
| `/app/trips/:id` | Edit trip — cities + dates |
| `/app/trips/:id/route` | Optimized route view |
| `/app/trips/:id/itinerary` | Day-by-day itinerary |
| `/app/trips/:id/matches` | Travelers with overlapping trips |
| `/app/profile` | Edit profile + travel preferences |

## Local Development

### Backend
```bash
cd backend
export DB_USER=root DB_PASSWORD=... DB_HOST=localhost DB_NAME=trailmates JWT_SECRET=...
go run .
```

### Frontend
```bash
cd frontend
npm install
npm run dev   # proxies /api to localhost:8080
```

## Deployment

Push to `main` → GitHub Actions builds frontend + backend and deploys to Lightsail via SSH.

## Documentation

- [Project Proposal](docs/project-proposal.md)
- [Architecture](docs/architecture.md)
- [Database Schema](docs/database-schema.md)
