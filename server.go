package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type proxyServer struct {}

// ----------------------------------------------
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Proxy-Connection",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func deleteHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

// ----------------------------------------------

func (ps *proxyServer) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Println("=====> ", req.RemoteAddr, " ", req.Method, " ", req.URL)

	if req.URL.Scheme != "http" {
		msg := "unsupported protocol scheme " + req.URL.Scheme
		http.Error(wr, msg, http.StatusBadRequest)
		log.Println(msg)
		return
	}

	client := &http.Client{}

	req.RequestURI = ""
	deleteHeaders(req.Header)

	//proxyReq, err := http.NewRequest(req.Method, req.RequestURI, req.Body)
	//copyHeader(proxyReq.Header, req.Header)

	log.Println("\n", formatRequest(req))

	//_, err = httputil.DumpRequest(proxyReq, true)
	//if err != nil {
	//	wr.WriteHeader(400)
	//	return
	//}
	// S A V E ----------------

	resp, err := client.Do(req)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP: ", err)
	}
	defer resp.Body.Close()

	log.Println("<===== ", resp.StatusCode, " ", resp.Status)

	deleteHeaders(resp.Header)
	copyHeader(wr.Header(), resp.Header)
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)

	//rsp, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("%v\n", string(rsp))
}

func formatRequest(r *http.Request) string {
	var request []string

	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)

	request = append(request, fmt.Sprintf("Host: %v", r.Host))

	for name, headers := range r.Header {
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n %v", r.Form.Encode())
	}

	return strings.Join(request, "\n")
}

//func formatResponse(wr http.ResponseWriter, req *http.Request) string {}
