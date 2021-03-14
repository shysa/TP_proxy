package main

import (
	"fmt"
	"github.com/shysa/TP_proxy/app/database"
	"github.com/shysa/TP_proxy/app/server"
	"github.com/shysa/TP_proxy/config"
	"github.com/shysa/TP_proxy/internal/proxy"
	"log"
	"net/http"
)

func main() {
	config.Cfg = config.Init()

	dbConn := database.NewDB(&config.Cfg.DB)
	if err := dbConn.Open(); err != nil {
		log.Fatal("connection refused: ", err)
		return
	}
	defer dbConn.Close()
	fmt.Println("connected to DB")

	srv := server.New(config.Cfg, dbConn)
	fmt.Println("listening on ", srv.Addr)
	srv.ListenAndServe()

	//handler := &proxyServer{}
	//
	//if err := http.ListenAndServe("127.0.0.1:8080", handler); err != nil {
	//	log.Fatal("Can't start proxy server: ", err)
	//}
}

