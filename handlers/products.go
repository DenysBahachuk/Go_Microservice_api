package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/DenysBahachuk/go_product_api/data"
	"github.com/gorilla/mux"
)

// Products is a http.Handler
type Products struct {
	logger *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{logger: l}
}

// getProducts returns the products from the data store
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle GET Products")
	// fetch the products from the datastore
	productsList := data.GetProducts()
	// serialize the list to JSON
	err := productsList.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle POST Product")

	product := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&product)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable convert id to int", http.StatusBadRequest)
		return
	}
	p.logger.Println("Handle PUT Product")

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := data.Product{}

		err := product.FromJSON(r.Body)
		if err != nil {
			p.logger.Println("[ERROR] deserializing product", err)
			http.Error(w, "Error reading product", http.StatusBadRequest)
			return
		}

		//validate the product

		err = product.Validate()
		if err != nil {
			p.logger.Println("[ERROR] validating product", err)
			http.Error(w, "Error validating product", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
