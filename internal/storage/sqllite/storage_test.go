package storage

import (
	"context"
	"fmt"
	"testing"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/GusevGrishaEm1/security-service/internal/model"
	"github.com/stretchr/testify/assert"
)

// SetupDatabase creates a new SQLite database and a users table.
func SetupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}
func TestFindUserByEmail(t *testing.T) {
	db, err := SetupDatabase()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	store := &storage{db: db}

	ctx := context.Background()

	// Insert a test user
	testUser := model.User{
		Email:    "test@example.com",
		Password: "password123",
	}
	insertUserQuery := `INSERT INTO users (email, password) VALUES (?, ?)`
	_, err = db.ExecContext(ctx, insertUserQuery, testUser.Email, testUser.Password)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Test finding the user by email
	user, err := store.FindUserByEmail(ctx, testUser.Email)
	if err != nil {
		t.Fatalf("Failed to find user by email: %v", err)
	}

	assert.Equal(t, testUser.Email, user.Email, "Expected user email to match")
	assert.Equal(t, testUser.Password, user.Password, "Expected user password to match")
}

func TestSaveUser(t *testing.T) {
	db, err := SetupDatabase()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	defer db.Close()

	store := &storage{db: db}

	ctx := context.Background()

	// Test user to save
	testUser := model.User{
		Email:    "newuser@example.com",
		Password: "newpassword123",
	}

	// Save the user
	err = store.SaveUser(ctx, testUser)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Verify the user was saved
	row := db.QueryRowContext(ctx, `SELECT email, password FROM users WHERE email = ?`, testUser.Email)

	var savedUser model.User
	err = row.Scan(&savedUser.Email, &savedUser.Password)
	if err != nil {
		t.Fatalf("Failed to query saved user: %v", err)
	}

	assert.Equal(t, testUser.Email, savedUser.Email, "Expected user email to match")
	assert.Equal(t, testUser.Password, savedUser.Password, "Expected user password to match")
}
