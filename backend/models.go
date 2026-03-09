package main

import "time"

// User represents a registered account.
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"display_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// Profile is the public traveler card linked to a user.
type Profile struct {
	ID          int      `json:"id"`
	UserID      int      `json:"user_id"`
	DisplayName string   `json:"display_name"`
	Bio         *string  `json:"bio"`
	TravelStyle string   `json:"travel_style"`
	Interests   []string `json:"interests"`
	IsVisible   bool     `json:"is_visible"`
}

// City is a seeded reference destination.
type City struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Country     string   `json:"country"`
	Region      *string  `json:"region"`
	Description *string  `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	CostLevel   string   `json:"cost_level"`
}

// Attraction is a point of interest within a city.
type Attraction struct {
	ID             int      `json:"id"`
	CityID         int      `json:"city_id"`
	Name           string   `json:"name"`
	Description    *string  `json:"description"`
	Category       *string  `json:"category"`
	EstimatedHours *float64 `json:"estimated_hours"`
	CostLevel      string   `json:"cost_level"`
}

// Trip is a user's multi-city travel plan.
type Trip struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Name        string     `json:"name"`
	StartDate   *string    `json:"start_date"`
	EndDate     *string    `json:"end_date"`
	BudgetStyle string     `json:"budget_style"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	Cities      []TripCity `json:"cities,omitempty"`
}

// TripCity is a city added to a trip with optional dates and ordering.
type TripCity struct {
	ID            int     `json:"id"`
	TripID        int     `json:"trip_id"`
	CityID        int     `json:"city_id"`
	CityName      string  `json:"city_name,omitempty"`
	Country       string  `json:"country,omitempty"`
	ArrivalDate   *string `json:"arrival_date"`
	DepartureDate *string `json:"departure_date"`
	OrderIndex    int     `json:"order_index"`
}

// ItineraryItem links an attraction to a trip day.
type ItineraryItem struct {
	ID             int     `json:"id"`
	TripID         int     `json:"trip_id"`
	AttractionID   int     `json:"attraction_id"`
	AttractionName string  `json:"attraction_name,omitempty"`
	CityName       string  `json:"city_name,omitempty"`
	ScheduledDate  *string `json:"scheduled_date"`
	Notes          *string `json:"notes"`
}

// Connection is a traveler-to-traveler request scoped to a trip.
type Connection struct {
	ID              int     `json:"id"`
	RequesterID     int     `json:"requester_id"`
	RecipientID     int     `json:"recipient_id"`
	TripID          int     `json:"trip_id"`
	Status          string  `json:"status"`
	Message         *string `json:"message"`
	RequesterName   string  `json:"requester_name,omitempty"`
	RecipientName   string  `json:"recipient_name,omitempty"`
	TripName        string  `json:"trip_name,omitempty"`
}

// RouteHop represents one leg between cities in the optimized route.
type RouteHop struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	DistanceKm float64 `json:"distance_km"`
	CostLevel  string  `json:"cost_level"`
}

// RouteResponse is returned by the route optimizer endpoint.
type RouteResponse struct {
	TripID         int        `json:"trip_id"`
	OptimizedOrder []TripCity `json:"optimized_order"`
	Hops           []RouteHop `json:"hops"`
}

// MatchedTraveler represents a user with overlapping trip cities/dates.
type MatchedTraveler struct {
	UserID        int      `json:"user_id"`
	DisplayName   string   `json:"display_name"`
	TravelStyle   string   `json:"travel_style"`
	Interests     []string `json:"interests"`
	OverlapCities []string `json:"overlap_cities"`
}
