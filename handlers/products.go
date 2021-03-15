package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jonathanbs9/go-building-microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// Serve Http
/*func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// GET
	if r.Method == http.MethodGet {
		p.GetProducts(rw, r)
		return
	}

	// POST
	if r.Method == http.MethodPost {
		p.AddProduct(rw, r)
		return
	}

	// PUT
	if r.Method == http.MethodPut {
		p.l.Println("Handle PUT Product")
		// expect the id in the URI
		p.l.Println("PUT => ", r.URL.Path)
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI => more than one ID")
			http.Error(rw, "Ivalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("Invalid URI => More than one capture group")
			http.Error(rw, "Ivalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI => unable to convert to number")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.l.Println("Got ID", id)
		p.updateProducts(id, rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}*/

// Get Product
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Hanle GET Products")
	lp := data.GetProducts()
	rw.Header().Add("Content-Type", "application/json")
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// Add Product
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(&prod)
}

// Update Product
func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID | ", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product | ", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Products) MiddleWareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Unable to Unmarshal json", http.StatusBadRequest)
			return
		}

		// Validate the product
		err = prod.Validate()
		if err != nil {

			p.l.Println("[ERROR] validating product")
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		// Add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
