package api

import (
	"bufio"
	"context"
	_ "context"
	"github.com/gin-gonic/gin"
	"github.com/shysa/TP_proxy/app/database"
	"github.com/shysa/TP_proxy/internal/models"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"net/url"
	"strconv"
	"strings"
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
	if err := h.repo.QueryRow(context.Background(), "select * from request where id = $1", id).Scan(&r.Id, &r.Method, &r.Scheme, &r.Path, &r.Proto, &r.Host, &r.Url, &r.RequestText); err != nil {
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

func (h *Handler) ScanRequest (c *gin.Context) {
	//query := "truncate forum, post, thread, users, votes cascade"
	//if _, err := h.repo.Exec(context.Background(), query); err != nil {
	//	log.Fatal("something went wrong: ", err.Error())
	//	return
	//}
	//c.Status(http.StatusOK)
}