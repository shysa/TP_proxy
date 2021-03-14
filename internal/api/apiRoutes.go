package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shysa/TP_proxy/app/database"
)

func AddServiceRoutes(r *gin.Engine, db *database.DB)  {
	handler := NewHandler(db)

	postGroup := r.Group("/service")
	{
		postGroup.GET("/status", handler.GetServiceStatus)
		postGroup.POST("/clear", handler.ServiceClear)
	}
}

