package server

import (
	"net/http"

	"toko-emas/internal/service"

	_ "toko-emas/docs" // swagger docs

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewHTTPRouter sets up all routes
func NewHTTPRouter(
	barangSvc *service.BarangService,
	userSvc *service.UserService,
	authSvc *service.AuthService,
	pembelianSvc *service.PembelianService,
	karatSvc *service.KaratService,
	bakiSvc *service.BakiService,
	penjualanSvc *service.PenjualanService,
) http.Handler {
	r := mux.NewRouter()

	// Swagger UI
	r.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth
	api.HandleFunc("/auth/login", authSvc.Login).Methods(http.MethodPost)
	api.HandleFunc("/auth/refresh", authSvc.Refresh).Methods(http.MethodPost)

	// Barang
	api.HandleFunc("/barang", barangSvc.List).Methods(http.MethodGet)
	api.HandleFunc("/barang/{id}", barangSvc.Get).Methods(http.MethodGet)
	api.HandleFunc("/barang", barangSvc.Create).Methods(http.MethodPost)
	api.HandleFunc("/barang/{id}", barangSvc.Update).Methods(http.MethodPut)
	api.HandleFunc("/barang/{id}", barangSvc.Delete).Methods(http.MethodDelete)

	// Users
	api.HandleFunc("/users", userSvc.List).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}", userSvc.Get).Methods(http.MethodGet)
	api.HandleFunc("/users", userSvc.Create).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", userSvc.Update).Methods(http.MethodPut)
	api.HandleFunc("/users/{id}", userSvc.Delete).Methods(http.MethodDelete)

	// Pembelian
	api.HandleFunc("/pembelian", pembelianSvc.List).Methods(http.MethodGet)
	api.HandleFunc("/pembelian/{id}", pembelianSvc.Get).Methods(http.MethodGet)
	api.HandleFunc("/pembelian", pembelianSvc.Create).Methods(http.MethodPost)
	api.HandleFunc("/pembelian/{id}", pembelianSvc.Update).Methods(http.MethodPut)
	api.HandleFunc("/pembelian/{id}", pembelianSvc.Delete).Methods(http.MethodDelete)

	// Karat
	api.HandleFunc("/karat", karatSvc.List).Methods(http.MethodGet)
	api.HandleFunc("/karat/{id}", karatSvc.Get).Methods(http.MethodGet)
	api.HandleFunc("/karat", karatSvc.Create).Methods(http.MethodPost)
	api.HandleFunc("/karat/{id}", karatSvc.Update).Methods(http.MethodPut)
	api.HandleFunc("/karat/{id}", karatSvc.Delete).Methods(http.MethodDelete)

	// Baki
	api.HandleFunc("/baki", bakiSvc.List).Methods(http.MethodGet)
	api.HandleFunc("/baki/{id}", bakiSvc.Get).Methods(http.MethodGet)
	api.HandleFunc("/baki", bakiSvc.Create).Methods(http.MethodPost)
	api.HandleFunc("/baki/{id}", bakiSvc.Update).Methods(http.MethodPut)
	api.HandleFunc("/baki/{id}", bakiSvc.Delete).Methods(http.MethodDelete)

	// Penjualan
	api.HandleFunc("/penjualan", penjualanSvc.List).Methods(http.MethodGet)
	api.HandleFunc("/penjualan/{id}", penjualanSvc.Get).Methods(http.MethodGet)
	api.HandleFunc("/penjualan", penjualanSvc.Create).Methods(http.MethodPost)
	api.HandleFunc("/penjualan/{id}", penjualanSvc.Update).Methods(http.MethodPut)
	api.HandleFunc("/penjualan/{id}", penjualanSvc.Delete).Methods(http.MethodDelete)
	api.HandleFunc("/penjualan/{id}/barang", penjualanSvc.AttachBarang).Methods(http.MethodPost)

	return r
}
