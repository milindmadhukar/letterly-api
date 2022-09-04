package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/touch-some-grass-bro/letterly-api/handlers"
)

// Function to handle routes
func (s *Server) HandleRoutes(mainRouter *chi.Mux) {
	mainRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Letterly API!"))
	})

  roomRouter := chi.NewRouter()
  roomRouter.Get("/{roomID}", handlers.GetRoom())
  roomRouter.Post("/", handlers.CreateRoom())
  roomRouter.Delete("/{roomID}", handlers.DeleteRoom())
  roomRouter.Post("/{roomID}/join", handlers.JoinRoom())
  // roomRouter.Get("/send/{roomID}/{message}", handlers.SendMessageToRoom())
  mainRouter.Mount("/room", roomRouter)

  gameRouter := chi.NewRouter()
  gameRouter.Post("/start", handlers.StartGame())
  mainRouter.Mount("/game", gameRouter)

}
