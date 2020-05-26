package controllers

import (
	"net/http"

	"github.com/deltegui/phoenix"
)

type JSONPresenter struct {
	renderer phoenix.JSONRenderer
}

func NewJSONPresenter(w http.ResponseWriter) JSONPresenter {
	return JSONPresenter{
		renderer: phoenix.JSONRenderer{w},
	}
}

func (b JSONPresenter) Present(data interface{}) {
	b.renderer.Render(data)
}

func (b JSONPresenter) PresentError(err error) {
	b.renderer.RenderError(err)
}
