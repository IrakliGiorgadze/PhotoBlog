package main

import (
	"flag"
	"fmt"
	"net/http"

	"Study/Web_Applications/PhotoBlog/controllers"
	"Study/Web_Applications/PhotoBlog/email"
	"Study/Web_Applications/PhotoBlog/middleware"
	"Study/Web_Applications/PhotoBlog/models"
	"Study/Web_Applications/PhotoBlog/rand"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {
	boolPtr := flag.Bool("prod", false, "set to true in production")
	flag.Parse()

	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)
	must(err)
	defer services.Close()

	//services.DestructiveReset()
	services.AutoMigrate()

	mgCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("jarvis.ge", "ai@sandboxe83804ac69bf46a88f406b0194c11ed8.mailgun.org"),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKey),
		)

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User, emailer)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	b, err := rand.Bytes(32)
	if err != nil {
		must(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	userMw := middleware.User{
		UserService: services.User,
	}

	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	r.HandleFunc("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")

	r.Handle("/forgot", usersC.ForgotPwView).Methods("GET")
	r.HandleFunc("/forgot", usersC.InitiateReset).Methods("POST")
	r.HandleFunc("/reset", usersC.ResetPw).Methods("GET")
	r.HandleFunc("/reset", usersC.CompleteReset).Methods("POST")

	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)

	fmt.Printf("server is starting on port:%d\n", cfg.Port)

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
