package auth

import (
	"akselander/sandbox/pkg/shared/api"
	"database/sql"
	"strings"

	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func LoginWallMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ac := c.(*Context)
		if ac.User.id != nil {
			return next(ac)
		}

		return redirectToAuth(c)
	}
}

func App(db *sql.DB, e *echo.Echo) {
	r := &userRepository{db}
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				c.Logger().Errorf("Could not retrieve session, reason: %v", err)
			}

			var id int
			var u *User

			id, ok := sess.Values["id"].(int)

			if ok {
				u, err = r.findById(id)
			}

			if err != nil {
				c.Logger().Errorf("Error while decoding session: %v", err)
			}

			if u == nil {
				u = &User{id: nil, hash: nil, Username: "", Password: ""}
			}

			ac := &Context{
				Context:        c,
				userRepository: r,
				User:           u,
			}

			return next(ac)
		}
	})
	g := e.Group("/auth")
	g.GET("/", redirect)
	g.GET("/logout/", logout)
	g.GET("/user/", getUser)

	s := g.Group("/signup")
	s.Use(userRedirect)
	s.GET("/", getSignup)
	s.GET("/submit/", getSubmitSignup)
	s.GET("/username/", getUsernameTaken)
	s.POST("/submit/", postSignupUser)

	l := g.Group("/login")
	l.Use(userRedirect)
	l.GET("/", getLogin)
	l.GET("/submit/", getSubmitLogin)
	l.POST("/submit/", postAuthenticateUser)
}

func userRedirect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ac := c.(*Context)
		if ac.User.id == nil {
			return next(ac)
		}

		return redirectAfterAuth(c)
	}
}

func redirect(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/auth/signup/")
}

func getSignup(c echo.Context) error {
	status := http.StatusOK
	ac := c.(*Context)
	u := ac.User
	u.Username = c.QueryParam("username")
	page := Page{User: *u, Mode: "signup", From: c.QueryParam("from")}

	replaceUsernameURI(c)
	return c.Render(status, api.EnsureProperTemplate(c, "auth-username", "auth-page"), page)
}

func getLogin(c echo.Context) error {
	status := http.StatusOK
	ac := c.(*Context)
	u := ac.User
	u.Username = c.QueryParam("username")
	page := Page{User: *u, Mode: "login", From: c.QueryParam("from")}

	replaceUsernameURI(c)
	return c.Render(status, api.EnsureProperTemplate(c, "auth-username", "auth-page"), page)
}

func getSubmitSignup(c echo.Context) error {
	status := http.StatusOK
	ac := c.(*Context)
	template := api.EnsureProperTemplate(c, "auth-password", "auth-password-page")
	u := ac.User
	u.Username = c.QueryParam("username")
	page := Page{User: *u, Mode: "signup", From: c.FormValue("from"), Errors: map[string]string{}}

	taken, err := u.nameTaken(ac.userRepository)

	if err != nil {
		return err
	}

	if taken != "" {
		status = http.StatusSeeOther
		template = api.EnsureProperTemplate(c, "auth-username", "auth-page")
		page.setError("username", taken)
		replaceUsernameURI(c)
	}

	return c.Render(status, template, page)
}

func getSubmitLogin(c echo.Context) error {
	status := http.StatusOK
	ac := c.(*Context)
	template := api.EnsureProperTemplate(c, "auth-password", "auth-password-page")
	u := ac.User
	u.Username = c.QueryParam("username")

	if u.Username == "" {
		template = api.EnsureProperTemplate(c, "auth-user", "auth-page")
	}
	page := Page{User: *u, Mode: "login", From: c.QueryParam("from")}

	return c.Render(status, template, page)
}

func getUsernameTaken(c echo.Context) error {
	ac := c.(*Context)
	status := http.StatusOK

	u := ac.User
	u.Username = c.QueryParam("username")
	taken, err := u.nameTaken(ac.userRepository)
	if err != nil {
		return err
	}

	if api.IsHtmxRequest(c) {
		return c.String(status, taken)
	}

	return c.Render(status, "auth-page", Page{
		User:   *u,
		Mode:   "login",
		Errors: map[string]string{"username": taken},
		From:   c.QueryParam("from"),
	})

}

func postAuthenticateUser(c echo.Context) error {
	c.Logger().Infof("Authentication attempt for user: %v", getUsernameTaken)
	template := api.EnsureProperTemplate(c, "auth-password", "auth-page-password")
	ac := c.(*Context)
	u := ac.User
	u.Username = c.FormValue("username")
	u.Password = c.FormValue("password")
	page := Page{User: *u, Mode: "login", From: c.FormValue("from"), Errors: map[string]string{}}
	id, result := u.authenticate(ac.userRepository)

	if result != "" {
		page.setError("password", result)
		status := http.StatusSeeOther
		return c.Render(status, template, page)
	}

	sess, err := session.Get("session", c)

	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	sess.Values["id"] = id
	sess.Save(c.Request(), c.Response())

	return redirectAfterAuth(c)
}

func postSignupUser(c echo.Context) error {
	template := api.EnsureProperTemplate(c, "auth-password", "auth-page-password")
	status := http.StatusOK
	ac := c.(*Context)
	u := ac.User
	u.Username = c.FormValue("username")
	u.Password = c.FormValue("password")
	page := Page{User: *u, Mode: "signup", From: c.FormValue("from"), Errors: map[string]string{}}

	taken, err := u.nameTaken(ac.userRepository)

	if err != nil {
		return err
	}

	if taken != "" {
		page.setError("password", taken)
		return c.Render(status, template, page)
	}

	id, passwordInvalid, err := u.save(ac.userRepository)

	if err != nil {
		return err
	}

	if passwordInvalid != "" {
		page.setError("password", passwordInvalid)
		return c.Render(status, template, page)
	}

	status = http.StatusFound
	sess, err := session.Get("session", c)

	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	sess.Values["id"] = id
	sess.Save(c.Request(), c.Response())

	return redirectAfterAuth(c)
}

func getUser(c echo.Context) error {
	ac := c.(*Context)
	u := ac.User

	if u.id == nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login/")
	}

	return c.Render(http.StatusSeeOther, api.EnsureProperTemplate(c, "auth-user", "auth-user-page"), Page{Mode: "user", User: *u})
}

func logout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["id"] = nil
	sess.Save(c.Request(), c.Response())

	if api.IsHtmxRequest(c) {
		return c.Render(http.StatusFound, "auth-username", Page{Mode: "login", User: User{}})
	}

	return c.Redirect(http.StatusFound, "/auth/login")
}

func redirectToAuth(c echo.Context) error {
	if api.IsHtmxRequest(c) {
		builder := strings.Builder{}
		builder.WriteString("/auth/signup/?from=")
		builder.WriteString(c.Request().URL.RequestURI())
		c.Response().Header().Set("HX-Replace-Url", builder.String())
		return c.Render(http.StatusSeeOther, "auth-page", Page{Mode: "signup", User: User{}, From: c.Request().URL.RequestURI()})
	}

	return c.Redirect(http.StatusSeeOther, "/auth/")
}

func redirectAfterAuth(c echo.Context) error {
	from := c.FormValue("from")

	if from != "" {

		if api.IsHtmxRequest(c) {
			c.Response().Header().Set("HX-Redirect", from)
			return c.NoContent(http.StatusSeeOther)
		}

		c.Response().Header().Set("HX-Retarget", "body")
		return c.Redirect(http.StatusSeeOther, from)
	}

	ac := c.(*Context)
	u := ac.User

	if api.IsHtmxRequest(c) {
		return c.Render(http.StatusSeeOther, "auth-user", Page{Mode: "user", User: *u})
	}

	return c.Redirect(http.StatusSeeOther, "/auth/user/")
}

func replaceUsernameURI(c echo.Context) {
	c.Response().Header().Set("HX-Replace-Url", strings.Replace(c.Request().URL.RequestURI(), "submit/", "", 1))
}

type Context struct {
	echo.Context
	Mode           string
	User           *User
	userRepository *userRepository
}

type Page struct {
	User   User
	Mode   string
	Errors map[string]string
	From   string
}

func (p *Page) setError(name string, error string) {
	p.Errors[name] = error
}
