package controller

import (
	"fmt"

	"GOProject/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/renderer"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rnd = renderer.New(
		renderer.Options{
			ParseGlobPattern: "html/*.html",
		},
	)
}

var rnd *renderer.Render

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func remove(slice []model.CartProduct, s int) []model.CartProduct {
	return append(slice[:s], slice[s+1:]...)
}

func (uc UserController) CalculateTotalPrice(slice []model.CartProduct) int64 {
	var totalPrice int64 = 0
	for i, _ := range slice {
		var product model.Product
		uc.session.DB("e-commerce").C("products").Find(bson.M{"_id": slice[i].ProductId}).One(&product)
		totalPrice += ((slice[i].ProductQty) * product.ProductPrice)
	}
	return totalPrice
}

func (uc UserController) CalculatePrice(cartproduct model.CartProduct, prevQty int64) int64 {
	var Price int64 = 0

	var product model.Product
	product.Id = cartproduct.ProductId
	fmt.Println(cartproduct.ProductId)
	fmt.Println(product.Id)
	if err := uc.session.DB("e-commerce").C("products").Find(bson.M{"_id": product.Id}).One(&product); err != nil {
		fmt.Println("Error we arer in")
		fmt.Println(err)
		fmt.Println("33", product)
	}
	fmt.Println("44", product)
	Price = ((cartproduct.ProductQty - prevQty) * product.ProductPrice)
	fmt.Println("Price calculated", Price)
	return Price
}

func (uc UserController) ProfileHandler(tokenstring string) (model.User, string) {
	var result model.User
	if tokenstring == "" {
		//json.NewEncoder(w).Encode("No Token Found")
		return result, "No Token Found"
	}

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte("secret"), nil
	})

	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)
		var userId = claims["userid"].(string)
		fmt.Printf("%T %s", claims["userid"], claims["userid"])
		oid := bson.ObjectIdHex(userId)
		fmt.Print(oid)
		if err := uc.session.DB("e-commerce").C("users").Find(bson.M{"_id": oid}).One(&result); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			//w.WriteHeader(404)
			//return
			return result, "User Not Found"
		}
		//json.NewEncoder(w).Encode(result)
		//return
		return result, "User Found"
	} else {
		res.Error = err.Error()
		//json.NewEncoder(w).Encode(res)
		//return
		return result, "Token invalid"
	}

}
func (uc UserController) RoleHandler(tokenstring string) (string, string) {
	var role string
	if tokenstring == "" {
		//json.NewEncoder(w).Encode("No Token Found")
		return "", "No Token Found"
	}

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte("secret"), nil
	})

	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role = claims["role"].(string)
		return role, "Token valid"
	} else {
		res.Error = err.Error()
		//json.NewEncoder(w).Encode(res)
		//return
		return "", "Token invalid"
	}

}

func (uc UserController) CheckProductsBeforeOrder(slice []model.CartProduct) (string, int) {
	flag := 1
	var productNotFound string
	for i, _ := range slice {
		product := model.Product{}
		if err := uc.session.DB("e-commerce").C("products").Find(bson.M{"_id": slice[i].ProductId}).One(&product); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			//w.WriteHeader(404)
		}
		if product.ProductQty >= slice[i].ProductQty {
		} else {
			flag = 0 //change here
			productNotFound = product.ProductName
			return productNotFound, flag
		}

	}
	return "", flag
}

func (uc UserController) UpdateProductsAfterOrder(slice []model.CartProduct) {
	for i, _ := range slice {
		product := model.Product{}
		if err := uc.session.DB("e-commerce").C("products").Find(bson.M{"_id": slice[i].ProductId}).One(&product); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			//w.WriteHeader(404)
			return
		}
		product.ProductQty -= slice[i].ProductQty

		if err := uc.session.DB("e-commerce").C("products").Update(bson.M{"_id": product.Id}, product); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
		}
	}
}
