package repositories

import (
	"chords_app/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
