package controllers

import (
	"log"
	"net/http"
	"time"

	"Study/Web_Applications/PhotoBlog/context"
	"Study/Web_Applications/PhotoBlog/email"
	"Study/Web_Applications/PhotoBlog/models"
	"Study/Web_Applications/PhotoBlog/rand"
	"Study/Web_Applications/PhotoBlog/views"
)

type Users struct {
	NewView      *views.View
	LoginView    *views.View
	ForgotPwView *views.View
	ResetPwView  *views.View
	us           models.UserService
	emailer      *email.Client
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

func NewUsers(us models.UserService, emailer *email.Client) *Users {
	return &Users{
		NewView:      views.NewView("bootstrap", "users/new"),
		LoginView:    views.NewView("bootstrap", "users/login"),
		ForgotPwView: views.NewView("bootstrap", "users/forgot_pw"),
		ResetPwView:  views.NewView("bootstrap", "users/reset_pw"),
		us:           us,
		emailer:      emailer,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	parseURLParams(r, &form)
	u.NewView.Render(w, r, form)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	go u.emailer.Welcome(user.Name, user.Email)
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Welcome to jarvis.ge, please behave do not break anything",
	}
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}
	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}
	err = u.emailer.ResetPw(form.Email, token)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Please check E-Mail",
	})
}

func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.SetAlert(err)
	}
	u.ResetPwView.Render(w, r, vd)
}

func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}
	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}
	u.signIn(w, user)
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Password reset was been successful, Enjoy!",
	})
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}
