package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	flag.Parse()
	handler := &proxyServer{}

	if err := http.ListenAndServe("127.0.0.1:8080", handler); err != nil {
		log.Fatal("Can't start proxy server: ", err)
	}
}

//func dumpTo(filename string) *os.File {
//	dumpTo := os.Stdout
//	if len(filename) > 0 {
//		file, err := os.Create(filename)
//		if err != nil {
//			log.Printf("Fail to open file %s, fallback to stdout", filename)
//		} else {
//			dumpTo = file
//		}
//	}
//	return dumpTo
//}


