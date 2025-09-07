package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/arieffadhlan/go-fitbyte/internal/models"
)

type UserRepository struct {
  db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) Repository {
  return &UserRepository{
    db: db,
  }
}

func (r *UserRepository) GetUserById(c context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, password, name, preference, weight_unit, height_unit, 
		       weight, height, image_uri, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user models.User
	err := r.db.GetContext(c, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserImage(c context.Context, id int, imageUri string) error {
  query := `UPDATE users SET image_uri = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
  
  _, err := r.db.ExecContext(c, query, imageUri, id)
	if err != nil {
		 return fmt.Errorf("failed to update user image: %w", err)
	}
  
	return nil
}