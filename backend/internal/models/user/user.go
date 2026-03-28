package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
)

// User is a player profile with gamification fields.
// DailyEmber tracks progress toward finishing today’s daily objective (e.g. points or percent scaled by your app).
type User struct {
	ID         int64          `json:"id"`
	Nickname   string         `json:"nickname"`
	Email      string         `json:"email"`
	Age        sql.NullInt32  `json:"age"`
	Gender     sql.NullString `json:"gender"`
	DailyEmber int            `json:"daily_ember"`
	Streak     int            `json:"streak"`
	Avatar     string         `json:"avatar"`
	CreatedAt  time.Time      `json:"created_at"`
}

var (
	ErrEmailRequired   = errors.New("user: email is required")
	ErrNotFound        = errors.New("user: not found")
	ErrNothingToUpdate = errors.New("user: neither streak nor daily ember provided to update")
)

// Save inserts a new user or updates an existing row matched by email.
func (u *User) Save(db *sql.DB) error {
	email := strings.TrimSpace(u.Email)
	if email == "" {
		return ErrEmailRequired
	}

	nick := strings.TrimSpace(u.Nickname)
	if nick == "" {
		gen := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())
		nick = gen.Generate()
	}
	u.Nickname = nick

	var age any
	if u.Age.Valid {
		age = u.Age.Int32
	}

	const q = `
INSERT INTO users (nickname, email, age, daily_ember, streak)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (email) DO UPDATE SET
	nickname = EXCLUDED.nickname,
	age = EXCLUDED.age,
	daily_ember = EXCLUDED.daily_ember,
	streak = EXCLUDED.streak
RETURNING id, created_at, avatar`

	err := db.QueryRow(
		q,
		nick,
		strings.ToLower(email),
		age,
		u.DailyEmber,
		u.Streak,
	).Scan(&u.ID, &u.CreatedAt, &u.Avatar)
	if err != nil {
		return fmt.Errorf("user save: %w", err)
	}
	return nil
}

// GetByEmail returns the user with the given email, or ErrNotFound.
func GetByEmail(db *sql.DB, email string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, ErrEmailRequired
	}

	const q = `
SELECT id, nickname, email, age, daily_ember, streak, avatar, created_at
FROM users
WHERE email = $1`

	var u User
	err := db.QueryRow(q, email).Scan(
		&u.ID,
		&u.Nickname,
		&u.Email,
		&u.Age,
		&u.DailyEmber,
		&u.Streak,
		&u.Avatar,
		&u.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user by email: %w", err)
	}
	return &u, nil
}

// UpdateStreakAndEmber updates daily_ember and/or streak for the row matching Email.
// Pass nil for a pointer to skip that column; 0 is valid when the pointer is non-nil.
// Returns ErrNothingToUpdate if both arguments are nil.
func (u *User) UpdateStreakAndEmber(db *sql.DB, dailyEmber *int, streak *int) error {
	email := strings.TrimSpace(strings.ToLower(u.Email))
	if email == "" {
		return ErrEmailRequired
	}
	if dailyEmber == nil && streak == nil {
		return ErrNothingToUpdate
	}

	var parts []string
	var args []any
	n := 1
	if dailyEmber != nil {
		parts = append(parts, fmt.Sprintf("daily_ember = $%d", n))
		args = append(args, *dailyEmber)
		n++
	}
	if streak != nil {
		parts = append(parts, fmt.Sprintf("streak = $%d", n))
		args = append(args, *streak)
		n++
	}
	args = append(args, email)
	q := fmt.Sprintf("UPDATE users SET %s WHERE email = $%d", strings.Join(parts, ", "), n)

	res, err := db.Exec(q, args...)
	if err != nil {
		return fmt.Errorf("user update streak/ember: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("user update streak/ember: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
