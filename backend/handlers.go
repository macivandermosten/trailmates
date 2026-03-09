package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// App holds the database connection shared across handlers.
type App struct {
	DB *sql.DB
}

// JSON helpers
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (a *App) requireDB(w http.ResponseWriter) bool {
	if a.DB == nil {
		writeError(w, http.StatusServiceUnavailable, "database not available")
		return false
	}
	return true
}

// --- Health ---

func (a *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]string{"app": "trailmates", "status": "ok"}
	if a.DB != nil {
		if err := a.DB.Ping(); err != nil {
			status["database"] = "error"
		} else {
			status["database"] = "connected"
		}
	} else {
		status["database"] = "not configured"
	}
	writeJSON(w, http.StatusOK, status)
}

// =====================
// CITIES (public, seeded)
// =====================

func (a *App) ListCities(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	rows, err := a.DB.Query("SELECT id, name, country, region, description, latitude, longitude, cost_level FROM cities ORDER BY country, name")
	if err != nil {
		log.Printf("ListCities error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to query cities")
		return
	}
	defer rows.Close()

	cities := []City{}
	for rows.Next() {
		var c City
		if err := rows.Scan(&c.ID, &c.Name, &c.Country, &c.Region, &c.Description, &c.Latitude, &c.Longitude, &c.CostLevel); err != nil {
			log.Printf("ListCities scan error: %v", err)
			continue
		}
		cities = append(cities, c)
	}
	writeJSON(w, http.StatusOK, cities)
}

func (a *App) GetCity(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid city id")
		return
	}

	var c City
	err = a.DB.QueryRow("SELECT id, name, country, region, description, latitude, longitude, cost_level FROM cities WHERE id = ?", id).
		Scan(&c.ID, &c.Name, &c.Country, &c.Region, &c.Description, &c.Latitude, &c.Longitude, &c.CostLevel)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "city not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query city")
		return
	}

	// Include attractions
	attrRows, err := a.DB.Query("SELECT id, city_id, name, description, category, estimated_hours, cost_level FROM attractions WHERE city_id = ? ORDER BY name", id)
	if err == nil {
		defer attrRows.Close()
		attractions := []Attraction{}
		for attrRows.Next() {
			var attr Attraction
			if err := attrRows.Scan(&attr.ID, &attr.CityID, &attr.Name, &attr.Description, &attr.Category, &attr.EstimatedHours, &attr.CostLevel); err == nil {
				attractions = append(attractions, attr)
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"city": c, "attractions": attractions})
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// =====================
// ATTRACTIONS (public)
// =====================

func (a *App) ListCityAttractions(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	cityID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid city id")
		return
	}
	rows, err := a.DB.Query("SELECT id, city_id, name, description, category, estimated_hours, cost_level FROM attractions WHERE city_id = ? ORDER BY name", cityID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query attractions")
		return
	}
	defer rows.Close()

	attractions := []Attraction{}
	for rows.Next() {
		var attr Attraction
		if err := rows.Scan(&attr.ID, &attr.CityID, &attr.Name, &attr.Description, &attr.Category, &attr.EstimatedHours, &attr.CostLevel); err == nil {
			attractions = append(attractions, attr)
		}
	}
	writeJSON(w, http.StatusOK, attractions)
}

func (a *App) GetAttraction(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid attraction id")
		return
	}
	var attr Attraction
	err = a.DB.QueryRow("SELECT id, city_id, name, description, category, estimated_hours, cost_level FROM attractions WHERE id = ?", id).
		Scan(&attr.ID, &attr.CityID, &attr.Name, &attr.Description, &attr.Category, &attr.EstimatedHours, &attr.CostLevel)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "attraction not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query attraction")
		return
	}
	writeJSON(w, http.StatusOK, attr)
}

// =====================
// PROFILE
// =====================

func (a *App) GetProfile(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)

	var p Profile
	var interestsJSON sql.NullString
	err := a.DB.QueryRow("SELECT id, user_id, display_name, bio, travel_style, interests, is_visible FROM profiles WHERE user_id = ?", uid).
		Scan(&p.ID, &p.UserID, &p.DisplayName, &p.Bio, &p.TravelStyle, &interestsJSON, &p.IsVisible)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "profile not found")
		return
	} else if err != nil {
		log.Printf("GetProfile error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to query profile")
		return
	}
	if interestsJSON.Valid {
		json.Unmarshal([]byte(interestsJSON.String), &p.Interests)
	}
	if p.Interests == nil {
		p.Interests = []string{}
	}
	writeJSON(w, http.StatusOK, p)
}

func (a *App) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)

	var p Profile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if p.DisplayName == "" {
		writeError(w, http.StatusBadRequest, "display_name is required")
		return
	}

	interestsBytes, _ := json.Marshal(p.Interests)

	_, err := a.DB.Exec(
		"UPDATE profiles SET display_name = ?, bio = ?, travel_style = ?, interests = ?, is_visible = ? WHERE user_id = ?",
		p.DisplayName, p.Bio, p.TravelStyle, string(interestsBytes), p.IsVisible, uid,
	)
	if err != nil {
		log.Printf("UpdateProfile error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}
	a.GetProfile(w, r)
}

// =====================
// TRIPS
// =====================

func (a *App) ListTrips(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)

	rows, err := a.DB.Query(
		"SELECT id, user_id, name, start_date, end_date, budget_style, status, created_at FROM trips WHERE user_id = ? ORDER BY created_at DESC", uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query trips")
		return
	}
	defer rows.Close()

	trips := []Trip{}
	for rows.Next() {
		var t Trip
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.StartDate, &t.EndDate, &t.BudgetStyle, &t.Status, &t.CreatedAt); err == nil {
			trips = append(trips, t)
		}
	}
	writeJSON(w, http.StatusOK, trips)
}

func (a *App) GetTrip(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}

	var t Trip
	err = a.DB.QueryRow(
		"SELECT id, user_id, name, start_date, end_date, budget_style, status, created_at FROM trips WHERE id = ? AND user_id = ?", id, uid,
	).Scan(&t.ID, &t.UserID, &t.Name, &t.StartDate, &t.EndDate, &t.BudgetStyle, &t.Status, &t.CreatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query trip")
		return
	}

	// Load trip cities
	cityRows, err := a.DB.Query(
		`SELECT tc.id, tc.trip_id, tc.city_id, c.name, c.country, tc.arrival_date, tc.departure_date, tc.order_index
		 FROM trip_cities tc JOIN cities c ON tc.city_id = c.id
		 WHERE tc.trip_id = ? ORDER BY tc.order_index`, id)
	if err == nil {
		defer cityRows.Close()
		for cityRows.Next() {
			var tc TripCity
			if err := cityRows.Scan(&tc.ID, &tc.TripID, &tc.CityID, &tc.CityName, &tc.Country, &tc.ArrivalDate, &tc.DepartureDate, &tc.OrderIndex); err == nil {
				t.Cities = append(t.Cities, tc)
			}
		}
	}
	if t.Cities == nil {
		t.Cities = []TripCity{}
	}

	writeJSON(w, http.StatusOK, t)
}

func (a *App) CreateTrip(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)

	var t Trip
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if t.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if t.BudgetStyle == "" {
		t.BudgetStyle = "mid-range"
	}

	res, err := a.DB.Exec(
		"INSERT INTO trips (user_id, name, start_date, end_date, budget_style) VALUES (?, ?, ?, ?, ?)",
		uid, t.Name, t.StartDate, t.EndDate, t.BudgetStyle,
	)
	if err != nil {
		log.Printf("CreateTrip error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create trip")
		return
	}
	id, _ := res.LastInsertId()
	t.ID = int(id)
	t.UserID = uid
	t.Status = "planning"
	t.Cities = []TripCity{}
	writeJSON(w, http.StatusCreated, t)
}

func (a *App) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}

	var t Trip
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if t.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if t.BudgetStyle == "" {
		t.BudgetStyle = "mid-range"
	}
	if t.Status == "" {
		t.Status = "planning"
	}

	res, err := a.DB.Exec(
		"UPDATE trips SET name = ?, start_date = ?, end_date = ?, budget_style = ?, status = ? WHERE id = ? AND user_id = ?",
		t.Name, t.StartDate, t.EndDate, t.BudgetStyle, t.Status, id, uid,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update trip")
		return
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}
	t.ID = id
	t.UserID = uid
	writeJSON(w, http.StatusOK, t)
}

func (a *App) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}

	res, err := a.DB.Exec("DELETE FROM trips WHERE id = ? AND user_id = ?", id, uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete trip")
		return
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// =====================
// TRIP CITIES
// =====================

func (a *App) AddTripCity(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	tripID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}

	// Verify trip ownership
	var ownerID int
	if err := a.DB.QueryRow("SELECT user_id FROM trips WHERE id = ?", tripID).Scan(&ownerID); err != nil || ownerID != uid {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}

	var tc TripCity
	if err := json.NewDecoder(r.Body).Decode(&tc); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if tc.CityID == 0 {
		writeError(w, http.StatusBadRequest, "city_id is required")
		return
	}

	// Get max order_index
	var maxOrder int
	a.DB.QueryRow("SELECT COALESCE(MAX(order_index), 0) FROM trip_cities WHERE trip_id = ?", tripID).Scan(&maxOrder)
	tc.OrderIndex = maxOrder + 1

	res, err := a.DB.Exec(
		"INSERT INTO trip_cities (trip_id, city_id, arrival_date, departure_date, order_index) VALUES (?, ?, ?, ?, ?)",
		tripID, tc.CityID, tc.ArrivalDate, tc.DepartureDate, tc.OrderIndex,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add city to trip")
		return
	}
	id, _ := res.LastInsertId()
	tc.ID = int(id)
	tc.TripID = tripID
	writeJSON(w, http.StatusCreated, tc)
}

func (a *App) UpdateTripCity(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	tripID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}
	cityID, err := strconv.Atoi(r.PathValue("cityId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid city id")
		return
	}

	var ownerID int
	if err := a.DB.QueryRow("SELECT user_id FROM trips WHERE id = ?", tripID).Scan(&ownerID); err != nil || ownerID != uid {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}

	var tc TripCity
	if err := json.NewDecoder(r.Body).Decode(&tc); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	_, err = a.DB.Exec(
		"UPDATE trip_cities SET arrival_date = ?, departure_date = ?, order_index = ? WHERE trip_id = ? AND city_id = ?",
		tc.ArrivalDate, tc.DepartureDate, tc.OrderIndex, tripID, cityID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update trip city")
		return
	}
	writeJSON(w, http.StatusOK, tc)
}

func (a *App) RemoveTripCity(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	tripID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}
	cityID, err := strconv.Atoi(r.PathValue("cityId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid city id")
		return
	}

	var ownerID int
	if err := a.DB.QueryRow("SELECT user_id FROM trips WHERE id = ?", tripID).Scan(&ownerID); err != nil || ownerID != uid {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}

	a.DB.Exec("DELETE FROM trip_cities WHERE trip_id = ? AND city_id = ?", tripID, cityID)
	w.WriteHeader(http.StatusNoContent)
}

// =====================
// ROUTE OPTIMIZER
// =====================

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func (a *App) GetRoute(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	tripID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}

	var ownerID int
	if err := a.DB.QueryRow("SELECT user_id FROM trips WHERE id = ?", tripID).Scan(&ownerID); err != nil || ownerID != uid {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}

	rows, err := a.DB.Query(
		`SELECT tc.id, tc.trip_id, tc.city_id, c.name, c.country, tc.arrival_date, tc.departure_date, tc.order_index, c.latitude, c.longitude, c.cost_level
		 FROM trip_cities tc JOIN cities c ON tc.city_id = c.id
		 WHERE tc.trip_id = ?`, tripID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query trip cities")
		return
	}
	defer rows.Close()

	type cityWithCoords struct {
		TripCity
		Lat       float64
		Lon       float64
		CostLevel string
	}
	var cities []cityWithCoords
	for rows.Next() {
		var c cityWithCoords
		var lat, lon sql.NullFloat64
		if err := rows.Scan(&c.ID, &c.TripID, &c.CityID, &c.CityName, &c.Country, &c.ArrivalDate, &c.DepartureDate, &c.OrderIndex, &lat, &lon, &c.CostLevel); err == nil {
			if lat.Valid {
				c.Lat = lat.Float64
			}
			if lon.Valid {
				c.Lon = lon.Float64
			}
			cities = append(cities, c)
		}
	}

	if len(cities) < 2 {
		resp := RouteResponse{TripID: tripID, OptimizedOrder: []TripCity{}, Hops: []RouteHop{}}
		for _, c := range cities {
			resp.OptimizedOrder = append(resp.OptimizedOrder, c.TripCity)
		}
		writeJSON(w, http.StatusOK, resp)
		return
	}

	// Nearest-neighbor greedy: start with leftmost (westernmost) city
	sort.Slice(cities, func(i, j int) bool { return cities[i].Lon < cities[j].Lon })

	visited := make([]bool, len(cities))
	order := []cityWithCoords{cities[0]}
	visited[0] = true

	for len(order) < len(cities) {
		last := order[len(order)-1]
		bestIdx := -1
		bestDist := math.MaxFloat64
		for i, c := range cities {
			if visited[i] {
				continue
			}
			d := haversine(last.Lat, last.Lon, c.Lat, c.Lon)
			if d < bestDist {
				bestDist = d
				bestIdx = i
			}
		}
		visited[bestIdx] = true
		order = append(order, cities[bestIdx])
	}

	// Build response and update order_index in DB
	resp := RouteResponse{TripID: tripID}
	for i, c := range order {
		c.TripCity.OrderIndex = i + 1
		resp.OptimizedOrder = append(resp.OptimizedOrder, c.TripCity)
		a.DB.Exec("UPDATE trip_cities SET order_index = ? WHERE id = ?", i+1, c.ID)

		if i > 0 {
			prev := order[i-1]
			dist := haversine(prev.Lat, prev.Lon, c.Lat, c.Lon)
			costLevel := "low"
			if dist > 1500 {
				costLevel = "high"
			} else if dist > 800 {
				costLevel = "medium"
			}
			resp.Hops = append(resp.Hops, RouteHop{
				From:       prev.CityName,
				To:         c.CityName,
				DistanceKm: math.Round(dist),
				CostLevel:  costLevel,
			})
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

// =====================
// ITINERARY
// =====================

func (a *App) ListItinerary(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	tripID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}

	var ownerID int
	if err := a.DB.QueryRow("SELECT user_id FROM trips WHERE id = ?", tripID).Scan(&ownerID); err != nil || ownerID != uid {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}

	rows, err := a.DB.Query(
		`SELECT ii.id, ii.trip_id, ii.attraction_id, a.name, c.name, ii.scheduled_date, ii.notes
		 FROM itinerary_items ii
		 JOIN attractions a ON ii.attraction_id = a.id
		 JOIN cities c ON a.city_id = c.id
		 WHERE ii.trip_id = ? ORDER BY ii.scheduled_date, a.name`, tripID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query itinerary")
		return
	}
	defer rows.Close()

	items := []ItineraryItem{}
	for rows.Next() {
		var item ItineraryItem
		if err := rows.Scan(&item.ID, &item.TripID, &item.AttractionID, &item.AttractionName, &item.CityName, &item.ScheduledDate, &item.Notes); err == nil {
			items = append(items, item)
		}
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *App) AddItineraryItem(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	tripID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}

	var ownerID int
	if err := a.DB.QueryRow("SELECT user_id FROM trips WHERE id = ?", tripID).Scan(&ownerID); err != nil || ownerID != uid {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}

	var item ItineraryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if item.AttractionID == 0 {
		writeError(w, http.StatusBadRequest, "attraction_id is required")
		return
	}

	res, err := a.DB.Exec(
		"INSERT INTO itinerary_items (trip_id, attraction_id, scheduled_date, notes) VALUES (?, ?, ?, ?)",
		tripID, item.AttractionID, item.ScheduledDate, item.Notes,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add itinerary item")
		return
	}
	id, _ := res.LastInsertId()
	item.ID = int(id)
	item.TripID = tripID
	writeJSON(w, http.StatusCreated, item)
}

func (a *App) RemoveItineraryItem(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	tripID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}
	itemID, err := strconv.Atoi(r.PathValue("itemId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var ownerID int
	if err := a.DB.QueryRow("SELECT user_id FROM trips WHERE id = ?", tripID).Scan(&ownerID); err != nil || ownerID != uid {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}

	a.DB.Exec("DELETE FROM itinerary_items WHERE id = ? AND trip_id = ?", itemID, tripID)
	w.WriteHeader(http.StatusNoContent)
}

// =====================
// MATCHING & CONNECTIONS
// =====================

func (a *App) GetMatches(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	tripID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}

	var ownerID int
	if err := a.DB.QueryRow("SELECT user_id FROM trips WHERE id = ?", tripID).Scan(&ownerID); err != nil || ownerID != uid {
		writeError(w, http.StatusNotFound, "trip not found")
		return
	}

	// Find users with overlapping cities in their trips who are visible
	rows, err := a.DB.Query(`
		SELECT DISTINCT p.user_id, p.display_name, p.travel_style, p.interests, c.name
		FROM trip_cities tc1
		JOIN trip_cities tc2 ON tc1.city_id = tc2.city_id
		JOIN trips t2 ON tc2.trip_id = t2.id
		JOIN profiles p ON t2.user_id = p.user_id
		JOIN cities c ON tc1.city_id = c.id
		WHERE tc1.trip_id = ? AND t2.user_id != ? AND p.is_visible = TRUE
		ORDER BY p.user_id, c.name`, tripID, uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query matches")
		return
	}
	defer rows.Close()

	matchMap := map[int]*MatchedTraveler{}
	for rows.Next() {
		var userID int
		var displayName, travelStyle, cityName string
		var interestsJSON sql.NullString
		if err := rows.Scan(&userID, &displayName, &travelStyle, &interestsJSON, &cityName); err == nil {
			if _, ok := matchMap[userID]; !ok {
				interests := []string{}
				if interestsJSON.Valid {
					json.Unmarshal([]byte(interestsJSON.String), &interests)
				}
				matchMap[userID] = &MatchedTraveler{
					UserID:      userID,
					DisplayName: displayName,
					TravelStyle: travelStyle,
					Interests:   interests,
				}
			}
			matchMap[userID].OverlapCities = append(matchMap[userID].OverlapCities, cityName)
		}
	}

	matches := []MatchedTraveler{}
	for _, m := range matchMap {
		matches = append(matches, *m)
	}
	writeJSON(w, http.StatusOK, matches)
}

func (a *App) ListConnections(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)

	rows, err := a.DB.Query(`
		SELECT c.id, c.requester_id, c.recipient_id, c.trip_id, c.status, c.message,
		       p1.display_name, p2.display_name, t.name
		FROM connections c
		JOIN profiles p1 ON c.requester_id = p1.user_id
		JOIN profiles p2 ON c.recipient_id = p2.user_id
		JOIN trips t ON c.trip_id = t.id
		WHERE c.requester_id = ? OR c.recipient_id = ?
		ORDER BY c.created_at DESC`, uid, uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query connections")
		return
	}
	defer rows.Close()

	connections := []Connection{}
	for rows.Next() {
		var conn Connection
		if err := rows.Scan(&conn.ID, &conn.RequesterID, &conn.RecipientID, &conn.TripID, &conn.Status, &conn.Message,
			&conn.RequesterName, &conn.RecipientName, &conn.TripName); err == nil {
			connections = append(connections, conn)
		}
	}
	writeJSON(w, http.StatusOK, connections)
}

func (a *App) CreateConnection(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)

	var conn Connection
	if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if conn.RecipientID == 0 || conn.TripID == 0 {
		writeError(w, http.StatusBadRequest, "recipient_id and trip_id are required")
		return
	}
	if conn.RecipientID == uid {
		writeError(w, http.StatusBadRequest, "cannot connect with yourself")
		return
	}

	res, err := a.DB.Exec(
		"INSERT INTO connections (requester_id, recipient_id, trip_id, message) VALUES (?, ?, ?, ?)",
		uid, conn.RecipientID, conn.TripID, conn.Message,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			writeError(w, http.StatusConflict, "connection request already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create connection")
		return
	}
	id, _ := res.LastInsertId()
	conn.ID = int(id)
	conn.RequesterID = uid
	conn.Status = "pending"
	writeJSON(w, http.StatusCreated, conn)
}

func (a *App) UpdateConnection(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	uid := getUserID(r)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid connection id")
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if body.Status != "accepted" && body.Status != "declined" {
		writeError(w, http.StatusBadRequest, "status must be 'accepted' or 'declined'")
		return
	}

	// Only the recipient can accept/decline
	res, err := a.DB.Exec(
		"UPDATE connections SET status = ? WHERE id = ? AND recipient_id = ?",
		body.Status, id, uid,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update connection")
		return
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		writeError(w, http.StatusNotFound, "connection not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": body.Status})
}
