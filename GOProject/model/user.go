package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id bson.ObjectId `json:"id" bson:"_id"`

	Name   string `json:"name" bson:"name"`
	Gender string `json:"gender" bson:"gender"`
	Age    int64  `json:"age" bson:"age"`

	PhoneNo string `json:"phoneno" bson:"phoneno"`
	EmailId string `json:"emailid" bson:"emailid"`
	Address string `json:"address" bson:"address"`

	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Token    string `json:"token" bson:"token"`
	Role     string `json:"role" bson:"role"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type Product struct {
	Id bson.ObjectId `json:"id" bson:"_id"`

	ProductImg   string `json:"pimg" bson:"pimg"`
	ProductName  string `json:"pname" bson:"pname"`
	ProductQty   int64  `json:"pqty" bson:"pqty"`     //int to int 64
	ProductPrice int64  `json:"pprice" bson:"pprice"` //float to int64
}

type Cart struct {
	Id           bson.ObjectId `json:"id" bson:"_id"`
	CartProducts []CartProduct `json:"cartproducts" bson:"cartproducts"`
	//Username     string        `json:"uname" bson:"uname"` //Add user Id here
	UserId     bson.ObjectId `json:"uid" bson:"uid"`
	TotalPrice int64         `json:"totalprice" bson:"totalprice"` //float to int64
}

type CartProduct struct {
	ProductName string `json:"pname" bson:"pname"`
	//ProductPrice float64 `json:"pprice" bson:"pprice"` //Add product ID instead of product name and price here
	ProductId  bson.ObjectId `json:"pid" bson:"pid"`
	ProductQty int64         `json:"pqty" bson:"pqty"` //int  to int64
}

type Order struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	OrderProducts []CartProduct `json:"orderproducts" bson:"orderproducts"`
	// Username       string        `json:"uname" bson:"uname"`
	UserId         bson.ObjectId `json:"uid" bson:"uid"`
	TotalPrice     int64         //float to int64
	OrderDate      time.Time
	DeliveryStatus string
	Address        string
}

type Payment struct {
	Id bson.ObjectId `json:"id" bson:"_id"`
	// Username   string        `json:"uname" bson:"uname"`
	UserId     bson.ObjectId `json:"uid" bson:"uid"`
	TotalPrice float64
}
