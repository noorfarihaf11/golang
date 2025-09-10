package service

import (
	"database/sql"
	"errors"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/utils"
)

func LoginService(db *sql.DB, req model.LoginRequest) (string, model.User, error) {
	var user model.User
	var passwordHash string

	err := db.QueryRow(`
		SELECT id, username, email, password_hash, role, created_at
		FROM users
		WHERE username = $1 OR email = $1
	`, req.Username).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		return "", user, errors.New("username atau password salah")
	}

	// cek password
	if !utils.CheckPassword(req.Password, passwordHash) {
		return "", user, errors.New("password salah")
	}

	// generate JWT
	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", user, errors.New("gagal generate token")
	}

	return token, user, nil
}
