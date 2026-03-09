-- TrailMates Database Schema

CREATE DATABASE IF NOT EXISTS trailmates;
USE trailmates;

CREATE TABLE IF NOT EXISTS users (
  id            INT AUTO_INCREMENT PRIMARY KEY,
  email         VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS profiles (
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

CREATE TABLE IF NOT EXISTS cities (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(100) NOT NULL,
  country     VARCHAR(100) NOT NULL,
  region      VARCHAR(100),
  description TEXT,
  latitude    DECIMAL(10,8),
  longitude   DECIMAL(11,8),
  cost_level  ENUM('low','medium','high') DEFAULT 'medium'
);

CREATE TABLE IF NOT EXISTS attractions (
  id               INT AUTO_INCREMENT PRIMARY KEY,
  city_id          INT NOT NULL,
  name             VARCHAR(255) NOT NULL,
  description      TEXT,
  category         VARCHAR(50),
  estimated_hours  DECIMAL(4,1),
  cost_level       ENUM('free','low','medium','high') DEFAULT 'free',
  FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS trips (
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

CREATE TABLE IF NOT EXISTS trip_cities (
  id               INT AUTO_INCREMENT PRIMARY KEY,
  trip_id          INT NOT NULL,
  city_id          INT NOT NULL,
  arrival_date     DATE,
  departure_date   DATE,
  order_index      INT NOT NULL DEFAULT 0,
  FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE,
  FOREIGN KEY (city_id) REFERENCES cities(id)
);

CREATE TABLE IF NOT EXISTS itinerary_items (
  id              INT AUTO_INCREMENT PRIMARY KEY,
  trip_id         INT NOT NULL,
  attraction_id   INT NOT NULL,
  scheduled_date  DATE,
  notes           TEXT,
  created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (trip_id)       REFERENCES trips(id)       ON DELETE CASCADE,
  FOREIGN KEY (attraction_id) REFERENCES attractions(id)
);

CREATE TABLE IF NOT EXISTS connections (
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
