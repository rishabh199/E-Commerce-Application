package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"encoding/json"

	"GOProject/model"

	"gopkg.in/mgo.v2/bson"
)

func (uc UserController) FetchPayment(w http.ResponseWriter, r *http.Request) {
	err := rnd.HTML(w, http.StatusOK, "payment2", nil)
	if err != nil {
		log.Fatal(err)
	}
}
func (uc UserController) PlaceOrder(w http.ResponseWriter, r *http.Request) {
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
	var user, checkString = uc.ProfileHandler(tokenString)

	// var user, checkString = uc.ProfileHandler(r.Header["Token"][0])
	// fmt.Printf("%T", r.Header["Token"][0])
	// fmt.Println("2121", r.Header["Token"][0])

	if checkString == "User Found" {
		cart := model.Cart{}

		if err := uc.session.DB("e-commerce").C("carts").Find(bson.M{"uid": user.Id}).One(&cart); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}
		productNotFound, flag := uc.CheckProductsBeforeOrder(cart.CartProducts)
		if flag == 0 {
			fmt.Println("Quantities insufficient:", productNotFound)
			var str2 string = "Quantities insufficient:" + productNotFound
			json.NewEncoder(w).Encode(str2)

		} else {

			cart := model.Cart{}
			if err := uc.session.DB("e-commerce").C("carts").Find(bson.M{"uid": user.Id}).One(&cart); err != nil {
				fmt.Println("Error we arer in")
				fmt.Println(err)
				w.WriteHeader(404)
				return
			}

			fmt.Println("I am here 5")

			order := model.Order{}
			order.OrderProducts = cart.CartProducts
			order.TotalPrice = cart.TotalPrice
			order.DeliveryStatus = "To be shipped!"
			order.OrderDate = time.Now()
			order.UserId = user.Id

			order.Id = bson.NewObjectId()

			//json.NewDecoder(r.Body).Decode(&cart)

			uc.session.DB("e-commerce").C("orders").Insert(order)

			fmt.Println("I am here 8")

			uc.UpdateProductsAfterOrder(order.OrderProducts)

			fmt.Println("Redirecting to Post Orders")

			http.Redirect(w, r, "/order", http.StatusMovedPermanently)
		}
	} else {

	}
}

func (uc UserController) DisplayOrders(w http.ResponseWriter, r *http.Request) {
	err := rnd.HTML(w, http.StatusOK, "orders", nil)
	if err != nil {
		log.Fatal(err)
	}
}
