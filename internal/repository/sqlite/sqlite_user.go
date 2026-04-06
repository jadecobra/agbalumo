package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func scanUser(s Scanner) (domain.User, error) {
	var u domain.User
	var createdAt time.Time
	err := s.Scan(&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.AvatarURL, &u.Role, &createdAt)
	if err == nil {
		u.CreatedAt = createdAt
	}
	return u, err
}

// SaveUser inserts or updates a user.
func (r *SQLiteRepository) SaveUser(ctx context.Context, u domain.User) error {
	updateQuery := `UPDATE users SET google_id=?, email=?, name=?, avatar_url=?, role=? WHERE id=?`
	res, err := r.writeDB.ExecContext(ctx, updateQuery,
		u.GoogleID, u.Email, u.Name, u.AvatarURL, u.Role, u.ID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows > 0 {
		return nil
	}

	insertQuery := `
	INSERT INTO users (id, google_id, email, name, avatar_url, role, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(google_id) DO UPDATE SET
		email = excluded.email,
		name = excluded.name,
		avatar_url = excluded.avatar_url,
		role = excluded.role;
	`
	_, err = r.writeDB.ExecContext(ctx, insertQuery,
		u.ID, u.GoogleID, u.Email, u.Name, u.AvatarURL, u.Role, u.CreatedAt,
	)
	return err
}

// FindUserByGoogleID retrieves a user by their Google ID.
func (r *SQLiteRepository) FindUserByGoogleID(ctx context.Context, googleID string) (domain.User, error) {
	query := `SELECT id, google_id, email, name, avatar_url, COALESCE(role, 'User'), created_at FROM users WHERE google_id = ?`
	row := r.readDB.QueryRowContext(ctx, query, googleID)

	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("user not found")
	}
	return u, err
}

// FindUserByID retrieves a user by their ID.
func (r *SQLiteRepository) FindUserByID(ctx context.Context, id string) (domain.User, error) {
	query := `SELECT id, google_id, email, name, avatar_url, COALESCE(role, 'User'), created_at FROM users WHERE id = ?`
	row := r.readDB.QueryRowContext(ctx, query, id)

	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("user not found")
	}
	return u, err
}

func (r *SQLiteRepository) GetUserCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	if err := r.readDB.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *SQLiteRepository) GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error) {
	query := `SELECT id, google_id, email, name, avatar_url, COALESCE(role, 'User'), created_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.readDB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var users []domain.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
