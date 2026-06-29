package data

import (
	"errors"

	"gorm.io/gorm"
)

type PembelianRepo struct {
	db *gorm.DB
}

func NewPembelianRepo(db *gorm.DB) *PembelianRepo {
	return &PembelianRepo{db: db}
}

func (r *PembelianRepo) FindAll() ([]Pembelian, error) {
	var items []Pembelian
	err := r.db.Where("deleted_at IS NULL").Find(&items).Error
	return items, err
}

func (r *PembelianRepo) FindByID(id uint) (*Pembelian, error) {
	var item Pembelian
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (r *PembelianRepo) Create(p *Pembelian) error {
	return r.db.Create(p).Error
}

func (r *PembelianRepo) Update(p *Pembelian) error {
	return r.db.Save(p).Error
}

func (r *PembelianRepo) Delete(id uint) error {
	return r.db.Model(&Pembelian{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
