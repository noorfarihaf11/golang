package repository

import (
	"database/sql"
	"time"
	"github.com/noorfarihaf11/clean-arc/app/model"
)

func RegisterUser(db *sql.DB, user *model.User) (*model.User, error) {
	query := `
        INSERT INTO users (username, email, password_hash, role, created_at) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, username, email, role, created_at
    `

	err := db.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		time.Now(),
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
