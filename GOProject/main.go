package main

import (
	"net/http"

	"GOProject/controller"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

func main() {

	// r := httprouter.New()
	r := mux.NewRouter()
	uc := controller.NewUserController(getSession())

	r.HandleFunc("/login", uc.Login) //only 2 for handlefunc
	r.HandleFunc("/register", uc.Register)

	r.HandleFunc("/products", uc.Products)
	r.HandleFunc("/product", uc.AddProduct)

	r.HandleFunc("/cart", uc.DisplayCart).Methods("GET")
	r.HandleFunc("/cart", uc.UpdateCart).Methods("POST")
	r.HandleFunc("/cart", uc.DeleteCart).Methods("DELETE")
	r.HandleFunc("/emptycart", uc.EmptyCart).Methods("GET")

	r.HandleFunc("/prepayment", uc.FetchPayment).Methods("GET")
	r.HandleFunc("/order", uc.PlaceOrder).Methods("POST")
	r.HandleFunc("/order", uc.DisplayOrders).Methods("GET")
	//http.HandleFunc("/getUser", uc.ProfileHandler)
	// r.GET("/users", uc.GetAllUsers)
	// r.GET("/user/:id", uc.GetUser)
	// r.POST("/user", uc.CreateUser)
	// r.DELETE("/user/:id", uc.DeleteUser)

	//r.GET("/products", uc.GetAllProducts)
	//r.GET("/product/:id", uc.GetProduct)

	// r.POST("/product", uc.CreateProduct)
	// r.DELETE("/product/:id", uc.DeleteProduct)

	// r.GET("/carts", uc.GetAllCarts)
	// r.POST("/carts", uc.CreateCart)
	// //r.DELETE("/user/:id/cart", uc.DeleteCart)

	// r.GET("/user/:id/cart", uc.GetCartUser)
	// r.POST("/user/:id/cart", uc.AddToCart)
	// r.DELETE("/user/:id/cart", uc.DeleteItemInCart)

	// r.GET("/user/:id/payment", uc.GetPayment)
	// r.POST("/user/:id/payment", uc.PostPayment)

	// r.GET("/user/:id/order", uc.PlaceOrder)

	http.ListenAndServe("localhost:8080", r)
}

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}
	return s
}
