package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Bengissimo/Citybike/pkg/citybike"
	"github.com/Bengissimo/Citybike/pkg/server"
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
	fmt.Println("Starting server")
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
