package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Bengissimo/Citybike/citybike"
	"github.com/Bengissimo/Citybike/server"
)

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

	s := server.New(db)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
