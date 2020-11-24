package main

import (
	"fmt"
	config "gia/config"
	control "gia/controllers"
	model "gia/model"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template
var mapUsers = map[string]control.User{}
var mapSessions = map[string]string{}
var ctl = control.Ctl{
	Users:    mapUsers,
	Sessions: mapSessions,
}

func init() {
	ctl.Logging = config.CreateLogging()
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	ctl.Template = tpl
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	ctl.Users["admin"] = control.User{
		Username: "admin",
		Password: bPassword,
		First:    "admin",
		Last:     "admin"}
	ctl.Model = model.InitModel()
	ctl.Model.AddVenue(model.Venue{
		Capacity: 1235,
		Kind:     "Stadium",
		Location: "East",
		Name:     "yio chu kang stadium",
		Desc:     "fitness corner",
	})
	ctl.Model.AddVenue(model.Venue{
		Capacity: 123,
		Kind:     "Hall",
		Location: "North",
		Name:     "LT123",
		Desc:     "boring",
	})
	ctl.Model.AddVenue(model.Venue{
		Capacity: 50,
		Kind:     "Room",
		Location: "South",
		Name:     "Room151",
		Desc:     "for lessons",
	})

}

func main() {
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	router := http.NewServeMux()
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", ctl.Index)
	router.HandleFunc("/browse", ctl.Browse)
	router.HandleFunc("/book", ctl.Book)
	router.HandleFunc("/confirmBook", ctl.ConfirmBook)
	router.HandleFunc("/viewBook", ctl.ViewBook)
	router.HandleFunc("/deleteBook", ctl.DeleteBook)
	router.HandleFunc("/addVenue", ctl.AddVenue)
	router.HandleFunc("/profile", ctl.Profile)
	router.HandleFunc("/editProfile", ctl.EditProfile)
	router.HandleFunc("/signup", ctl.Signup)
	router.HandleFunc("/login", ctl.Login)
	router.HandleFunc("/logout", ctl.Logout)
	router.Handle("/favicon.ico", http.NotFoundHandler())
	server := &http.Server{
		Addr:         config.PORT,
		Handler:      config.Tracing(nextRequestID)(ctl.Logging.Infologging()(router)),
		ErrorLog:     ctl.Logging.Error,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	//http.ListenAndServe(config.PORT, nil)
	server.ListenAndServeTLS(config.CertPath, config.KeyPath)
}
