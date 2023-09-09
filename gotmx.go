package main

import (
	"fmt"
	"log"
	"net/http"
)

var PORT string = "8080"

func handleLanding(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>gotmx landing page</h1>")
}

func main() {
	fmt.Println("Running gotmx")

	http.HandleFunc("/", handleLanding)

	fmt.Printf("Listening to: http://localhost:%s\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}
