package handlers

import (
	"log"
	"net/http"
	"product-api/product-api/data"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}
	if r.Method == http.MethodPut {

		regex := regexp.MustCompile(`/([0-9]+)`)
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)
		p.l.Println("G", g, len(g), g[0][1])
		if len(g) != 1 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, _ := strconv.Atoi(idString)
		p.l.Println("GOt Id", id)
		p.updateProduct(id, w, r)
		return

	}
	w.WriteHeader(http.StatusMethodNotAllowed)

}
func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle put", id)
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Bad Request in parsing", http.StatusBadRequest)
		return
	}
	p.l.Println("product ", prod)
	errorp := data.UpdateProduct(id, prod)
	if errorp != nil {
		p.l.Println("eroor data", errorp)
		http.Error(w, "Erro", http.StatusMethodNotAllowed)
		return
	}
	if errorp == data.ErrProductNotFound {
		http.Error(w, "Erro", http.StatusBadRequest)
		return
	}

}
func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Printf("Handle POST")
	prod := &data.Product{}
	p.l.Printf("data %#v", prod)
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	p.l.Printf("Prod:  %#v", prod)
	data.AddProduct(prod)
}
func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to parse data", http.StatusInternalServerError)
	}
}
