package service

import (
	"net/http"
	"strconv"
	
	v1 "toko-emas/api/v1"
	"toko-emas/internal/data"
)

type DashboardService struct {
	repo *data.DashboardRepo
}

func NewDashboardService(repo *data.DashboardRepo) *DashboardService {
	return &DashboardService{repo: repo}
}

// GetDashboard godoc
// @Summary      Get dashboard charts data
// @Tags         Dashboard
// @Produce      json
// @Success      200  {object}  v1.Response
// @Router       /api/v1/dashboard [get]
func (s *DashboardService) GetDashboard(w http.ResponseWriter, r *http.Request) {
	barangChart, err := s.repo.GetBarangChart()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	penjualanChart, err := s.repo.GetPenjualanChart()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	pembelianChart, err := s.repo.GetPembelianChart()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	penjualanHariIni, err := s.repo.GetPenjualanHariIni()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	barangTersedia, err := s.repo.CountBarangTersedia()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	penjualanBulanan, err := s.repo.GetPenjualanBulananChart()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	penjualanPerGrup, err := s.repo.GetPenjualanGrupChart()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	var penjualanHariIniSales int64
	var penjualanBulananSales int64
	var penjualanTahunanSalesChart []data.ChartResult

	kodeSalesStr := r.URL.Query().Get("kode_sales")
	if kodeSalesStr != "" {
		if kodeSales, err := strconv.ParseInt(kodeSalesStr, 10, 64); err == nil {
			penjualanHariIniSales, _ = s.repo.GetPenjualanHariIniBySales(kodeSales)
			penjualanBulananSales, _ = s.repo.GetPenjualanBulananBySales(kodeSales)
			penjualanTahunanSalesChart, _ = s.repo.GetPenjualanTahunanChartBySales(kodeSales)
		}
	}

	resp := v1.DashboardResponse{
		BarangChart:                make([]v1.ChartData, len(barangChart)),
		PenjualanChart:             make([]v1.ChartData, len(penjualanChart)),
		PembelianChart:             make([]v1.ChartData, len(pembelianChart)),
		PenjualanHariIni:           penjualanHariIni,
		BarangTersedia:             barangTersedia,
		PenjualanBulanan:           make([]v1.ChartData, len(penjualanBulanan)),
		PenjualanPerGrup:           make([]v1.ChartData, len(penjualanPerGrup)),
		PenjualanHariIniSales:      penjualanHariIniSales,
		PenjualanBulananSales:      penjualanBulananSales,
		PenjualanTahunanSalesChart: make([]v1.ChartData, len(penjualanTahunanSalesChart)),
	}

	for i, b := range barangChart {
		resp.BarangChart[i] = v1.ChartData{Label: b.Label, Value: b.Value}
	}
	for i, p := range penjualanChart {
		resp.PenjualanChart[i] = v1.ChartData{Label: p.Label, Value: p.Value}
	}
	for i, p := range pembelianChart {
		resp.PembelianChart[i] = v1.ChartData{Label: p.Label, Value: p.Value}
	}
	for i, b := range penjualanBulanan {
		resp.PenjualanBulanan[i] = v1.ChartData{Label: b.Label, Value: b.Value}
	}
	for i, g := range penjualanPerGrup {
		resp.PenjualanPerGrup[i] = v1.ChartData{Label: g.Label, Value: g.Value}
	}
	for i, c := range penjualanTahunanSalesChart {
		resp.PenjualanTahunanSalesChart[i] = v1.ChartData{Label: c.Label, Value: c.Value}
	}

	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: resp})
}
