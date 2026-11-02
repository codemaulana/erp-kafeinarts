package repository

import (
	"errors"

	"erp-kafeinart/model"

	"gorm.io/gorm"
)

// UserRepository mendefinisikan kontrak akses data untuk entitas User.
type UserRepository interface {
	FindAll(search, role string) ([]model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository menginisialisasi implementasi UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(search, role string) ([]model.User, error) {
	var users []model.User
	query := r.db.Model(&model.User{})
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if err := query.Order("id ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
