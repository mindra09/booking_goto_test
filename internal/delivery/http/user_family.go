package http

import (
	"booking_togo/internal/model"
	"booking_togo/internal/usecase"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type UserFamilyHandler struct {
	usecaseuser usecase.IUserUsecase
}

func NewUserFamilyHandler(usecaseuser usecase.IUserUsecase) *UserFamilyHandler {
	return &UserFamilyHandler{
		usecaseuser: usecaseuser,
	}

}

func (h *UserFamilyHandler) RegisterRoutes(r *mux.Router) {
	// Register routes related to user family here
	r.HandleFunc("/user", h.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/user", h.CreateUserFamily).Methods(http.MethodPost)
	r.HandleFunc("/user/{id}", h.UserDetail).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", h.UserUpdate).Methods(http.MethodPut)
	r.HandleFunc("/user/{id}", h.UserDelete).Methods(http.MethodDelete)
	r.HandleFunc("/user/{id}/family/{family_id}", h.FamilyDelete).Methods(http.MethodDelete)
}

func (h *UserFamilyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	users, err := h.usecaseuser.GetAll(ctx)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (h *UserFamilyHandler) CreateUserFamily(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var userFamilyPayload model.User
	if err := json.NewDecoder(r.Body).Decode(&userFamilyPayload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := h.usecaseuser.Create(ctx, &userFamilyPayload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, "User created successfully")
}

func (h *UserFamilyHandler) UserDetail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userDetail, userDetailErr := h.usecaseuser.Detail(ctx, id)
	if userDetailErr != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": userDetailErr.Error()})
		return
	}
	writeJSON(w, http.StatusOK, userDetail)
}

func (h *UserFamilyHandler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var userFamilyPayload model.User
	if err := json.NewDecoder(r.Body).Decode(&userFamilyPayload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userFamilyPayload.UserID = id
	userUpdatelErr := h.usecaseuser.Update(ctx, &userFamilyPayload)
	if userUpdatelErr != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": userUpdatelErr.Error()})
		return
	}
	writeJSON(w, http.StatusOK, "User updated successfully")
}

func (h *UserFamilyHandler) UserDelete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userDetailErr := h.usecaseuser.Delete(ctx, id)
	if userDetailErr != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": userDetailErr.Error()})
		return
	}
	writeJSON(w, http.StatusOK, "User deleted successfully")
}

func (h *UserFamilyHandler) FamilyDelete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	familyID, err := strconv.Atoi(mux.Vars(r)["family_id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	familyDeleteErr := h.usecaseuser.DeleteFamily(ctx, id, familyID)
	if familyDeleteErr != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": familyDeleteErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, "family deleted successfully")
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
