package main

import (
	"fmt"
	"github.com/shysa/TP_proxy/app/database"
	"github.com/shysa/TP_proxy/app/server"
	"github.com/shysa/TP_proxy/config"
	"log"
	_ "net/http"
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
	apiSrv := server.NewApi(config.Cfg, dbConn, srv)

	fmt.Println("Api listening on ", apiSrv.Addr)
	go apiSrv.ListenAndServe()

	fmt.Println("Proxy listening on ", srv.Addr)
	srv.ListenAndServe()
}

