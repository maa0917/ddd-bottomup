package presentation

import (
	"ddd-bottomup/usecase"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	createUserUseCase *usecase.CreateUserUseCase
	getUserUseCase    *usecase.GetUserUseCase
	updateUserUseCase *usecase.UpdateUserUseCase
	deleteUserUseCase *usecase.DeleteUserUseCase
}

func NewUserHandler(
	createUserUseCase *usecase.CreateUserUseCase,
	getUserUseCase *usecase.GetUserUseCase,
	updateUserUseCase *usecase.UpdateUserUseCase,
	deleteUserUseCase *usecase.DeleteUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
		getUserUseCase:    getUserUseCase,
		updateUserUseCase: updateUserUseCase,
		deleteUserUseCase: deleteUserUseCase,
	}
}

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateUserResponse struct {
	UserID string `json:"userId"`
}

type GetUserResponse struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Email     *string `json:"email,omitempty"`
}

type UpdateUserResponse struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := usecase.CreateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	output, err := h.createUserUseCase.Execute(input)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "already exists") ||
			strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "cannot be empty") {
			status = http.StatusBadRequest
		}
		h.writeError(w, err.Error(), status)
		return
	}

	response := CreateUserResponse{
		UserID: output.UserID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "userID")
	input := usecase.GetUserInput{
		UserID: userID,
	}

	output, err := h.getUserUseCase.Execute(input)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "invalid") {
			status = http.StatusBadRequest
		}
		h.writeError(w, err.Error(), status)
		return
	}

	response := GetUserResponse{
		UserID:    output.UserID,
		FirstName: output.FirstName,
		LastName:  output.LastName,
		Email:     output.Email,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "userID")
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := usecase.UpdateUserInput{
		UserID:    userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	output, err := h.updateUserUseCase.Execute(input)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "already exists") ||
			strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "cannot be empty") {
			status = http.StatusBadRequest
		}
		h.writeError(w, err.Error(), status)
		return
	}

	response := UpdateUserResponse{
		UserID:    output.UserID,
		FirstName: output.FirstName,
		LastName:  output.LastName,
		Email:     output.Email,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "userID")
	input := usecase.DeleteUserInput{
		UserID: userID,
	}

	err := h.deleteUserUseCase.Execute(input)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "invalid") {
			status = http.StatusBadRequest
		}
		h.writeError(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) writeError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
