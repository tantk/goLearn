//Package controller :
package controller

import (
	"encoding/json"
	"fmt"
	"goms/model"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
}

//ResData : Response data
type ResData struct {
	TotalItems  int            `json:"totalItems"`
	Courses     []model.Course `json:"courses"`
	CurrentPage int            `json:"currentPage"`
	TotalPages  int            `json:"totalPages"`
}

//Allcourses : Query all data and does pagination
func Allcourses(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	resdata := ResData{}

	v := r.URL.Query()
	_, exists := v["size"]

	if exists {

		page, _ := strconv.Atoi(v["page"][0])
		size, _ := strconv.Atoi(v["size"][0])
		name, nameExists := v["name"]
		limit := fmt.Sprint(size)
		offset := fmt.Sprint(page * size)
		resdata.TotalItems, _ = strconv.Atoi(model.CountRecords("course"))
		resdata.CurrentPage = page
		resdata.TotalPages = int(math.Ceil(float64(resdata.TotalItems) / float64(size)))
		if nameExists {
			resdata.Courses = model.GetRecordsPaginationName("course", limit, offset, name[0])
		} else {
			resdata.Courses = model.GetRecordsPagination("course", limit, offset)
		}

	} else {
		data := model.GetRecords("course")
		resdata.Courses = data
	}

	json.NewEncoder(w).Encode(resdata)
}

//CourseByID : handlers GET DELETE and PUT, all requiring ID
//GET returns a single course json
//DELETE returns 204 if succeed, 404  if ID is not found
//PUT returns  204 if succeed, 404 if ID is not found
func CourseByID(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	var course model.Course
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["courseid"])
	switch r.Method {
	case "GET":
		course, _ = model.GetRecordByID("course", id)
		json.NewEncoder(w).Encode(course)

	case "DELETE":
		_, ok := model.GetRecordByID("course", id)
		if ok {
			model.DeleteCourseRecord(id)
			w.WriteHeader(http.StatusNoContent)

		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	case "PUT":
		if r.Header.Get("Content-type") == "application/json" {
			var reqC model.ReqCourse
			_, ok := model.GetRecordByID("course", id)
			if ok {
				reqBody, err := ioutil.ReadAll(r.Body)
				if err == nil {
					json.Unmarshal(reqBody, &reqC)
					model.EditCourseRecord(reqC)
					w.WriteHeader(http.StatusNoContent)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}

}

//NewCourse : checks for json obj and does a post, returns 201 if suceed
func NewCourse(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if r.Header.Get("Content-type") == "application/json" && r.Method == "POST" {
		var reqC model.ReqCourse
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqBody, &reqC)
			model.InsertCourseRecord(reqC)
			w.WriteHeader(http.StatusCreated)
		}
	}

}
