package server

import (
	"fmt"
	_ "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/shysa/TP_proxy/app/database"
	"github.com/shysa/TP_proxy/config"
	"net/http"
	_ "net/http/pprof"
)

func New(cfg *config.Config, db *database.DB) *http.Server {
	router := gin.Default()

	router.RouterGroup = *router.Group("/api")
	{

	}

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Address, cfg.Server.Port),
		Handler: router,
	}
}