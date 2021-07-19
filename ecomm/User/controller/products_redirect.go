package controller

import (

	//"C:/Users/Dell/Desktop/GOProject/model"

	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (uc UserController) GetAllProducts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("in getallproducts")
	http.Redirect(w, r, "http://localhost:8082/products", http.StatusMovedPermanently)

}

func (uc UserController) GetProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("in getproducts")
	http.Redirect(w, r, "http://localhost:8082/product/"+p.ByName("id"), http.StatusMovedPermanently)
}

func (uc UserController) CreateProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("in CreateProduct")
	http.Redirect(w, r, "http://localhost:8082/products", http.StatusMovedPermanently)

}
