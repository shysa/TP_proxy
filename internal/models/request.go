package models

type Request struct {
	Id          int    `json:"id"`
	Method      string `json:"method,omitempty"`
	Scheme      string `json:"scheme,omitempty"`
	Path        string `json:"path,omitempty"`
	Proto       string `json:"proto,omitempty"`
	Host        string `json:"host,omitempty"`
	Url         string `json:"url,omitempty"`
	RequestText string `json:"request_text,omitempty"`
}
