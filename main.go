package main

import (
	"fmt"
	"github.com/speps/go-hashids"
	"github.com/ssdb/gossdb/ssdb"
	"log"
	"net/http"
	"strings"
)

var (
	ip   = "127.0.0.1"
	port = 8888
)

func main() {
	http.HandleFunc("/a/", shorterHandler)
	http.HandleFunc("/s/", redirectHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//shorter is short url to 8-len token
func shorter(url string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = url
	hd.MinLength = 8
	h, err := hashids.NewWithData(hd)
	if err != nil {
		err = fmt.Errorf("problem with newwithdata: %s", err)
		return "", err
	}
	e, err := h.Encode([]int{len(url)})
	if err != nil {
		err = fmt.Errorf("problem with encode: %s", err)
		return "", err
	}
	return e, nil
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
	fmt.Fprintf(w, err.Error())
}

func shorterHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	_, ok := q["url"]
	if !ok || q["url"][0] == "" {
		err := fmt.Errorf("don't find url")
		handleError(w, err)
		return
	}
	url := q["url"][0]
	fmt.Println(url)
	short, err := shorter(url)
	if err != nil {
		handleError(w, err)
		return
	}

	db, err := ssdb.Connect(ip, port)
	if err != nil {
		handleError(w, err)
		return
	}
	defer db.Close()

	db.Set(short, url)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", short)
}

func checkToken(token string) error {
	if len(token) != 8 {
		return fmt.Errorf("invalid token: len = %d", len(token))
	}
	for _, c := range token {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			continue
		}
		return fmt.Errorf("invalid token")
	}
	return nil
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	if len(s) < 3 {
		err := fmt.Errorf("don't have token")
		handleError(w, err)
		return
	}
	if len(s) > 3 {
		err := fmt.Errorf("too much parameters")
		handleError(w, err)
		return
	}
	token := s[2]
	if err := checkToken(token); err != nil {
		err = fmt.Errorf("%s", err)
		handleError(w, err)
		return
	}
	db, err := ssdb.Connect(ip, port)
	if err != nil {
		err = fmt.Errorf("can't connect to data base %s", err)
		handleError(w, err)
		return
	}
	defer db.Close()

	url, err := db.Get(s[2])
	if url == nil {
		log.Printf("Can't find tokes %s", token)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Don't find token")
		return
	}
	if err != nil {
		err = fmt.Errorf("problem with get from db by token %s: %s", token, err)
		handleError(w, err)
		return
	}
	http.Redirect(w, r, url.(string), http.StatusFound)
}
