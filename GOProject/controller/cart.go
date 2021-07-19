package controller

import (
	"GOProject/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

func (uc UserController) DisplayCart(w http.ResponseWriter, r *http.Request) {
	// funcMap := template.FuncMap{
	// 	// The name "inc" is what the function will be called in the template text.
	// 	"inc": func(i int) int {
	// 		return i + 1
	// 	},
	// }
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

	//var user, checkString = uc.ProfileHandler(r.Header["Token"][0])
	// fmt.Printf("%T", r.Header["Token"][0])
	// fmt.Println("2121", r.Header["Token"][0])

	if checkString == "User Found" {
		type DisplayCartProduct struct {
			ProductName  string
			ProductImg   string
			ProductPrice int64
			ProductId    bson.ObjectId
			ProductQty   int64
		}
		var cartproducts []DisplayCartProduct
		cart := model.Cart{}

		if err := uc.session.DB("e-commerce").C("carts").Find(bson.M{"uid": user.Id}).One(&cart); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}
		oldcartproducts := cart.CartProducts

		for i, _ := range cart.CartProducts {
			var cartproduct DisplayCartProduct
			cartproduct.ProductName = cart.CartProducts[i].ProductName
			cartproduct.ProductQty = cart.CartProducts[i].ProductQty
			cartproduct.ProductId = cart.CartProducts[i].ProductId

			var product model.Product
			product.Id = cart.CartProducts[i].ProductId
			fmt.Println(cartproduct.ProductId)
			fmt.Println(product.Id)
			if err := uc.session.DB("e-commerce").C("products").Find(bson.M{"_id": product.Id}).One(&product); err != nil {
				fmt.Println("Error we arer in")
				fmt.Println(err)
				fmt.Println("33", product)
			}

			cartproduct.ProductImg = product.ProductImg
			cartproduct.ProductPrice = product.ProductPrice * cartproduct.ProductQty

			cartproducts = append(cartproducts, cartproduct)
		}
		fmt.Println(oldcartproducts)
		fmt.Println(cartproducts)
		tpls := []string{"html/cart/layout.tmpl", "html/cart/index.tmpl", "html/cart/partial.tmpl"}

		err := rnd.Template(w, http.StatusOK, tpls, cartproducts)
		if err != nil {
			log.Fatal(err) //respond with error page or message
		}
	} else {
		fmt.Println("232", checkString)
	}

}

func (uc UserController) UpdateCart(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	token, err := r.Cookie("token")
	if err != nil {
		fmt.Println("cookie was not found")
		return
	}

	tokenString := token.String()[6:]
	var user, checkString = uc.ProfileHandler(tokenString)

	//var user, checkString = uc.ProfileHandler(r.Header["Token"][0])

	if checkString == "User Found" {

		var cartproducts model.CartProduct
		cartproducts.ProductName = r.Form["product_name"][0]
		// fmt.Println("121", r.Form["product_name"][0])
		// fmt.Println("7212", r.Form["product_id"][0][13:37])
		// fmt.Printf("%T", r.Form["product_id"][0])

		cartproducts.ProductId = bson.ObjectIdHex(r.Form["product_id"][0][13:37])
		//cartproducts.ProductId = bson.ObjectIdHex(r.Form["product_id"][0]) // works for postman
		qty, _ := strconv.ParseInt(r.Form["product_quantity"][0], 10, 64)
		cartproducts.ProductQty = qty
		fmt.Print(cartproducts)

		cart := model.Cart{}
		if err := uc.session.DB("e-commerce").C("carts").Find(bson.M{"uid": user.Id}).One(&cart); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}

		oldcartproducts := cart.CartProducts
		fmt.Println("1", oldcartproducts)
		flag := 0

		for i, _ := range oldcartproducts {
			fmt.Println("oldcartproducts[i].ProductName", oldcartproducts[i].ProductName)
			if oldcartproducts[i].ProductId == cartproducts.ProductId {
				var prevQty = oldcartproducts[i].ProductQty
				oldcartproducts[i].ProductQty = cartproducts.ProductQty
				flag = 1
				fmt.Println("2", oldcartproducts)

				Price := uc.CalculatePrice(oldcartproducts[i], prevQty)
				var totalPrice = cart.TotalPrice + Price
				fmt.Println("Total price:", totalPrice)

				if err := uc.session.DB("e-commerce").C("carts").Update(bson.M{"_id": cart.Id}, bson.M{"$set": bson.M{"cartproducts": oldcartproducts, "totalprice": totalPrice}}); err != nil {
					fmt.Println("Error we arer in")
					fmt.Println(err)
					w.WriteHeader(404)
					return
				}
			}
		}
		if flag == 0 {
			oldcartproducts = append(oldcartproducts, cartproducts)
			fmt.Println("2", oldcartproducts)

			fmt.Println("222", cart.Id)

			Price := uc.CalculatePrice(cartproducts, 0)
			var totalPrice = cart.TotalPrice + Price
			fmt.Println("Total price:", totalPrice)

			if err := uc.session.DB("e-commerce").C("carts").Update(bson.M{"_id": cart.Id}, bson.M{"$set": bson.M{"cartproducts": oldcartproducts, "totalprice": totalPrice}}); err != nil {
				fmt.Println("Error we arer in")
				fmt.Println(err)
				w.WriteHeader(404)
				return
			}
		}
		http.Redirect(w, r, "/cart", http.StatusMovedPermanently)
	} else {
		fmt.Println(checkString)
	}

}

func (uc UserController) DeleteCart(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")
	//var res model.ResponseResult
	//var userid string = "" //get user id

	var user, checkString = uc.ProfileHandler(r.Header["Token"][0])

	if checkString == "User Found" {
		fmt.Println("2")
		body, err := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
		fmt.Println("4")
		fmt.Println(err, "6")
		var cartproducts model.CartProduct
		fmt.Println("221", r.Form["product_name"])
		cartproducts.ProductName = r.Form["product_name"][0]
		cartproducts.ProductId = bson.ObjectIdHex(r.Form["product_id"][0])
		qty, _ := strconv.ParseInt(r.Form["product_quantity"][0], 10, 64)
		cartproducts.ProductQty = qty
		fmt.Print(cartproducts)

		cart := model.Cart{}
		if err := uc.session.DB("e-commerce").C("carts").Find(bson.M{"uid": user.Id}).One(&cart); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}
		oldcartproducts := cart.CartProducts
		fmt.Println("1", oldcartproducts)

		for i, _ := range oldcartproducts {
			fmt.Println("oldcartproducts[i].ProductName", oldcartproducts[i].ProductName)
			if oldcartproducts[i].ProductId == cartproducts.ProductId {
				oldcartproducts = remove(oldcartproducts, i)
			}
		}

		fmt.Println("2", oldcartproducts)

		fmt.Println("222", cart.Id)

		Price := uc.CalculatePrice(cartproducts, 0)
		var totalPrice = cart.TotalPrice - Price
		fmt.Println("Total price:", totalPrice)

		if err := uc.session.DB("e-commerce").C("carts").Update(bson.M{"_id": cart.Id}, bson.M{"$set": bson.M{"cartproducts": oldcartproducts, "totalprice": totalPrice}}); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}
	} else {
		fmt.Println(checkString)
	}
}

func (uc UserController) EmptyCart(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	token, err := r.Cookie("token")
	if err != nil {
		fmt.Println("cookie was not found")
		return
	}

	tokenString := token.String()[6:]
	var user, checkString = uc.ProfileHandler(tokenString)
	fmt.Println("****************************************************************")
	//var user, checkString = uc.ProfileHandler(r.Header["Token"][0])

	if checkString == "User Found" {

		fmt.Println("****************************************************************")
		var oldcartproducts []model.CartProduct
		totalPrice := 0
		fmt.Println("1", oldcartproducts)

		if err := uc.session.DB("e-commerce").C("carts").Update(bson.M{"uid": user.Id}, bson.M{"$set": bson.M{"cartproducts": oldcartproducts, "totalprice": totalPrice}}); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}
		var cart model.Cart
		if err := uc.session.DB("e-commerce").C("carts").Find(bson.M{"uid": user.Id}).One(&cart); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}
		fmt.Println("22222222222", oldcartproducts, cart)

		//http.Redirect(w, r, "/products", http.StatusMovedPermanently)
		json.NewEncoder(w).Encode(cart)
	} else {
		fmt.Println(checkString)
	}

}

func (uc UserController) EmptyCart2(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Inside empty cart222")
}
