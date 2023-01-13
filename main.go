package main

import (
	"github.com/Bengissimo/Citybike/citybike"
	"html/template"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func main() {
	db, err := citybike.New("citybike.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.LoadData(); err != nil {
		log.Fatal(err)
	}
	
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)

}
