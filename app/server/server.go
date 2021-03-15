package server

import (
	"fmt"
	_ "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/shysa/TP_proxy/app/database"
	"github.com/shysa/TP_proxy/config"
	"github.com/shysa/TP_proxy/internal/proxy"
	"net/http"
	_ "net/http/pprof"
)

func New(cfg *config.Config, db *database.DB) *http.Server {
	handler := &proxy.ProxyServer{
		Db: proxy.NewDumper(db),
	}

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Address, cfg.Server.Port),
		Handler: handler,
	}
}

func NewApi(cfg *config.Config, db *database.DB) *http.Server {
	router := gin.Default()

	// добавить апишные роуты

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Api.Address, cfg.Api.Port),
		Handler: router,
	}
}