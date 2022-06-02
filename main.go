package main

import (
	groupie "groupie/pkg"
	"log"
	"net/http"
)

func main() {
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.Handle("/map/", http.StripPrefix("/map/", http.FileServer(http.Dir("map"))))
	groupie.Init()
	log.Println("\nThe server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
