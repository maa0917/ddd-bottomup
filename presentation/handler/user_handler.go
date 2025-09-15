package handler

import (
	"ddd-bottomup/usecase"
	"encoding/json"
	"net/http"
	"strings"
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

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	path := strings.TrimPrefix(r.URL.Path, "/users")
	
	switch r.Method {
	case http.MethodPost:
		if path == "" || path == "/" {
			h.CreateUser(w, r)
		} else {
			h.writeError(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodGet:
		if path == "" || path == "/" {
			h.writeError(w, "User ID required", http.StatusBadRequest)
		} else {
			userID := strings.TrimPrefix(path, "/")
			h.GetUser(w, r, userID)
		}
	case http.MethodPut:
		if path == "" || path == "/" {
			h.writeError(w, "User ID required", http.StatusBadRequest)
		} else {
			userID := strings.TrimPrefix(path, "/")
			h.UpdateUser(w, r, userID)
		}
	case http.MethodDelete:
		if path == "" || path == "/" {
			h.writeError(w, "User ID required", http.StatusBadRequest)
		} else {
			userID := strings.TrimPrefix(path, "/")
			h.DeleteUser(w, r, userID)
		}
	default:
		h.writeError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, userID string) {
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

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, userID string) {
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

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, userID string) {
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