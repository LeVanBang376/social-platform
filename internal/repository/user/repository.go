package user

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"social-platform/internal/model"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(
	ctx context.Context,
	db *gorm.DB,
	user *model.User,
) error {
	return db.
		WithContext(ctx).
		Create(user).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	userID uuid.UUID,
) (*model.User, error) {
	var user model.User

	err := db.
		WithContext(ctx).
		First(&user, "user_id = ?", userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByEmail(
	ctx context.Context,
	db *gorm.DB,
	email string,
) (*model.User, error) {
	var user model.User

	err := db.
		WithContext(ctx).
		First(&user, "email = ?", email).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	user *model.User,
) error {
	return db.
		WithContext(ctx).
		Save(user).
		Error
}
