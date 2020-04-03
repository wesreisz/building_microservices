package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wesreisz/building-microservices/product-api/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	http.ListenAndServe(":9090", nil)
}
