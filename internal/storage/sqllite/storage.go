package storage

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"

	"github.com/GusevGrishaEm1/security-service/internal/config"
	"github.com/GusevGrishaEm1/security-service/internal/model"
)

type storage struct {
	db *sql.DB
}

func NewAuthStorage(config *config.Config) (*storage, error) {
	db, err := sql.Open("sqlite3", config.StoragePath)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatal(err)
	}
	return &storage{db: db}, nil
}

func (s *storage) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = ?`
	row := s.db.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, nil
		}
		return model.User{}, err
	}

	return user, nil
}

func (s *storage) SaveUser(ctx context.Context, user model.User) error {
	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, query, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
