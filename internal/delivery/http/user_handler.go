package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/LuigiAzevedo/public-library-v2/internal/domain/entity"
	uc "github.com/LuigiAzevedo/public-library-v2/internal/ports/usecase"
)

type userHandler struct {
	UserUsecase uc.UserUsecase
}

// NewUserHandler creates a new instance of userHandler
func NewUserHandler(r *chi.Mux, useCase uc.UserUsecase) {
	handler := &userHandler{
		UserUsecase: useCase,
	}

	r.Route("/v1/users", func(r chi.Router) {
		r.Get("/{id}", handler.GetUser)
		r.Post("/", handler.CreateUser)
		r.Put("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
	})
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	u, err := h.UserUsecase.GetUser(ctx, id)
	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	json.NewEncoder(w).Encode(u)
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u *entity.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := h.UserUsecase.CreateUser(ctx, u)
	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u *entity.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.ID, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.UserUsecase.UpdateUser(ctx, u)
	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request timed out", http.StatusGatewayTimeout)
		default:
			if errors.Cause(err).Error() == "user not found" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.UserUsecase.DeleteUser(ctx, id)
	if err != nil {
		if errors.Cause(err).Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
