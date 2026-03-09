package main

import (
	"log"
	"net/http"
)

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	initJWTSecret()

	db := ConnectDB()
	if db != nil {
		defer db.Close()
	}

	app := &App{DB: db}
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /health", cors(app.HealthHandler))

	// Auth
	mux.HandleFunc("POST /register", cors(app.RegisterHandler))
	mux.HandleFunc("POST /login", cors(app.LoginHandler))
	mux.HandleFunc("GET /auth/verify", cors(requireAuth(app.VerifyHandler)))

	// Profile (authenticated)
	mux.HandleFunc("GET /profile", cors(requireAuth(app.GetProfile)))
	mux.HandleFunc("PUT /profile", cors(requireAuth(app.UpdateProfile)))

	// Cities (public)
	mux.HandleFunc("GET /cities", cors(app.ListCities))
	mux.HandleFunc("GET /cities/{id}", cors(app.GetCity))

	// Attractions (public)
	mux.HandleFunc("GET /cities/{id}/attractions", cors(app.ListCityAttractions))
	mux.HandleFunc("GET /attractions/{id}", cors(app.GetAttraction))

	// Trips (authenticated)
	mux.HandleFunc("GET /trips", cors(requireAuth(app.ListTrips)))
	mux.HandleFunc("POST /trips", cors(requireAuth(app.CreateTrip)))
	mux.HandleFunc("GET /trips/{id}", cors(requireAuth(app.GetTrip)))
	mux.HandleFunc("PUT /trips/{id}", cors(requireAuth(app.UpdateTrip)))
	mux.HandleFunc("DELETE /trips/{id}", cors(requireAuth(app.DeleteTrip)))

	// Trip Cities (authenticated)
	mux.HandleFunc("POST /trips/{id}/cities", cors(requireAuth(app.AddTripCity)))
	mux.HandleFunc("PUT /trips/{id}/cities/{cityId}", cors(requireAuth(app.UpdateTripCity)))
	mux.HandleFunc("DELETE /trips/{id}/cities/{cityId}", cors(requireAuth(app.RemoveTripCity)))

	// Route optimizer (authenticated)
	mux.HandleFunc("GET /trips/{id}/route", cors(requireAuth(app.GetRoute)))

	// Itinerary (authenticated)
	mux.HandleFunc("GET /trips/{id}/itinerary", cors(requireAuth(app.ListItinerary)))
	mux.HandleFunc("POST /trips/{id}/itinerary", cors(requireAuth(app.AddItineraryItem)))
	mux.HandleFunc("DELETE /trips/{id}/itinerary/{itemId}", cors(requireAuth(app.RemoveItineraryItem)))

	// Matching & Connections (authenticated)
	mux.HandleFunc("GET /trips/{id}/matches", cors(requireAuth(app.GetMatches)))
	mux.HandleFunc("GET /connections", cors(requireAuth(app.ListConnections)))
	mux.HandleFunc("POST /connections", cors(requireAuth(app.CreateConnection)))
	mux.HandleFunc("PUT /connections/{id}", cors(requireAuth(app.UpdateConnection)))

	// OPTIONS preflight
	mux.HandleFunc("OPTIONS /", cors(func(w http.ResponseWriter, r *http.Request) {}))

	addr := ":8080"
	log.Printf("TrailMates backend running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
