package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Bengissimo/Citybike/citybike"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func main() {
	downloadData := flag.Bool("download", false, "download data")
	flag.Parse()

	db, err := citybike.New("citybike.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if *downloadData {
		if err = db.LoadData(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Download completed successfully")
		return
	}

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)

}
