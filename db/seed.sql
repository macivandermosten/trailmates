-- TrailMates Seed Data
USE trailmates;

-- =====================
-- CITIES (15 European cities)
-- =====================
INSERT INTO cities (name, country, region, description, latitude, longitude, cost_level) VALUES
('Lisbon', 'Portugal', 'Southern Europe', 'Hilly coastal capital with pastel buildings, trams, and world-class pastéis de nata.', 38.72225700, -9.13933700, 'low'),
('Porto', 'Portugal', 'Southern Europe', 'Port wine cellars, azulejo-tiled streets, and the Douro River.', 41.14961100, -8.61099400, 'low'),
('Barcelona', 'Spain', 'Southern Europe', 'Gaudí architecture, beaches, and late-night tapas culture.', 41.38506400, 2.17340400, 'medium'),
('Madrid', 'Spain', 'Southern Europe', 'Royal palaces, world-class art museums, and bustling plazas.', 40.41676500, -3.70379200, 'medium'),
('Paris', 'France', 'Western Europe', 'The City of Light — iconic landmarks, cafés, and art at every turn.', 48.85661400, 2.35222200, 'high'),
('Amsterdam', 'Netherlands', 'Western Europe', 'Canal rings, cycling culture, and the Van Gogh Museum.', 52.36676000, 4.89454200, 'high'),
('Berlin', 'Germany', 'Central Europe', 'Vibrant nightlife, street art, history, and affordable living.', 52.52000700, 13.40495400, 'low'),
('Prague', 'Czech Republic', 'Central Europe', 'Gothic spires, cheap beer, and a stunning Old Town Square.', 50.07553800, 14.43780200, 'low'),
('Vienna', 'Austria', 'Central Europe', 'Imperial palaces, classical music, and legendary coffeehouse culture.', 48.20817400, 16.37381900, 'medium'),
('Rome', 'Italy', 'Southern Europe', 'Ancient ruins, Renaissance art, and the best pizza and pasta on Earth.', 41.90278300, 12.49636500, 'medium'),
('Florence', 'Italy', 'Southern Europe', 'Birthplace of the Renaissance — the Uffizi, Duomo, and gelato.', 43.76923600, 11.25588900, 'medium'),
('Budapest', 'Hungary', 'Central Europe', 'Thermal baths, ruin bars, and dramatic Danube views.', 47.49791200, 19.04023500, 'low'),
('Dubrovnik', 'Croatia', 'Southern Europe', 'Walled Old Town on the Adriatic — Game of Thrones fame.', 42.65066500, 18.09442200, 'medium'),
('Athens', 'Greece', 'Southern Europe', 'The Acropolis, souvlaki, and island-hopping launchpad.', 37.98381300, 23.72753600, 'low'),
('Marrakech', 'Morocco', 'North Africa', 'Weekend add-on: souks, riads, and the Atlas Mountains on the doorstep.', 31.62970000, -7.98120000, 'low')
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- =====================
-- ATTRACTIONS (3-4 per city)
-- =====================
INSERT INTO attractions (city_id, name, description, category, estimated_hours, cost_level) VALUES
-- Lisbon (1)
(1, 'Belém Tower', 'Iconic 16th-century riverside fortress and UNESCO site.', 'history', 1.5, 'low'),
(1, 'Time Out Market', 'Gourmet food hall in the Cais do Sodré neighborhood.', 'food', 2.0, 'medium'),
(1, 'Alfama District Walk', 'Wander narrow streets of the oldest neighborhood and hear fado music.', 'outdoor', 2.5, 'free'),
-- Porto (2)
(2, 'Livraria Lello', 'Stunning neo-Gothic bookshop that inspired Harry Potter.', 'history', 1.0, 'low'),
(2, 'Port Wine Cellars', 'Tour and taste at the cellars along the Douro in Vila Nova de Gaia.', 'food', 2.0, 'low'),
(2, 'Ribeira District', 'Colorful waterfront UNESCO neighborhood.', 'outdoor', 2.0, 'free'),
-- Barcelona (3)
(3, 'La Sagrada Família', 'Gaudí''s unfinished masterpiece basilica.', 'history', 2.0, 'medium'),
(3, 'Park Güell', 'Mosaic-covered hilltop park with city views.', 'outdoor', 2.0, 'low'),
(3, 'La Boqueria Market', 'Famous covered market on La Rambla with fresh produce and tapas.', 'food', 1.5, 'low'),
-- Madrid (4)
(4, 'Museo del Prado', 'One of the world''s finest art collections — Velázquez, Goya, El Greco.', 'museum', 3.0, 'medium'),
(4, 'Retiro Park', 'Sprawling royal park with a boating lake and crystal palace.', 'outdoor', 2.0, 'free'),
(4, 'Mercado de San Miguel', 'Gourmet tapas market near Plaza Mayor.', 'food', 1.5, 'medium'),
-- Paris (5)
(5, 'Louvre Museum', 'World''s largest art museum — home of the Mona Lisa.', 'museum', 4.0, 'medium'),
(5, 'Montmartre & Sacré-Cœur', 'Hilltop bohemian quarter with panoramic views.', 'outdoor', 2.5, 'free'),
(5, 'Eiffel Tower', 'The iconic iron tower — best at sunset.', 'history', 2.0, 'medium'),
-- Amsterdam (6)
(6, 'Van Gogh Museum', 'The world''s largest collection of Van Gogh paintings.', 'museum', 2.5, 'medium'),
(6, 'Vondelpark', 'Amsterdam''s beloved green lung — perfect for a picnic.', 'outdoor', 2.0, 'free'),
(6, 'Jordaan Neighborhood', 'Charming canals, indie shops, and brown cafés.', 'outdoor', 2.0, 'free'),
-- Berlin (7)
(7, 'Brandenburg Gate', 'Neoclassical symbol of German reunification.', 'history', 1.0, 'free'),
(7, 'East Side Gallery', 'Open-air gallery on the longest remaining stretch of the Berlin Wall.', 'museum', 1.5, 'free'),
(7, 'Markthalle Neun', 'Weekly Street Food Thursday in a historic market hall.', 'food', 2.0, 'low'),
-- Prague (8)
(8, 'Charles Bridge', 'Gothic stone bridge lined with Baroque statues.', 'history', 1.0, 'free'),
(8, 'Prague Castle', 'The world''s largest ancient castle complex.', 'history', 3.0, 'low'),
(8, 'Old Town Square', 'Astronomical Clock, Týn Church, and lively café scene.', 'outdoor', 1.5, 'free'),
-- Vienna (9)
(9, 'Schönbrunn Palace', 'Habsburg summer residence with magnificent gardens.', 'history', 3.0, 'medium'),
(9, 'Naschmarkt', 'Vienna''s most popular open-air market since the 16th century.', 'food', 1.5, 'low'),
(9, 'MuseumsQuartier', 'One of the world''s largest cultural complexes.', 'museum', 3.0, 'medium'),
-- Rome (10)
(10, 'Colosseum', 'Ancient amphitheater — the icon of Imperial Rome.', 'history', 2.5, 'medium'),
(10, 'Vatican Museums', 'Sistine Chapel ceiling and miles of Renaissance masterpieces.', 'museum', 4.0, 'medium'),
(10, 'Trastevere', 'Cobblestone neighborhood with the best trattorias in Rome.', 'food', 2.0, 'low'),
-- Florence (11)
(11, 'Uffizi Gallery', 'Botticelli''s Birth of Venus and top Renaissance art.', 'museum', 3.0, 'medium'),
(11, 'Ponte Vecchio', 'Medieval stone bridge lined with jewelry shops.', 'history', 1.0, 'free'),
(11, 'Mercato Centrale', 'Two-story food market with traditional Florentine fare.', 'food', 1.5, 'low'),
-- Budapest (12)
(12, 'Széchenyi Thermal Bath', 'Largest medicinal bath in Europe — soak outdoors year-round.', 'outdoor', 3.0, 'low'),
(12, 'Ruin Bars (Szimpla Kert)', 'Iconic ruin pub in a converted apartment building.', 'nightlife', 3.0, 'low'),
(12, 'Fisherman''s Bastion', 'Neo-Gothic terrace with sweeping views of Parliament.', 'history', 1.5, 'free'),
-- Dubrovnik (13)
(13, 'City Walls Walk', '2 km walk atop the medieval walls surrounding Old Town.', 'outdoor', 2.0, 'medium'),
(13, 'Lokrum Island', 'Short ferry to a car-free island with a monastery and swimming.', 'outdoor', 3.0, 'low'),
(13, 'Stradun', 'Limestone-paved main street through Old Town.', 'history', 1.0, 'free'),
-- Athens (14)
(14, 'Acropolis', 'The Parthenon and ancient citadel overlooking Athens.', 'history', 3.0, 'medium'),
(14, 'Plaka District', 'Oldest neighborhood — winding streets, tavernas, and souvenir shops.', 'outdoor', 2.0, 'free'),
(14, 'Central Market', 'Bustling Athens market with meat, fish, and spices.', 'food', 1.0, 'free'),
-- Marrakech (15)
(15, 'Jemaa el-Fnaa', 'Legendary night market square with street food and performers.', 'food', 3.0, 'low'),
(15, 'Majorelle Garden', 'Stunning blue villa and botanical garden once owned by Yves Saint Laurent.', 'outdoor', 1.5, 'low'),
(15, 'Medina Souks', 'Labyrinth of market stalls selling leather, spices, and ceramics.', 'outdoor', 2.5, 'free')
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- =====================
-- SAMPLE USERS (password is "password123" for all — bcrypt hash)
-- =====================
INSERT INTO users (id, email, password_hash) VALUES
(1, 'alex@example.com',  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'),
(2, 'jordan@example.com','$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'),
(3, 'sam@example.com',   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy')
ON DUPLICATE KEY UPDATE email = VALUES(email);

INSERT INTO profiles (user_id, display_name, bio, travel_style, interests, is_visible) VALUES
(1, 'Alex', 'Gap-year backpacker from Canada. Love hiking and street food.', 'budget', '["hiking","food","history"]', TRUE),
(2, 'Jordan', 'Digital nomad hopping around Europe. Into museums and coffee.', 'mid-range', '["museums","food","nightlife"]', TRUE),
(3, 'Sam', 'Solo traveler from Australia. Outdoor adventures and local culture.', 'budget', '["hiking","outdoor","history"]', TRUE)
ON DUPLICATE KEY UPDATE display_name = VALUES(display_name);

-- =====================
-- SAMPLE TRIPS
-- =====================
INSERT INTO trips (id, user_id, name, start_date, end_date, budget_style, status) VALUES
(1, 1, 'Summer 2026 Iberia', '2026-06-10', '2026-06-25', 'budget', 'planning'),
(2, 2, 'Central Europe Art Tour', '2026-06-15', '2026-07-01', 'mid-range', 'planning'),
(3, 3, 'Mediterranean Loop', '2026-06-12', '2026-06-28', 'budget', 'planning')
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- Trip 1: Alex — Lisbon → Madrid → Barcelona
INSERT INTO trip_cities (trip_id, city_id, arrival_date, departure_date, order_index) VALUES
(1, 1, '2026-06-10', '2026-06-14', 1),
(1, 4, '2026-06-14', '2026-06-18', 2),
(1, 3, '2026-06-18', '2026-06-25', 3)
ON DUPLICATE KEY UPDATE order_index = VALUES(order_index);

-- Trip 2: Jordan — Vienna → Prague → Berlin → Amsterdam
INSERT INTO trip_cities (trip_id, city_id, arrival_date, departure_date, order_index) VALUES
(2, 9, '2026-06-15', '2026-06-19', 1),
(2, 8, '2026-06-19', '2026-06-23', 2),
(2, 7, '2026-06-23', '2026-06-27', 3),
(2, 6, '2026-06-27', '2026-07-01', 4)
ON DUPLICATE KEY UPDATE order_index = VALUES(order_index);

-- Trip 3: Sam — Rome → Florence → Barcelona → Lisbon
INSERT INTO trip_cities (trip_id, city_id, arrival_date, departure_date, order_index) VALUES
(3, 10, '2026-06-12', '2026-06-16', 1),
(3, 11, '2026-06-16', '2026-06-19', 2),
(3, 3,  '2026-06-19', '2026-06-23', 3),
(3, 1,  '2026-06-23', '2026-06-28', 4)
ON DUPLICATE KEY UPDATE order_index = VALUES(order_index);

-- =====================
-- SAMPLE ITINERARY ITEMS
-- =====================
INSERT INTO itinerary_items (trip_id, attraction_id, scheduled_date, notes) VALUES
(1, 1, '2026-06-11', 'Morning visit — arrive early to beat the crowds'),
(1, 3, '2026-06-12', 'Afternoon walk through Alfama, catch fado at night'),
(1, 10, '2026-06-15', 'Spend full morning at the Prado'),
(2, 27, '2026-06-16', 'Book Schönbrunn tickets in advance'),
(2, 24, '2026-06-20', 'Walk Charles Bridge at sunrise for photos'),
(3, 31, '2026-06-13', 'Colosseum first thing, then Forum');

-- =====================
-- SAMPLE CONNECTION
-- =====================
INSERT INTO connections (requester_id, recipient_id, trip_id, status, message) VALUES
(1, 3, 1, 'pending', 'Hey Sam! I see we overlap in Barcelona and Lisbon. Want to explore together?')
ON DUPLICATE KEY UPDATE status = VALUES(status);
