# Project Proposal: TrailMates

## Project Name
TrailMates

## Target Audience
Solo travelers and backpackers aged 18–35 planning multi-city Europe trips who want a smarter route and a lightweight way to connect with other travelers nearby.

## Problem Statement
Planning a Europe backpacking route is overwhelming. Travelers bounce between cities based on vibes or TikTok, then realize they are backtracking, overpaying for transport, or missing easy add-ons (like a Spain → Morocco weekend). Solo travel can also feel isolating without a safe, low-friction way to meet people on similar routes.

## Value Proposition
TrailMates helps solo travelers plan the most efficient multi-city Europe trip and connect with compatible travelers by generating optimized routes and matching people by dates, cities, and interests.

---

## MVP Must-Have Features

1. **User accounts + profiles** — Sign up, log in, and set travel style, interests, and a bio. Control visibility for matching.
2. **Trip builder** — Select cities, set a date range and budget style, and save a trip.
3. **Route optimizer** — Reorder selected cities to minimize backtracking and flag expensive or illogical hops between cities.
4. **Attractions browser** — Browse top attractions by city and add them to a trip's itinerary.
5. **Traveler matching** — See other users whose trips overlap by city and date window; send a connection request.

---

## Explicitly Not Building (MVP Scope)

- **In-app payments, bookings, or affiliate links** — TrailMates does not facilitate or monetize travel purchases.
- **Real-time chat** — MVP uses connection requests only; no messaging system.
- **Global coverage** — MVP covers Europe plus a small set of "weekend add-on" destinations (e.g., Morocco). No Asia, Americas, etc.
- **Live airfare or hotel price scraping** — Cost levels are stored estimates ("budget / mid-range / luxury"), not live data.
- **Mobile native app** — Web only for MVP.

---

## Pages and User Flow

### Public (no login required)
| Page | Path | Description |
|------|------|-------------|
| Landing | `/` | Value proposition, how it works, sign-up CTA |
| Explore Cities | `/cities` | Grid of supported cities with country and cost level |
| City Detail | `/cities/:id` | City overview, top attractions, "add to my trip" |

### Auth
| Page | Path | Description |
|------|------|-------------|
| Sign Up | `/signup` | Email + password + display name |
| Login | `/login` | Email + password |

### App (login required)
| Page | Path | Description |
|------|------|-------------|
| Dashboard | `/app` | List of saved trips + "create trip" button |
| Create / Edit Trip | `/app/trips/new`, `/app/trips/:id` | Select cities, date range, budget style |
| Route View | `/app/trips/:id/route` | Optimized city order + hop summary (distance, cost level) |
| Itinerary | `/app/trips/:id/itinerary` | Day-by-day view of saved attractions |
| Matches | `/app/trips/:id/matches` | Travelers overlapping this trip; send connection request |
| Profile | `/app/profile` | Edit bio, travel style, interests, matching visibility |

### Navigation Flow
Landing → Sign Up / Login → Dashboard → Create Trip → Route View → Itinerary → Matches

---

## Success Criteria for Module Completion
- Repository contains `CLAUDE.md`, `README.md`, and `docs/` folder with proposal, architecture, and schema
- GitHub Issues created for all MVP features
- Hello World page accessible over HTTPS at `macicapstone.duckdns.org`
- GitHub Actions deployment pipeline succeeds
