// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta

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

// Get Product
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
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
