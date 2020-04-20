package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
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

func main() {
	//http://localhost:8080/api/search поиск
	//http://localhost:8080/api/cookie?NAME=12&V=ggg установка кук
	stopchan := make(chan os.Signal)

	router := chi.NewRouter()
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
