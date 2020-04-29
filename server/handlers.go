package server

import (
	"../my_models"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/gorilla/sessions"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func (serv *Server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {

	templateName := serv.indexTemplate
	file, err := os.Open(path.Join(serv.rootDir, serv.templatesDir, templateName))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("page").Parse(string(data))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	posts, err := my_models.GetAllTaskItems(serv.db)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	serv.Page.Posts = posts

	if err := templ.Execute(w, serv.Page); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}
func (serv *Server) getAddPost(w http.ResponseWriter, r *http.Request) {
	templateName := serv.indexTemplate
	file, err := os.Open(path.Join(serv.rootDir, serv.templatesDir, templateName))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("pageADD").Parse(string(data))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	cookie := &http.Cookie{
		Name:  "token",
		Value: token,
	}
	http.SetCookie(w, cookie)
	serv.Page.Token = token

	if err := templ.Execute(w, serv.Page); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}
func (serv *Server) postNewPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token, _ := r.Cookie("token")

	if token.Value != r.Form["Token"][0] {
		serv.getTemplateHandler(w, r)
		return
	}

	if r.Form["Title"][0] == "" || r.Form["Text"][0] == "" {
		err := errors.New("Not Title or Text")
		serv.SendInternalErr(w, err)
		return
	}

	post := my_models.PostItem{Text: r.Form["Text"][0], Title: r.Form["Title"][0]}

	if err := post.Insert(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:  "token",
		Value: "",
	}
	http.SetCookie(w, cookie)

	serv.getTemplateHandler(w, r)

}
func (serv *Server) changeGet(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "id")

	post, err := my_models.GetPost(serv.db, postId)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	templateName := serv.indexTemplate
	file, err := os.Open(path.Join(serv.rootDir, serv.templatesDir, templateName))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("pageChange").Parse(string(data))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))

	token := fmt.Sprintf("%x", h.Sum(nil))
	cookie := &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
	serv.Post.Token = token
	serv.Post.Posts = post

	http.SetCookie(w, cookie)

	if err := templ.Execute(w, serv.Post); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}
func (serv *Server) postChangePost(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	token, err := r.Cookie("token")
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	if token.Value != r.Form["Token"][0] {
		serv.getTemplateHandler(w, r)
		return
	}

	if r.Form["Title"][0] == "" || r.Form["Text"][0] == "" || r.Form["ID"][0] == "" {
		err := errors.New("Not Title or Text")
		serv.SendInternalErr(w, err)
		return
	}

	post := my_models.PostItem{Text: r.Form["Text"][0], Title: r.Form["Title"][0], ID: r.Form["ID"][0]}

	if err := post.Update(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:  "token",
		Value: "",
	}
	http.SetCookie(w, cookie)

	serv.getTemplateHandler(w, r)

}
