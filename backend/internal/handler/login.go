package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/KoiralaSam/Mindcare/backend/internal/models/user"
)

// LoginRequest is the JSON body for POST /api/login.
type LoginRequest struct {
	Email  string  `json:"email"`
	Age    *int    `json:"age,omitempty"`
	Gender *string `json:"gender,omitempty"`
}

// UserJSON is an API-friendly view of a user (nullable fields as pointers).
type UserJSON struct {
	ID         int64     `json:"id"`
	Nickname   string    `json:"nickname"`
	Email      string    `json:"email"`
	Age        *int      `json:"age,omitempty"`
	Gender     *string   `json:"gender,omitempty"`
	DailyEmber int       `json:"daily_ember"`
	Streak     int       `json:"streak"`
	Avatar     string    `json:"avatar"`
	CreatedAt  string    `json:"created_at"`
}

// LoginResponse is returned after a successful login or registration.
type LoginResponse struct {
	User    UserJSON `json:"user"`
	Created bool     `json:"created"`
}

func userToJSON(u *user.User) UserJSON {
	out := UserJSON{
		ID:         u.ID,
		Nickname:   u.Nickname,
		Email:      u.Email,
		DailyEmber: u.DailyEmber,
		Streak:     u.Streak,
		Avatar:     u.Avatar,
		CreatedAt:  u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
	if u.Age.Valid {
		v := int(u.Age.Int32)
		out.Age = &v
	}
	if u.Gender.Valid {
		s := u.Gender.String
		out.Gender = &s
	}
	return out
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// LoginHandler handles POST /api/login: existing user by email returns profile; otherwise creates via Save.
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req LoginRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		email := strings.TrimSpace(req.Email)
		if email == "" {
			writeError(w, http.StatusBadRequest, "email is required")
			return
		}

		existing, err := user.GetByEmail(db, email)
		if err == nil {
			writeJSON(w, http.StatusOK, LoginResponse{User: userToJSON(existing), Created: false})
			return
		}
		if !errors.Is(err, user.ErrNotFound) {
			log.Printf("login get by email: %v", err)
			writeError(w, http.StatusInternalServerError, "could not look up user")
			return
		}

		u := user.User{Email: email}
		if req.Age != nil {
			u.Age = sql.NullInt32{Int32: int32(*req.Age), Valid: true}
		}
		if req.Gender != nil {
			g := strings.TrimSpace(*req.Gender)
			if g != "" {
				u.Gender = sql.NullString{String: g, Valid: true}
			}
		}
		if err := u.Save(db); err != nil {
			if errors.Is(err, user.ErrEmailRequired) {
				writeError(w, http.StatusBadRequest, "email is required")
				return
			}
			log.Printf("login save user: %v", err)
			writeError(w, http.StatusInternalServerError, "could not create user")
			return
		}

		writeJSON(w, http.StatusCreated, LoginResponse{User: userToJSON(&u), Created: true})
	}
}
