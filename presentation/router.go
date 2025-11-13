package presentation

import (
	"ddd-bottomup/usecase"
	"net/http"
)

func NewRouter(
	createUserUseCase *usecase.CreateUserUseCase,
	getUserUseCase *usecase.GetUserUseCase,
	updateUserUseCase *usecase.UpdateUserUseCase,
	deleteUserUseCase *usecase.DeleteUserUseCase,
) *http.ServeMux {
	mux := http.NewServeMux()

	userHandler := NewUserHandler(
		createUserUseCase,
		getUserUseCase,
		updateUserUseCase,
		deleteUserUseCase,
	)

	mux.Handle("/users", userHandler)
	mux.Handle("/users/", userHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	return mux
}
