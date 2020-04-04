package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/util"

	//"github.com/kataras/iris/v12/sessions"
)

type OrderController struct {
	Ctx     iris.Context
	Session *sessions.Session
}

var orderStaticView = mvc.View{
	Name: "order.html",
	Data: iris.Map{"Title": "User Registration"},
}

type Result struct {
	Id int	`json:"id"`
	Title string	`json:"title"`

	Orders []Order
}

type Order struct {
	Id int 	`json:"id"`
	Title string 	`json:"title"`
}
func (c *OrderController) Get() mvc.Result {

	cookieName := c.Ctx.GetCookie(util.COOKEI_NAME)
	fmt.Println("cookieName=", cookieName)

	//var s []string
	//s = append(s, "ddddd")
	//s = append(s, "aaaaa")
	var s []Order
	result := new( Result)
	result.Id = 1001
	s1 := Order{
		Id:1,
		Title: "titile1",
	}
	s = append(s, s1)
	s2 := Order{
		Id:1,
		Title: "titile2",
	}
	s = append(s, s2)
	result.Orders = s
	//c.Ctx.ViewData("Result", result)
	//mvc.View{}.Data = s

	if len(c.Session.GetString(util.SessionUserName)) > 0 {
		fmt.Println("用户已经登录")
		return mvc.View{
			Name: "order.html",
			Data: iris.Map{
				"OrderCount": "10",
				"UserId": c.Session.GetString(util.SessionUserName),
				"Orders": s,
			},
		}
	} else {
		fmt.Println("用户没有登录")
		return mvc.View{
			Name: "order.html",
			Data: iris.Map{
				"Result": result,
				"OrderCount": "10002",
				"UserId": c.Session.GetString(util.SessionUserName),
			},
			//Data: map[string] interface{}{
			//	"Result": result,
			//},
		}

	}
}

func (c *OrderController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(util.SessionUserName, 0)
	return userID
}

func (c *OrderController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *OrderController) logout() {
	c.Session.Destroy()
}
