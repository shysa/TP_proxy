package api

import (
	"bufio"
	"bytes"
	"context"
	_ "context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shysa/TP_proxy/app/database"
	"github.com/shysa/TP_proxy/internal/models"
	"io"
	"io/ioutil"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	repo *database.DB
	proxy *http.Server
}

func NewHandler(db *database.DB, p *http.Server) *Handler {
	return &Handler{
		repo: db,
		proxy: p,
	}
}

func (h *Handler) GetRequests(c *gin.Context) {
	rows, err := h.repo.Query(context.Background(), "select id, request_text from request")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	}
	defer rows.Close()

	var requests []models.Request
	var r models.Request

	for rows.Next() {
		if err := rows.Scan(&r.Id, &r.RequestText); err != nil {
			log.Println("Can't read request rows from db: ", err)
		}

		requests = append(requests, r)
	}

	c.JSON(http.StatusOK, requests)
}

func (h *Handler) GetRequestById(c *gin.Context) {
	p := c.Param("id")

	id, err := strconv.Atoi(p)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r := models.Request{Id: id}

	if err := h.repo.QueryRow(context.Background(), "select request_text from request where id = $1", r.Id).Scan(&r.RequestText); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *Handler) RepeatRequest(c *gin.Context) {
	p := c.Param("id")

	id, err := strconv.Atoi(p)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r := models.Request{}
	if err := h.repo.QueryRow(context.Background(), "select id, url, request_text from request where id = $1", id).Scan(&r.Id, &r.Url, &r.RequestText); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	b := bufio.NewReader(strings.NewReader(r.RequestText))

	var req *http.Request
	req, err = http.ReadRequest(b)
	if err != nil {
		log.Println("Can't read request from db for repeat: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	req.RequestURI = ""
	req.URL, _ = url.Parse(r.Url)

	h.proxy.Handler.ServeHTTP(c.Writer, req)
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

func doCheck(c *gin.Context, req *http.Request) bool {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 5 * time.Second,
	}

	req.RequestURI = ""

	resp, err := client.Do(req)
	if err != nil {
		http.Error(c.Writer, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP: ", err)
	}

	c.Writer.Flush()

	for _, h := range hopHeaders {
		resp.Header.Del(h)
	}
	for k, headers := range resp.Header {
		for _, h := range headers {
			c.Writer.Header().Add(k, h)
		}
	}
	c.Writer.WriteHeader(resp.StatusCode)

	var l bytes.Buffer
	rsp := io.Writer(&l)
	if _, err := io.Copy(rsp, resp.Body); err != nil {
		log.Printf("Failed to read body response: %v", err)
	}
	resp.Body.Close()

	if strings.Contains(l.String(), "root:") {
		return true
	}
	return false
}

func (h *Handler) ScanRequest (c *gin.Context) {
	p := c.Param("id")

	id, err := strconv.Atoi(p)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r := models.Request{}
	if err := h.repo.QueryRow(context.Background(), "select id, method, url, request_text from request where id = $1", id).Scan(&r.Id, &r.Method, &r.Url, &r.RequestText); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	b := bufio.NewReader(strings.NewReader(r.RequestText))

	var req *http.Request
	req, err = http.ReadRequest(b)
	if err != nil {
		log.Println("Can't read request from db for repeat: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	req.RequestURI = ""
	req.URL, _ = url.Parse(r.Url)

	checkStrings := []string{`;cat /etc/passwd;`, `|cat /etc/passwd|`, "`cat /etc/passwd`"}

	switch req.Method {
	case http.MethodGet:
		body := req.URL.String()

		for _, injection := range checkStrings {
			req.URL, _ = url.Parse(body + injection)

			vulnerable := doCheck(c, req)
			fmt.Printf("[%v: %s] request №%d\n", vulnerable, injection, r.Id)
		}

	case http.MethodPost:
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("Can't read body of request for adding CI: ", err)
		}

		for _, injection := range checkStrings {
			bodyString := string(body) + injection
			req.Body = ioutil.NopCloser(strings.NewReader(bodyString))
			req.ContentLength = int64(len(bodyString))

			vulnerable := doCheck(c, req)
			fmt.Printf("[%v: %s] request №%d\n", vulnerable, injection, r.Id)
		}

	default:
		log.Println("Unsupported method for scan vulnerabilities: ", req.Method)
	}

}