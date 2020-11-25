package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"goms/conf"
	ctrl "goms/controller"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/courses", ctrl.Allcourses)
	router.HandleFunc("/api/v1/courses/new", ctrl.NewCourse)
	router.HandleFunc("/api/v1/courses/{courseid}", ctrl.CourseByID).Methods(
		"GET", "PUT", "DELETE",
	)
	//need to use cors to allow using REST api in the same system
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowCredentials: true,
	})
	handler := c.Handler((router))
	conf := conf.GetConfig()
	fmt.Println("Listening at port ", conf.RESTport)
	log.Fatal(http.ListenAndServe(":"+conf.RESTport, handler))

}
