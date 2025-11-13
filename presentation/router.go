package presentation

import (
	"ddd-bottomup/usecase"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	createUserUseCase *usecase.CreateUserUseCase,
	getUserUseCase *usecase.GetUserUseCase,
	updateUserUseCase *usecase.UpdateUserUseCase,
	deleteUserUseCase *usecase.DeleteUserUseCase,
) *chi.Mux {
	r := chi.NewRouter()

	userHandler := NewUserHandler(
		createUserUseCase,
		getUserUseCase,
		updateUserUseCase,
		deleteUserUseCase,
	)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", userHandler.GetUser)
			r.Put("/", userHandler.UpdateUser)
			r.Delete("/", userHandler.DeleteUser)
		})
	})

	return r
}
