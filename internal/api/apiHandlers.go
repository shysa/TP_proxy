package api

import (
	_ "context"
	"github.com/gin-gonic/gin"
	"github.com/shysa/TP_proxy/app/database"
	_ "log"
	_ "net/http"
)

type Handler struct {
	repo *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{
		repo: db,
	}
}

func (h *Handler) GetServiceStatus(c *gin.Context) {
	//s := models.Status{}
	//query := "select " +
	//	"(select count(*) from forum) as f, " +
	//	"(select count(*) from post) as p, " +
	//	"(select count(*) from thread) as t, " +
	//	"(select count(*) from users) as u "
	//if err := h.repo.QueryRow(context.Background(), query).Scan(&s.Forum, &s.Post, &s.Thread, &s.User); err != nil {
	//	log.Fatal("something went wrong: ", err.Error())
	//	return
	//}
	//c.JSON(http.StatusOK, s)
}

func (h *Handler) ServiceClear(c *gin.Context) {
	//query := "truncate forum, post, thread, users, votes cascade"
	//if _, err := h.repo.Exec(context.Background(), query); err != nil {
	//	log.Fatal("something went wrong: ", err.Error())
	//	return
	//}
	//c.Status(http.StatusOK)
}