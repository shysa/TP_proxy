package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type ProxyServer struct {
	Db *Handler
}

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

func (ps *ProxyServer) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	if req.URL.Scheme != "http" {
		http.Error(wr, "Unsupported protocol scheme ", http.StatusBadRequest)
		log.Println("Unsupported protocol scheme ", req.URL.Scheme)
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 5 * time.Second,
	}

	req.RequestURI = ""
	deleteHeaders(req.Header)

	log.Println("=====> ", req.Method, " ", req.URL, " ", req.Proto)
	//log.Println("\n", formatRequest(req))

	sr, err := httputil.DumpRequest(req, true)
	if err != nil {
		wr.WriteHeader(400)
		return
	}

	var reqId int
	if reqId, err = ps.Db.SaveRequest(req, string(sr)); err != nil {
		log.Println("Can't save request to DB: ", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP: ", err)
	}
	defer resp.Body.Close()

	deleteHeaders(resp.Header)
	copyHeader(wr.Header(), resp.Header)
	wr.WriteHeader(resp.StatusCode)

	var l bytes.Buffer
	rsp := io.MultiWriter(wr, &l)
	if _, err := io.Copy(rsp, resp.Body); err != nil {
		log.Printf("Failed to read body response: %v", err)
	}

	log.Println("<===== ", resp.StatusCode, " ", resp.Status, " ", resp.Request.URL)
	//fmt.Printf("%v\n", l.String())

	if err = ps.Db.SaveResponse(l.String(), reqId); err != nil {
		log.Println("Can't save response to DB: ", err)
		return
	}
}