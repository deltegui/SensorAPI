package controllers

import (
	"github.com/deltegui/locomotive"
	"net/http"
)

type ErrorController struct{}

func NewErrorController() ErrorController {
	return ErrorController{}
}

func (ErrorController ErrorController) NotFound(w http.ResponseWriter, req *http.Request) {
	presenter := locomotive.JSONPresenter{w}
	presenter.Present(struct{Code string}{Code: "404"})
}

func (ErrorController ErrorController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Get, Handler: ErrorController.NotFound, Endpoint: "404"},
	}
}
