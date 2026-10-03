package data

import (
	"gorm.io/gorm"
)

type DashboardRepo struct {
	db *gorm.DB
}

func NewDashboardRepo(db *gorm.DB) *DashboardRepo {
	return &DashboardRepo{db: db}
}

type ChartResult struct {
	Label string
	Value int64
}

func (r *DashboardRepo) GetBarangChart() ([]ChartResult, error) {
	var results []ChartResult
	err := r.db.Model(&Barang{}).
		Joins("LEFT JOIN karat k ON k.id = barang.karat_id").
		Select("COALESCE(k.name, 'Tanpa Karat') as label, COUNT(barang.id) as value").
		Where("barang.deleted_at IS NULL").
		Group("k.name").
		Find(&results).Error
	return results, err
}

func (r *DashboardRepo) GetPenjualanChart() ([]ChartResult, error) {
	var results []ChartResult
	err := r.db.Model(&Penjualan{}).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as label, SUM(total_harga) as value").
		Where("deleted_at IS NULL").
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("label ASC").
		Find(&results).Error
	return results, err
}

func (r *DashboardRepo) GetPembelianChart() ([]ChartResult, error) {
	var results []ChartResult
	err := r.db.Model(&Pembelian{}).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as label, SUM(harga_deal) as value").
		Where("deleted_at IS NULL").
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("label ASC").
		Find(&results).Error
	return results, err
}

// GetPenjualanHariIni returns total sales amount of today.
func (r *DashboardRepo) GetPenjualanHariIni() (int64, error) {
	var total int64
	err := r.db.Model(&Penjualan{}).
		Select("COALESCE(SUM(total_harga), 0)").
		Where("deleted_at IS NULL AND created_at::date = CURRENT_DATE").
		Scan(&total).Error
	return total, err
}

// GetPenjualanHariIniBySales returns total sales amount of today for a specific sales code.
func (r *DashboardRepo) GetPenjualanHariIniBySales(kodeSales int64) (int64, error) {
	var total int64
	err := r.db.Model(&Penjualan{}).
		Select("COALESCE(SUM(total_harga), 0)").
		Where("deleted_at IS NULL AND created_at::date = CURRENT_DATE AND kode_sales = ?", kodeSales).
		Scan(&total).Error
	return total, err
}

// GetPenjualanBulananBySales returns total sales amount of this month for a specific sales code.
func (r *DashboardRepo) GetPenjualanBulananBySales(kodeSales int64) (int64, error) {
	var total int64
	err := r.db.Model(&Penjualan{}).
		Select("COALESCE(SUM(total_harga), 0)").
		Where("deleted_at IS NULL AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE) AND kode_sales = ?", kodeSales).
		Scan(&total).Error
	return total, err
}

// CountBarangTersedia returns the number of available (unsold) items.
func (r *DashboardRepo) CountBarangTersedia() (int64, error) {
	var count int64
	err := r.db.Model(&Barang{}).
		Where("deleted_at IS NULL AND penjualan_id IS NULL").
		Count(&count).Error
	return count, err
}

// GetPenjualanBulananChart returns monthly sales totals for the last 12 months.
func (r *DashboardRepo) GetPenjualanBulananChart() ([]ChartResult, error) {
	var results []ChartResult
	err := r.db.Model(&Penjualan{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as label, SUM(total_harga) as value").
		Where("deleted_at IS NULL AND created_at >= date_trunc('month', CURRENT_DATE) - INTERVAL '11 months'").
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("label ASC").
		Find(&results).Error
	return results, err
}

// GetPenjualanGrupChart returns the count of sold items grouped by grup.
func (r *DashboardRepo) GetPenjualanGrupChart() ([]ChartResult, error) {
	var results []ChartResult
	err := r.db.Model(&Barang{}).
		Joins("JOIN penjualan p ON p.id = barang.penjualan_id").
		Select("COALESCE(NULLIF(barang.grup, ''), 'Tanpa Grup') as label, COUNT(barang.id) as value").
		Where("barang.deleted_at IS NULL AND p.deleted_at IS NULL").
		Group("COALESCE(NULLIF(barang.grup, ''), 'Tanpa Grup')").
		Order("value DESC").
		Find(&results).Error
	return results, err
}

// GetPenjualanTahunanChartBySales returns monthly sales totals for the current year for a specific sales code.
func (r *DashboardRepo) GetPenjualanTahunanChartBySales(kodeSales int64) ([]ChartResult, error) {
	var results []ChartResult
	err := r.db.Model(&Penjualan{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as label, SUM(total_harga) as value").
		Where("deleted_at IS NULL AND EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM CURRENT_DATE) AND kode_sales = ?", kodeSales).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("label ASC").
		Find(&results).Error
	return results, err
}
