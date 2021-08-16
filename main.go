package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

type Data struct {
	db map[string]string
}

func main() {

	data := &Data{db: make(map[string]string)}

	// SERVER
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()

	r.HandleFunc("/", data.homepage)
	r.HandleFunc("/to/{key}", data.redirect)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(port, nil))
}

func short() string {

	var letters string = "_QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890"

	randLetters := make([]byte, 10)
	for i := range randLetters {
		randLetters[i] = letters[rand.Intn(len(letters))]
	}
	return string(randLetters)
}

func (data *Data) redirect(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	k := vars["key"]
	http.Redirect(w, r, data.db[k], http.StatusSeeOther)
	w.WriteHeader(http.StatusOK)

}

func (data *Data) homepage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, "Ошибка шаблона: %s/n", err.Error())
	}
	var link string
	if r.Method == "POST" {
		link = r.FormValue("link")
		if !ValidUrl(link) {
			link = "Неверный формат ссылки"
		} else {
			sh := short()
			shortlink := "http://localhost:8080/to/" + sh
			data.db[sh] = link
			link = shortlink
			for key, value := range data.db {
				fmt.Printf("%s === %s \n", key, value)
			}
			fmt.Println("____________________________________")

		}
	}
	temp.Execute(w, link)
}

func ValidUrl(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}
