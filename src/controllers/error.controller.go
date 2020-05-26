package controllers

import (
	"net/http"

	"github.com/deltegui/phoenix"
)

func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		presenter := phoenix.JSONRenderer{w}
		presenter.Render(struct {
			Code string `json:"code"`
		}{Code: "404"})
	}
}
