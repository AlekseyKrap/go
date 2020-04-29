package server

import (
	"../my_models"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/middleware"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// Server - объект сервера
type Server struct {
	lg            *logrus.Logger
	db            *sql.DB
	rootDir       string
	templatesDir  string
	indexTemplate string
	Page          my_models.Page
	Post          my_models.Post
}

// New - создаёт новый экземпляр сервера
func New(lg *logrus.Logger, rootDir string, db *sql.DB) *Server {
	return &Server{
		lg:            lg,
		db:            db,
		rootDir:       rootDir,
		templatesDir:  "/static",
		indexTemplate: "/index.html",
		Page: my_models.Page{
			Posts: my_models.PostItemSlice{
				//{ID: "0", Text: "123", Completed: false},
				//{ID: "1", Text: "test", Completed: true},
				//{ID: "2", Text: "test 2", Completed: false},
			},
		},
	}
}

// Start - запускает сервер
func (serv *Server) Start(addr string) error {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	serv.bindRoutes(r)
	serv.lg.Debug("server is started ...")
	return http.ListenAndServe(addr, r)
}

// SendErr - отправляет ошибку пользователю и логирует её
func (serv *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	serv.lg.WithField("data", obj).WithError(err).Error("server error")
	w.WriteHeader(code)
	errModel := my_models.ErrorModel{
		Code:     code,
		Err:      err.Error(),
		Desc:     "server error",
		Internal: obj,
	}
	data, _ := json.Marshal(errModel)
	w.Write(data)
}

// SendInternalErr - отправляет 500 ошибку
func (serv *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	serv.SendErr(w, err, http.StatusInternalServerError, obj)
}
