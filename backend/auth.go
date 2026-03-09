package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func initJWTSecret() {
	s := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if s == "" {
		s = "dev-secret-change-me"
		log.Println("WARNING: JWT_SECRET not set, using insecure default")
	}
	jwtSecret = []byte(s)
}

type contextKey string

const userIDKey contextKey = "user_id"

// generateToken creates a signed JWT for a user.
func generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

// requireAuth middleware validates the JWT and injects user_id into context.
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid claims")
			return
		}
		uid, ok := claims["user_id"].(float64)
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid user_id in token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, int(uid))
		next(w, r.WithContext(ctx))
	}
}

// getUserID extracts user_id from the request context.
func getUserID(r *http.Request) int {
	return r.Context().Value(userIDKey).(int)
}

// --- Register ---

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func (a *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		writeError(w, http.StatusBadRequest, "email, password, and display_name are required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	tx, err := a.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", req.Email, string(hash))
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		log.Printf("Register insert user error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	userID64, _ := res.LastInsertId()
	userID := int(userID64)

	_, err = tx.Exec("INSERT INTO profiles (user_id, display_name) VALUES (?, ?)", userID, req.DisplayName)
	if err != nil {
		log.Printf("Register insert profile error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create profile")
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	token, err := generateToken(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"token": token,
		"user":  map[string]any{"id": userID, "email": req.Email, "display_name": req.DisplayName},
	})
}

// --- Login ---

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !a.requireDB(w) {
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	var user User
	var displayName string
	err := a.DB.QueryRow(
		"SELECT u.id, u.email, u.password_hash, p.display_name FROM users u LEFT JOIN profiles p ON p.user_id = u.id WHERE u.email = ?",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &displayName)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	} else if err != nil {
		log.Printf("Login query error: %v", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  map[string]any{"id": user.ID, "email": user.Email, "display_name": displayName},
	})
}

// --- Verify ---

func (a *App) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	var email, displayName string
	if a.DB != nil {
		a.DB.QueryRow(
			"SELECT u.email, p.display_name FROM users u LEFT JOIN profiles p ON p.user_id = u.id WHERE u.id = ?", uid,
		).Scan(&email, &displayName)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user": map[string]any{"id": uid, "email": email, "display_name": displayName},
	})
}
