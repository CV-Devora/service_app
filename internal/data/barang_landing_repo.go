package data

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BarangLandingRepo struct {
	db *gorm.DB
}

func NewBarangLandingRepo(db *gorm.DB) *BarangLandingRepo {
	return &BarangLandingRepo{db: db}
}

func (r *BarangLandingRepo) FindAll() ([]BarangLanding, error) {
	var items []BarangLanding
	err := r.db.Where("deleted_at IS NULL").Order("created_at ASC").Find(&items).Error
	return items, err
}

func (r *BarangLandingRepo) FindByID(id uuid.UUID) (*BarangLanding, error) {
	var item BarangLanding
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (r *BarangLandingRepo) Create(item *BarangLanding) error {
	item.ID = uuid.New()
	return r.db.Create(item).Error
}

func (r *BarangLandingRepo) Update(item *BarangLanding) error {
	return r.db.Save(item).Error
}

func (r *BarangLandingRepo) Delete(id uuid.UUID) error {
	return r.db.Model(&BarangLanding{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
