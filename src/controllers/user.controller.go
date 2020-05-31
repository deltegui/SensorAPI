package controllers

import (
	"log"
	"net/http"
	"sensorapi/src/domain"

	"github.com/deltegui/phoenix"
	"github.com/gorilla/csrf"
)

func UserIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		(phoenix.HTMLRenderer{w}).RenderData("userindex.html", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(req),
		})
	}
}

func ProcessUserLogin(execUserCase domain.LoginUserCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		user := req.Form.Get("name")
		pass := req.Form.Get("password")
		log.Printf("User: %s, password: %s\n", user, pass)
		response, err := execUserCase(domain.LoginUserRequest{
			UserName:     user,
			UserPassword: pass,
		})
		if err != nil {
			http.Redirect(w, req, "/user/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, req, "/user/panel", http.StatusSeeOther)
	}
}

func registerUserRoutes() {
	phoenix.MapGroup("/user", func(m phoenix.Mapper) {
		m.MapAll([]phoenix.Mapping{
			{Method: phoenix.Get, Builder: UserIndex, Endpoint: "/login"},
			{Method: phoenix.Post, Builder: ProcessUserLogin, Endpoint: "/login"},
		}, phoenix.NewCSRFMiddleware())
	})
}
