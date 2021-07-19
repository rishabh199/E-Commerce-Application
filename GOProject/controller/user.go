package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"encoding/json"

	"GOProject/model"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (uc UserController) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("17271")
	fmt.Println("method:", r.Method) //get request method

	if r.Method == "GET" {
		err := rnd.HTML(w, http.StatusOK, "login", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")

		var user model.User
		str1 := r.Form["username"]
		str2 := r.Form["password"]
		fmt.Println(r.Form["password"])
		user.Username = strings.Join(str1, " ")
		user.Password = strings.Join(str2, "")
		fmt.Println(user.Password)
		var foundUser model.User
		var res model.ResponseResult

		err := uc.session.DB("e-commerce").C("users").Find(bson.M{"username": user.Username}).One(&foundUser)

		if err != nil {
			res.Error = "Invalid username"
			json.NewEncoder(w).Encode(res)
			return
		}

		fmt.Println("im near type")
		fmt.Printf("%T", []byte(user.Password))
		fmt.Printf("after type")

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))

		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
		fmt.Println(string(hash))
		fmt.Println(foundUser.Password)
		if err != nil {
			fmt.Println(err)
			res.Error = "Invalid password"
			json.NewEncoder(w).Encode(res)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": foundUser.Username,
			"userid":   (foundUser.Id),
			"role":     foundUser.Role,
		})

		tokenString, err := token.SignedString([]byte("secret"))

		if err != nil {
			res.Error = "Error while generating token,Try again"
			json.NewEncoder(w).Encode(res)
			return
		}

		foundUser.Token = tokenString
		foundUser.Password = ""

		cart := model.Cart{}
		if err := uc.session.DB("e-commerce").C("carts").Find(bson.M{"uid": foundUser.Id}).One(&cart); err != nil {
			fmt.Println("Error we arer in")
			fmt.Println(err)
			w.WriteHeader(404)
			return
		}

		json.NewEncoder(w).Encode(foundUser)

	}
}

func (uc UserController) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("17272")
	fmt.Println("method:", r.Method) //get request method

	if r.Method == "GET" {
		err := rnd.HTML(w, http.StatusOK, "register", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"][0])
		fmt.Println("age:", r.Form["age"])
		fmt.Println("gender:", r.Form["gender"])
		fmt.Println("password:", r.Form["psw"])

		w.Header().Set("Content-Type", "application/json")

		var user model.User
		str1 := r.Form["username"]
		str2 := r.Form["psw"]
		user.Username = strings.Join(str1, " ")
		user.Password = strings.Join(str2, "")

		if r.Form["role"] != nil {
			user.Role = r.Form["role"][0]
		} else {
			user.Role = "customer"
		}

		var res model.ResponseResult

		var foundUser model.User
		err := uc.session.DB("e-commerce").C("users").Find(bson.M{"username": user.Username}).One(&foundUser)

		if err != nil {
			fmt.Println("I am inside post", err)
			if err.Error() == "not found" {
				hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

				if err != nil {
					res.Error = "Error While Hashing Password, Try Again"
					json.NewEncoder(w).Encode(res)
					return
				}
				user.Password = string(hash)
				user.Id = bson.NewObjectId()

				user.Name = strings.Join(r.Form["name"], " ")
				age, _ := strconv.ParseInt((r.Form["age"][0]), 10, 64)
				user.Age = age
				user.Gender = r.Form["gender"][0]
				user.PhoneNo = r.Form["phoneno"][0]
				user.Address = strings.Join(r.Form["address"], " ")
				user.EmailId = r.Form["email"][0]

				err2 := uc.session.DB("e-commerce").C("users").Insert(user)
				if err2 != nil {
					res.Error = "Error While Creating User, Try Again"
					json.NewEncoder(w).Encode(res)
					return
				}

				//new code
				cart := model.Cart{}
				cart.Id = bson.NewObjectId()
				cart.UserId = user.Id
				cart.TotalPrice = 0
				uc.session.DB("e-commerce").C("carts").Insert(cart)

				res.Result = "Registration Successful"
				json.NewEncoder(w).Encode(res)
				return
			}

			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Result = "Username already Exists!!"
		json.NewEncoder(w).Encode(res)
		return

	}
}
