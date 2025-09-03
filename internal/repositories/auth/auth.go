package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) Repository {
	return &AuthRepository{db: db}
}

func (ar *AuthRepository) Create(ctx context.Context, user *models.User) (int, error) {
	var newId int
	currentTime := time.Now()

	query := `
		INSERT INTO users (email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := ar.db.QueryRowContext(ctx, query,
		user.Email,
		user.Password,
		currentTime,
		currentTime,
	).Scan(&newId)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, errors.New("email already exists")
		} else {
			return 0, errors.New("register user failed")
		}
	}

	return newId, nil
}

func (ar *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	uData := &models.User{}
	query := `SELECT id, email, password FROM users WHERE email = $1`

	err := ar.db.QueryRowContext(ctx,
		query,
		email,
	).Scan(&uData.ID, &uData.Email, &uData.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("user not found")
	}

	return uData, nil
}
