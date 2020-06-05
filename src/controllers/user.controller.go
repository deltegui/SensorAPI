package controllers

import (
	"log"
	"net/http"
	"sensorapi/src/domain"

	"github.com/deltegui/phoenix"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

func UserIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		phoenix.NewHTMLRenderer(w).RenderData("userindex.html", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(req),
		})
	}
}

func ProcessUserLogin(execUserCase domain.LoginUserCase, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		user := req.Form.Get("name")
		pass := req.Form.Get("password")
		log.Printf("User: %s, password: %s\n", user, pass)
		_, err := execUserCase(domain.LoginUserRequest{
			UserName:     user,
			UserPassword: pass,
		})
		if err != nil {
			http.Redirect(w, req, "/user/login", http.StatusSeeOther)
			return
		}
		session, err := store.Get(req, "session")
		if err != nil {
			log.Fatalln(err)
		}
		session.Values["username"] = user
		session.Save(req, w)
		http.Redirect(w, req, "/user/panel", http.StatusSeeOther)
	}
}

func registerUserRoutes(app phoenix.App) {
	app.MapGroup("/user", func(m phoenix.Mapper) {
		csrf := phoenix.NewCSRFMiddleware()
		m.Get("/login", UserIndex, csrf)
		m.Post("/login", ProcessUserLogin, csrf)
	})
}
