package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shysa/TP_proxy/app/database"
	"net/http"
)

func AddServiceRoutes(r *gin.Engine, db *database.DB, p *http.Server)  {
	handler := NewHandler(db, p)

	r.GET("/requests", handler.GetRequests)
	r.GET("/requests/:id", handler.GetRequestById)
	r.GET("/repeat/:id", handler.RepeatRequest)
	r.GET("/scan/:id", handler.ScanRequest)
}

