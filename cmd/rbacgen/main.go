package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := envOrDefault("DB_HOST", "edu-planner-user-management-postgres")
	dbPort := envOrDefault("DB_PORT", "5432")

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminFirstName := envOrDefault("ADMIN_FIRST_NAME", "Admin")
	adminLastName := envOrDefault("ADMIN_LAST_NAME", "User")

	if adminEmail == "" || adminPassword == "" {
		return errors.New("ADMIN_EMAIL and ADMIN_PASSWORD must be set")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/user_management?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return fmt.Errorf("connect db: %w", err)
	}

	var count int
	if err = db.QueryRow(`SELECT COUNT(*) FROM "user" WHERE email = $1 AND is_deleted = false`, adminEmail).Scan(&count); err != nil {
		return fmt.Errorf("check admin: %w", err)
	}
	if count > 0 {
		log.Printf("admin %s already exists, skipping", adminEmail)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var userID uint64
	if err = tx.QueryRow(`
		INSERT INTO "user" (email, role, is_active, created_at)
		VALUES ($1, 'admin', true, $2)
		RETURNING id`,
		adminEmail, now,
	).Scan(&userID); err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	if _, err = tx.Exec(`
		INSERT INTO user_profile (user_id, first_name, last_name, created_at)
		VALUES ($1, $2, $3, $4)`,
		userID, adminFirstName, adminLastName, now,
	); err != nil {
		return fmt.Errorf("insert user_profile: %w", err)
	}

	if _, err = tx.Exec(`
		INSERT INTO user_credential (user_id, password_hash, updated_at)
		VALUES ($1, $2, $3)`,
		userID, string(hash), now,
	); err != nil {
		return fmt.Errorf("insert user_credential: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	log.Printf("admin user %s created", adminEmail)
	return nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
