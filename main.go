package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

func search(str string, urls []string) (answer []string, err error) {

	group := struct {
		errgroup.Group
		sync.Mutex
		urls []string
	}{
		urls: make([]string, 0, len(urls)),
	}
	for _, value := range urls {
		u := value
		group.Go(func() error {
			resp, err := http.Get(u)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			i := strings.Index(string(body), str)

			if i >= 0 {
				group.Lock()
				group.urls = append(group.urls, u)
				group.Unlock()
			}
			return nil

		})
	}
	err = group.Wait()
	answer = group.urls
	return answer, err
}

type Resp struct {
	Search string
	Sites  []string
}

type Err struct {
	Message  error
	ErrorNum int
}
type Ansver struct {
	Sites []string
	Err
}

func getListUrl(w http.ResponseWriter, r *http.Request) {

	buffer := new(bytes.Buffer)

	buffer.ReadFrom(r.Body)

	bytesSlice := buffer.Bytes()

	var answer Ansver

	var response Resp
	json.Unmarshal(bytesSlice, &response)

	resSearch, err := search(response.Search, response.Sites)
	if err != nil {
		answer.ErrorNum = http.StatusBadGateway
		answer.Message = err
	}

	answer.Sites = resSearch
	answerJSON, err := json.Marshal(answer)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)

	w.Write(answerJSON)
}

func changeCookie(w http.ResponseWriter, r *http.Request) {

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		panic(err)
	}

	for key, value := range m {
		cookie := &http.Cookie{
			Name:  key,
			Value: value[0],
		}
		http.SetCookie(w, cookie)
	}

	cookies := r.Cookies()

	for _, value := range cookies {
		w.Write([]byte(value.Name + " = " + value.Value + " ; "))
	}

}

type BlockItems []Block
type Block struct {
	Title string
	Text  string
	Id    int
}

type Server struct {
	lg     *logrus.Logger
	Title  string
	Blocks BlockItems
}

func (s *Server) HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/static/index.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", s)
	if err != nil {
		s.lg.WithError(err).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func (s *Server) AddPost(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("./www/static/add.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", s)
	if err != nil {
		s.lg.WithError(err).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func (s *Server) NewPost(w http.ResponseWriter, r *http.Request) {
	buffer := new(bytes.Buffer)

	buffer.ReadFrom(r.Body)

	bytesSlice := buffer.Bytes()

	m, errP := url.ParseQuery(string(bytesSlice))
	if errP != nil {
		s.lg.WithError(errP).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if m.Get("Title") != "" && m.Get("Text") != "" {
		item := Block{Text: m.Get("Text"), Title: m.Get("Title"), Id: len(s.Blocks) + 1}
		s.Blocks = append(s.Blocks, item)
	}

	buffer.Reset()
	r.Body.Close()

	file, _ := os.Open("./www/static/index.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", s)
	if err != nil {
		s.lg.WithError(err).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func (s *Server) ChgangetPost(w http.ResponseWriter, r *http.Request) {
	buffer := new(bytes.Buffer)

	buffer.ReadFrom(r.Body)

	bytesSlice := buffer.Bytes()

	m, errP := url.ParseQuery(string(bytesSlice))
	if errP != nil {
		s.lg.WithError(errP).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Println(m)

	if m.Get("Title") != "" && m.Get("Text") != "" && m.Get("Id") != "" {
		for i, v := range s.Blocks {
			if id, _ := strconv.Atoi(m.Get("Id")); v.Id == id {
				s.Blocks[i].Title = m.Get("Title")
				s.Blocks[i].Text = m.Get("Text")
			}
		}
	}

	buffer.Reset()
	r.Body.Close()

	file, _ := os.Open("./www/static/index.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", s)
	if err != nil {
		s.lg.WithError(err).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) СhangeGet(w http.ResponseWriter, r *http.Request) {

	postId, _ := strconv.Atoi(chi.URLParam(r, "id"))
	file, _ := os.Open("./www/static/change.html")
	data, _ := ioutil.ReadAll(file)
	var post Block
	for _, v := range s.Blocks {
		if v.Id == postId {
			post = v
		}
	}
	templ := template.Must(template.New("page").Parse(string(data)))
	err := templ.ExecuteTemplate(w, "page", post)
	if err != nil {
		s.lg.WithError(err).Error("template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {

	stopchan := make(chan os.Signal)

	router := chi.NewRouter()

	lg := logrus.New()
	// router.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("www/static"))))

	serv := Server{
		lg:    lg,
		Title: "Блог",
		Blocks: BlockItems{
			{Text: " Text1 Text1 Text1 Text1 Text1 Text1 Text1 Text1 Text1 Text1", Title: "Title1", Id: 1},
		},
	}

	router.Route("/", func(r chi.Router) {
		r.Get("/", serv.HandleGetIndex)
		r.Get("/add", serv.AddPost)
		r.Post("/newPost", serv.NewPost)
		r.Post("/chgangetPost", serv.ChgangetPost)
		r.Get("/change/{id}", serv.СhangeGet)
	})
	router.Route("/api", func(r chi.Router) {
		r.Route("/search", func(r chi.Router) {
			r.Post("/", getListUrl)
		})
		r.Route("/cookie", func(r chi.Router) {
			r.Get("/", changeCookie)
		})

	})

	go func() {
		err := http.ListenAndServe(":8080", router)
		log.Fatal(err)
	}()

	signal.Notify(stopchan, os.Interrupt, os.Kill)
	<-stopchan
	log.Print("gracefull shutdown")
}
