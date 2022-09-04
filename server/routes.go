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

  // TODO: Change to respective method types
  roomRouter.Get("/{roomID}", handlers.GetRoom())
  roomRouter.Post("/", handlers.CreateRoom())
  roomRouter.Delete("/delete/{roomID}", handlers.DeleteRoom())

  // roomRouter.Get("/send/{roomID}/{message}", handlers.SendMessageToRoom())

  mainRouter.Mount("/room", roomRouter)

  // mainRouter.Get()

}
