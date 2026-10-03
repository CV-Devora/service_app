package service

import (
	"encoding/json"
	"net/http"

	v1 "toko-emas/api/v1"
	"toko-emas/internal/data"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ---- Barang Landing ----

type BarangLandingService struct {
	repo *data.BarangLandingRepo
}

func NewBarangLandingService(repo *data.BarangLandingRepo) *BarangLandingService {
	return &BarangLandingService{repo: repo}
}

// ListBarangLanding godoc
// @Summary      List semua barang landing page
// @Tags         Barang Landing
// @Produce      json
// @Success      200  {object}  v1.Response
// @Router       /api/v1/barang-landing [get]
func (s *BarangLandingService) List(w http.ResponseWriter, r *http.Request) {
	items, err := s.repo.FindAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: items})
}

// GetBarangLanding godoc
// @Summary      Get barang landing by ID
// @Tags         Barang Landing
// @Produce      json
// @Param        id   path      string  true  "Barang Landing ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/barang-landing/{id} [get]
func (s *BarangLandingService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	item, err := s.repo.FindByID(id)
	if err != nil || item == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: item})
}

// CreateBarangLanding godoc
// @Summary      Tambah barang landing
// @Tags         Barang Landing
// @Accept       json
// @Produce      json
// @Param        body  body      v1.CreateBarangLandingRequest  true  "Data barang landing"
// @Success      201   {object}  v1.Response
// @Router       /api/v1/barang-landing [post]
func (s *BarangLandingService) Create(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateBarangLandingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	item := &data.BarangLanding{
		Nama:  req.Nama,
		Karat: req.Karat,
		Berat: req.Berat,
		Harga: req.Harga,
		Photo: req.Photo,
	}
	if err := s.repo.Create(item); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, v1.Response{Code: 201, Message: "created", Data: item})
}

// UpdateBarangLanding godoc
// @Summary      Update barang landing
// @Tags         Barang Landing
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "Barang Landing ID"
// @Param        body  body      v1.UpdateBarangLandingRequest  true  "Data barang landing"
// @Success      200   {object}  v1.Response
// @Router       /api/v1/barang-landing/{id} [put]
func (s *BarangLandingService) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	existing, err := s.repo.FindByID(id)
	if err != nil || existing == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	var req v1.UpdateBarangLandingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	existing.Nama = req.Nama
	existing.Karat = req.Karat
	existing.Berat = req.Berat
	existing.Harga = req.Harga
	existing.Photo = req.Photo
	if err := s.repo.Update(existing); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "updated", Data: existing})
}

// DeleteBarangLanding godoc
// @Summary      Hapus barang landing
// @Tags         Barang Landing
// @Produce      json
// @Param        id   path      string  true  "Barang Landing ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/barang-landing/{id} [delete]
func (s *BarangLandingService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	if err := s.repo.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "deleted"})
}
