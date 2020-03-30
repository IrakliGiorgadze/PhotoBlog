package main

import (
	"fmt"
	"net/http"

	"Study/Web_Applications/PhotoBlog/controllers"
	"Study/Web_Applications/PhotoBlog/models"

	"github.com/gorilla/mux"
)

const (
	host     = "172.30.1.147"
	port     = 5432
	user     = "dev"
	password = "dev69"
	dbname   = "photo_blog_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	must(err)
	defer us.Close()

	us.DestructiveReset()
	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()

	r.Handle("/", staticC.Home).Methods("GET")

	r.Handle("/contact", staticC.Contact).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	fmt.Println("starting server on :8080 ...")
	http.ListenAndServe(":8080", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
