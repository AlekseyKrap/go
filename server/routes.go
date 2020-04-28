package server

import (
	"github.com/go-chi/chi"
)

func (serv *Server) bindRoutes(r *chi.Mux) {

	r.Route("/", func(r chi.Router) {

		r.Get("/", serv.getTemplateHandler)
		r.Get("/add", serv.getAddPost)
		r.Post("/newPost", serv.postNewPost)
		r.Get("/newPost", serv.getTemplateHandler)
		r.Get("/change/{id}", serv.changeGet)
		r.Post("/changePost", serv.postChangePost)

	})
}
