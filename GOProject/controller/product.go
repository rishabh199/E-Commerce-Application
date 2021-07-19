package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"encoding/json"

	"GOProject/model"

	"gopkg.in/mgo.v2/bson"
)

func (uc UserController) Products(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		products := []model.Product{}

		if err := uc.session.DB("e-commerce").C("products").Find(nil).All(&products); err != nil {
			w.WriteHeader(404)
			return
		}

		tpls := []string{"html/products/layout.tmpl", "html/products/index.tmpl", "html/products/partial.tmpl"}

		err := rnd.Template(w, http.StatusOK, tpls, products)
		if err != nil {
			log.Fatal(err) //respond with error page or message
		}
	} else {
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		var res model.ResponseResult
		var product model.Product
		product.ProductName = r.Form["productname"][0]

		err := uc.session.DB("e-commerce").C("products").Find(bson.M{"pname": product.ProductName}).One(&product)
		if err != nil {
			fmt.Println("I am inside post", err)
			if err.Error() == "not found" {

				product.ProductImg = r.Form["productimage"][0]
				qty, _ := strconv.ParseInt(r.Form["productquantity"][0], 10, 64)
				product.ProductQty = qty
				price, _ := strconv.ParseInt(r.Form["productprice"][0], 10, 64)
				product.ProductPrice = price
				product.Id = bson.NewObjectId()

				err2 := uc.session.DB("e-commerce").C("products").Insert(product)
				if err2 != nil {
					res.Error = "Error While Creating Product, Try Again"
					fmt.Println(err2)
					json.NewEncoder(w).Encode(res)
					return
				}
				res.Result = "Product Added"
				json.NewEncoder(w).Encode(res)
				return
			}
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		res.Result = "Product already Exists!!"
		json.NewEncoder(w).Encode(res)
		return
	}

}

func (uc UserController) AddProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")

		token, err := r.Cookie("token")
		if err != nil {
			fmt.Println("cookie was not found")
			return
		}

		tokenString := token.String()[6:]
		fmt.Printf("%T", tokenString)
		fmt.Println("23222", tokenString)
		var role, checkString = uc.RoleHandler(tokenString)

		if role == "admin" {
			err := rnd.HTML(w, http.StatusOK, "productAdd", nil)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println(role, checkString)
			json.NewEncoder(w).Encode("Admin access required!!...:((")
		}
	}
}
