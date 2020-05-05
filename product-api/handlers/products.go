package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/wesreisz/building-microservices/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
	}

	if r.Method == http.MethodPut {
		p.updateProduct(rw, r)
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Update from Put")

	regex := regexp.MustCompile("([0-9]+)")
	g := regex.FindAllStringSubmatch(r.URL.Path, -1)
	if len(g) != 1 {
		p.l.Println("Invalid URL more than one")
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	if len(g[0]) != 2 {
		p.l.Println("Invalid URL more than one capture group", g)
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	idString := g[0][1]
	id, err := strconv.Atoi(idString)

	if err != nil {
		p.l.Println("Unable to convert to a number")
		http.Error(rw, "Unable to convert to a number", http.StatusBadRequest)
	}

	p.updateProducts(rw, id, r)
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall JSON", http.StatusBadRequest)
	}
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get Products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) updateProducts(rw http.ResponseWriter, id int, r *http.Request) {
	p.l.Println("Handle Put Products")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall JSON", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)

	err = data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
