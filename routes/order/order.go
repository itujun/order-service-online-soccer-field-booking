package routes

import (
	"github.com/gin-gonic/gin"
	"order-service/clients"
	"order-service/constants"
	"order-service/controllers/http"
	"order-service/middlewares"
)

type OrderRoute struct {
	controllers.IControllerRegistry
	client clients.IClientRegistry
	group  *gin.RouterGroup
}

type IOrderRoute interface {
	Run()
}

func NewOrderRoute(
	group *gin.RouterGroup,
	controller controllers.IControllerRegistry,
	client clients.IClientRegistry,
) IOrderRoute {
	return &OrderRoute{
		IControllerRegistry: controller,
		client:              client,
		group:               group,
	}
}

func (o *OrderRoute) Run() {
	group := o.group.Group("/order")
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, o.client), o.GetOrder().GetAllWithPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, o.client), o.GetOrder().GetByUUID)
	group.GET("/user", middlewares.CheckRole([]string{
		constants.Customer,
	}, o.client), o.GetOrder().GetOrderByUserID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Customer,
	}, o.client), o.GetOrder().Create)
}
