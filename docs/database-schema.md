# Database Schema: TrailMates

## Overview

MySQL database. All tables use `INT AUTO_INCREMENT` primary keys. Foreign keys cascade on delete unless noted.

---

## Entity-Relationship Summary

```
users ──< profiles          (1-to-1, user_id)
users ──< trips             (1-to-many, user_id)
trips ──< trip_cities       (1-to-many, trip_id)
cities ──< trip_cities      (1-to-many, city_id)
cities ──< attractions      (1-to-many, city_id)
trips ──< itinerary_items   (1-to-many, trip_id)
attractions ──< itinerary_items  (1-to-many, attraction_id)
users ──< connections       (requester_id, recipient_id)
trips ──< connections       (trip_id — requester's trip)
```

---

## Tables

### `users`
Stores login credentials. No personally identifiable info beyond email.

| Column | Type | Notes |
|--------|------|-------|
| `id` | INT PK AUTO_INCREMENT | |
| `email` | VARCHAR(255) UNIQUE NOT NULL | |
| `password_hash` | VARCHAR(255) NOT NULL | bcrypt |
| `created_at` | TIMESTAMP DEFAULT NOW() | |
| `updated_at` | TIMESTAMP ON UPDATE NOW() | |

---

### `profiles`
One-to-one with `users`. Stores the public-facing traveler card shown in matching.

| Column | Type | Notes |
|--------|------|-------|
| `id` | INT PK AUTO_INCREMENT | |
| `user_id` | INT UNIQUE FK → users(id) | CASCADE delete |
| `display_name` | VARCHAR(100) NOT NULL | |
| `bio` | TEXT | optional |
| `travel_style` | ENUM('budget','mid-range','luxury') DEFAULT 'mid-range' | |
| `interests` | JSON | e.g. `["hiking","food","museums"]` |
| `is_visible` | BOOLEAN DEFAULT TRUE | opt-in to traveler matching |
| `created_at` | TIMESTAMP DEFAULT NOW() | |
| `updated_at` | TIMESTAMP ON UPDATE NOW() | |

---

### `cities`
Seeded reference data. Not user-editable.

| Column | Type | Notes |
|--------|------|-------|
| `id` | INT PK AUTO_INCREMENT | |
| `name` | VARCHAR(100) NOT NULL | |
| `country` | VARCHAR(100) NOT NULL | |
| `region` | VARCHAR(100) | e.g. "Western Europe" |
| `description` | TEXT | |
| `latitude` | DECIMAL(10,8) | for distance calculations |
| `longitude` | DECIMAL(11,8) | |
| `cost_level` | ENUM('low','medium','high') DEFAULT 'medium' | |

---

### `attractions`
Seeded top attractions per city. Not user-editable.

| Column | Type | Notes |
|--------|------|-------|
| `id` | INT PK AUTO_INCREMENT | |
| `city_id` | INT FK → cities(id) | CASCADE delete |
| `name` | VARCHAR(255) NOT NULL | |
| `description` | TEXT | |
| `category` | VARCHAR(50) | museum, food, outdoor, nightlife, history, etc. |
| `estimated_hours` | DECIMAL(4,1) | typical visit duration |
| `cost_level` | ENUM('free','low','medium','high') DEFAULT 'free' | |

---

### `trips`
A user's saved multi-city trip plan.

| Column | Type | Notes |
|--------|------|-------|
| `id` | INT PK AUTO_INCREMENT | |
| `user_id` | INT FK → users(id) | CASCADE delete |
| `name` | VARCHAR(255) NOT NULL | e.g. "Summer 2026 Iberia" |
| `start_date` | DATE | nullable — can plan without fixed dates |
| `end_date` | DATE | nullable |
| `budget_style` | ENUM('budget','mid-range','luxury') DEFAULT 'mid-range' | |
| `status` | ENUM('planning','active','completed') DEFAULT 'planning' | |
| `created_at` | TIMESTAMP DEFAULT NOW() | |
| `updated_at` | TIMESTAMP ON UPDATE NOW() | |

---

### `trip_cities`
Junction table — the cities in a trip, with dates and optimized order.

| Column | Type | Notes |
|--------|------|-------|
| `id` | INT PK AUTO_INCREMENT | |
| `trip_id` | INT FK → trips(id) | CASCADE delete |
| `city_id` | INT FK → cities(id) | no cascade (city is reference data) |
| `arrival_date` | DATE | nullable |
| `departure_date` | DATE | nullable |
| `order_index` | INT NOT NULL DEFAULT 0 | set by route optimizer |

---

### `itinerary_items`
Attractions the user has added to a specific trip day.

| Column | Type | Notes |
|--------|------|-------|
| `id` | INT PK AUTO_INCREMENT | |
| `trip_id` | INT FK → trips(id) | CASCADE delete |
| `attraction_id` | INT FK → attractions(id) | no cascade |
| `scheduled_date` | DATE | nullable — can add without a specific day |
| `notes` | TEXT | user's own notes |
| `created_at` | TIMESTAMP DEFAULT NOW() | |

---

### `connections`
Traveler connection requests between users, scoped to the requester's trip.

| Column | Type | Notes |
|--------|------|-------|
| `id` | INT PK AUTO_INCREMENT | |
| `requester_id` | INT FK → users(id) | CASCADE delete |
| `recipient_id` | INT FK → users(id) | CASCADE delete |
| `trip_id` | INT FK → trips(id) | CASCADE delete; requester's trip |
| `status` | ENUM('pending','accepted','declined') DEFAULT 'pending' | |
| `message` | TEXT | optional intro message |
| `created_at` | TIMESTAMP DEFAULT NOW() | |
| `updated_at` | TIMESTAMP ON UPDATE NOW() | |

**Unique constraint:** `(requester_id, recipient_id, trip_id)` — one request per pairing per trip.

---

## SQL (Create Tables)

```sql
CREATE TABLE users (
  id            INT AUTO_INCREMENT PRIMARY KEY,
  email         VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE profiles (
  id            INT AUTO_INCREMENT PRIMARY KEY,
  user_id       INT NOT NULL UNIQUE,
  display_name  VARCHAR(100) NOT NULL,
  bio           TEXT,
  travel_style  ENUM('budget','mid-range','luxury') DEFAULT 'mid-range',
  interests     JSON,
  is_visible    BOOLEAN DEFAULT TRUE,
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE cities (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(100) NOT NULL,
  country     VARCHAR(100) NOT NULL,
  region      VARCHAR(100),
  description TEXT,
  latitude    DECIMAL(10,8),
  longitude   DECIMAL(11,8),
  cost_level  ENUM('low','medium','high') DEFAULT 'medium'
);

CREATE TABLE attractions (
  id               INT AUTO_INCREMENT PRIMARY KEY,
  city_id          INT NOT NULL,
  name             VARCHAR(255) NOT NULL,
  description      TEXT,
  category         VARCHAR(50),
  estimated_hours  DECIMAL(4,1),
  cost_level       ENUM('free','low','medium','high') DEFAULT 'free',
  FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE CASCADE
);

CREATE TABLE trips (
  id           INT AUTO_INCREMENT PRIMARY KEY,
  user_id      INT NOT NULL,
  name         VARCHAR(255) NOT NULL,
  start_date   DATE,
  end_date     DATE,
  budget_style ENUM('budget','mid-range','luxury') DEFAULT 'mid-range',
  status       ENUM('planning','active','completed') DEFAULT 'planning',
  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE trip_cities (
  id               INT AUTO_INCREMENT PRIMARY KEY,
  trip_id          INT NOT NULL,
  city_id          INT NOT NULL,
  arrival_date     DATE,
  departure_date   DATE,
  order_index      INT NOT NULL DEFAULT 0,
  FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE,
  FOREIGN KEY (city_id) REFERENCES cities(id)
);

CREATE TABLE itinerary_items (
  id              INT AUTO_INCREMENT PRIMARY KEY,
  trip_id         INT NOT NULL,
  attraction_id   INT NOT NULL,
  scheduled_date  DATE,
  notes           TEXT,
  created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (trip_id)       REFERENCES trips(id)       ON DELETE CASCADE,
  FOREIGN KEY (attraction_id) REFERENCES attractions(id)
);

CREATE TABLE connections (
  id            INT AUTO_INCREMENT PRIMARY KEY,
  requester_id  INT NOT NULL,
  recipient_id  INT NOT NULL,
  trip_id       INT NOT NULL,
  status        ENUM('pending','accepted','declined') DEFAULT 'pending',
  message       TEXT,
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uq_connection (requester_id, recipient_id, trip_id),
  FOREIGN KEY (requester_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (trip_id)      REFERENCES trips(id) ON DELETE CASCADE
);
```
