package data

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindAll() ([]User, error) {
	var users []User
	err := r.db.Where("deleted_at IS NULL").Find(&users).Error
	return users, err
}

func (r *UserRepo) FindByID(id uuid.UUID) (*User, error) {
	var user User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepo) FindByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepo) Create(u *User) error {
	u.ID = uuid.New()
	return r.db.Create(u).Error
}

func (r *UserRepo) Update(u *User) error {
	return r.db.Save(u).Error
}

func (r *UserRepo) Delete(id uuid.UUID) error {
	return r.db.Model(&User{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
