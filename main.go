package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// Default Handler/
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Default route\n")
}

// Main entrypoint. Listens for requests on the port below.
func main() {
	router := httprouter.New()
	router.GET("/", Index)
	log.Fatal(http.ListenAndServe(":8009", router))
}
