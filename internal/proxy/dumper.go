package proxy

import (
	"context"
	"github.com/shysa/TP_proxy/app/database"
	"net/http"
	"net/url"
	"strings"
)

type Handler struct {
	repo *database.DB
}

func NewDumper(db *database.DB) *Handler {
	return &Handler{
		repo: db,
	}
}

func (h *Handler) SaveRequest(r *http.Request, req string) (int, error) {
	u, _ := url.Parse(r.URL.String())
	p := strings.Replace(r.URL.String(), u.Scheme+"://"+r.Host, "", -1)

	var id int
	query := "insert into request(method, scheme, path, proto, host, url, request_text) values($1, $2, $3, $4, $5, $6, $7) returning id"

	if err := h.repo.QueryRow(context.Background(), query, r.Method, u.Scheme, p, r.Proto, r.Host, r.URL.String(), req).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (h *Handler) SaveResponse(r string, reqId int) error {
	if _, err := h.repo.Exec(context.Background(), "insert into response(response_text, request_id) values($1, $2)", r, reqId); err != nil {
		return err
	}
	return nil
}