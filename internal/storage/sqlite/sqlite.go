package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/storage"

	//"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/mattn/go-sqlite3"
)



type Storage struct {
	db *sql.DB

}


func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db : db}, nil

}


func (a *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {

	const op = "storage.sqlite.SaveUser" 

	statement, err := a.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := statement.ExecContext(ctx, email, passHash)

	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}



func (a *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	statement, err := a.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email == ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s %w", op, err)
	}

	row := statement.QueryRowContext(ctx, email)

	var user models.User

	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
		
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}


func (a *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {

	const op = "storage.sqlite.IsAdmin"


	statement, err := a.db.Prepare("SELECT is_admin FROM users WHERE id == ?")

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := statement.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil

}


func (a *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "sqlite.Storage.App"

	statement, err := a.db.Prepare("SELECT id, name, secret FROM apps WHERE id == ?")

	if err != nil {
		return models.App{}, fmt.Errorf("%s %w", op, err)
	}

	var app models.App

	row := statement.QueryRowContext(ctx, appID)

	err = row.Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}