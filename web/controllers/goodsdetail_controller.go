package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/models"
	//"github.com/kataras/iris/v12/sessions"
)

type GoodsDetailController struct {
	Ctx     iris.Context
	Session *sessions.Session
}

type GoodsDetailResult struct {
	Title  string `json:"title"`
	Goods *models.Goods
	Item   int `json:"item"`
}

//type Order struct {
//	Id int 	`json:"id"`
//	Title string 	`json:"title"`
//}
func (c *GoodsDetailController) GetBy( id int64) mvc.Result {

	result := new(GoodsDetailResult)
	//orderItems := new(models.OrderItems)
	result.Goods = 	models.GoodsFindById(id)
	//result.Orders = orderItems.Items
	//result.Id = 1001
	return mvc.View{
		Name: "goods-detail.html",
		Data: iris.Map{
			"Result":     result,
			"OrderCount": "10",
			//"UserId":     c.Session.GetString(util.SessionUserName),
		},
	}
}

