//Package model :
package model

import (
	"database/sql"
	"fmt"
	"goms/conf"

	"github.com/jmoiron/sqlx"
)

//Course :
type Course struct {
	ID         int             `json:"ID" db:"id"`
	Name       string          `json:"Name" db:"course_name"`
	Provider   sql.NullString  `json:"Provider" db:"course_provider"`
	CertType   sql.NullString  `json:"CertType" db:"course_cert_type"`
	Rating     sql.NullFloat64 `json:"Rating" db:"course_rating"`
	Difficulty sql.NullString  `json:"Difficulty" db:"course_difficulty"`
	Enrolled   sql.NullInt32   `json:"Enrolled" db:"course_students_enrolled"`
	Desc       sql.NullString  `json:"Desc" db:"course_desc"`
	Require    sql.NullString  `json:"Require" db:"course_req"`
	Topics     sql.NullString  `json:"Topics" db:"course_topics"`
}

//ReqCourse : Course data format from frontend
type ReqCourse struct {
	ID         int     `json:"ID"`
	Name       string  `json:"Name"`
	Enroll     int     `json:"Enroll"`
	CertType   string  `json:"CertType"`
	Provider   string  `json:"Provider"`
	Rating     float32 `json:"Rating"`
	Desc       string  `json:"Desc"`
	Difficulty string  `json:"Difficulty"`
}

// sqlx package
func dbConn() (db *(sqlx.DB)) {
	conf := conf.GetConfig()
	dataSourceName := conf.DbUser + ":" + conf.DbPW + "@tcp(" + conf.DbHost + ":" + conf.DbPort + ")/" + conf.DbName
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

//sql driver
func dbConnRaw() (db *(sql.DB)) {
	conf := conf.GetConfig()
	dataSourceName := conf.DbUser + ":" + conf.DbPW + "@tcp(" + conf.DbHost + ":" + conf.DbPort + ")/" + conf.DbName
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

//GetRecords :
func GetRecords(table string) []Course {
	db := dbConn()

	courses := []Course{}
	err := db.Select(&courses, "SELECT * FROM "+table)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	return courses
}

//CountRecords :
func CountRecords(table string) string {
	db := dbConn()

	var count []uint8
	err := db.Get(&count, "select count(*) as count FROM "+table)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	return string(count)
}

//GetRecordsPagination :
func GetRecordsPagination(table string, limit string, offset string) []Course {
	db := dbConn()
	courses := []Course{}
	err := db.Select(&courses, "SELECT * FROM "+table+" LIMIT "+limit+" OFFSET "+offset)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	return courses
}

//GetRecordsPaginationName :
func GetRecordsPaginationName(table string, limit string, offset string, name string) []Course {
	db := dbConn()
	courses := []Course{}
	err := db.Select(&courses, "SELECT * FROM "+table+" WHERE course_name LIKE '%"+name+"%' LIMIT "+limit+" OFFSET "+offset)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	return courses
}

//GetRecordByID :
func GetRecordByID(table string, id int) (Course, bool) {
	course := Course{}
	db := dbConn()
	query := fmt.Sprintf(`SELECT * FROM `+table+` where id=%d`, id)
	err := db.Get(&course, query)
	if err != nil {
		fmt.Println(err)
		return course, false
	}
	defer db.Close()
	return course, true
}

//InsertCourseRecord :
func InsertCourseRecord(rc ReqCourse) {
	db := dbConn()
	query := fmt.Sprintf(
		`INSERT INTO course (
			course_name,course_provider,course_cert_type,course_rating,
			course_difficulty,course_students_enrolled,course_desc
			) VALUES ('%s', '%s', '%s', %f,'%s',%d,'%s')`,
		rc.Name, rc.Provider, rc.CertType, rc.Rating,
		rc.Difficulty, rc.Enroll, rc.Desc)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}

//EditCourseRecord :
func EditCourseRecord(rc ReqCourse) {
	db := dbConn()
	query := fmt.Sprintf(
		`UPDATE course SET course_name="%s",
		course_provider="%s",
		course_cert_type="%s",
		course_rating=%f,
		course_difficulty="%s",
		course_students_enrolled=%d,
		course_desc="%s"
		WHERE id=%d`,
		rc.Name, rc.Provider, rc.CertType, rc.Rating,
		rc.Difficulty, rc.Enroll, rc.Desc, rc.ID)
	fmt.Println(query)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}

//DeleteCourseRecord :
func DeleteCourseRecord(ID int) {
	db := dbConn()
	query := fmt.Sprintf(
		"DELETE FROM course WHERE ID='%d'", ID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}
