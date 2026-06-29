package service

import (
	"encoding/json"
	"net/http"

	v1 "toko-emas/api/v1"
	"toko-emas/internal/data"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *data.UserRepo
}

func NewUserService(repo *data.UserRepo) *UserService {
	return &UserService{repo: repo}
}

// ListUsers godoc
// @Summary      List semua user
// @Description  Mendapatkan semua data user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  v1.Response
// @Router       /api/v1/users [get]
func (s *UserService) List(w http.ResponseWriter, r *http.Request) {
	users, err := s.repo.FindAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: users})
}

// GetUser godoc
// @Summary      Get user by ID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/users/{id} [get]
func (s *UserService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	user, err := s.repo.FindByID(id)
	if err != nil || user == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: user})
}

// CreateUser godoc
// @Summary      Tambah user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body  body      v1.CreateUserRequest  true  "Data user"
// @Success      201   {object}  v1.Response
// @Router       /api/v1/users [post]
func (s *UserService) Create(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: "hash error"})
		return
	}
	u := &data.User{
		Nama:     req.Nama,
		Username: req.Username,
		Password: string(hash),
		Role:     req.Role,
	}
	if err := s.repo.Create(u); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, v1.Response{Code: 201, Message: "created", Data: u})
}

// UpdateUser godoc
// @Summary      Update user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "User ID"
// @Param        body  body      v1.UpdateUserRequest  true  "Data user"
// @Success      200   {object}  v1.Response
// @Router       /api/v1/users/{id} [put]
func (s *UserService) Update(w http.ResponseWriter, r *http.Request) {
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
	var req v1.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	existing.Nama = req.Nama
	existing.Username = req.Username
	existing.Role = req.Role
	if err := s.repo.Update(existing); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "updated", Data: existing})
}

// DeleteUser godoc
// @Summary      Hapus user
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/users/{id} [delete]
func (s *UserService) Delete(w http.ResponseWriter, r *http.Request) {
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
