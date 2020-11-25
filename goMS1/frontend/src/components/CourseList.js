import React, { useState, useEffect } from "react";
import Pagination from "@material-ui/lab/Pagination";
import { Link } from "react-router-dom";

import CourseDataService from "../services/CourseService";

const CourseList = () => {
  const [courses, setCourses] = useState([]);
  const [currentCourse, setCurrentCourse] = useState(null);
  const [currentIndex, setCurrentIndex] = useState(-1);
  const [searchName, setSearchName] = useState("");

  const [page, setPage] = useState(1);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);

  const pageSizes = [10, 20, 50, 100];


  const onChangeSearchName = e => {
    const searchName = e.target.value;
    setSearchName(searchName);
  };

  const getRequestParams = (searchName, page, pageSize) => {
    let params = {};

    if (searchName) {
      params["name"] = searchName;
    }

    if (page) {
      params["page"] = page - 1;
    }

    if (pageSize) {
      params["size"] = pageSize;
    }

    return params;
  };

  const retrieveCourses = () => {
    const params = getRequestParams(searchName, page, pageSize);
    CourseDataService.getAll(params)
      .then(response => {
        const { courses, totalPages } = response.data;
        setCourses(courses);
        setCount(totalPages);
        console.log(response.data);
      })
      .catch(e => {
        console.log(e);
      });
  };
  useEffect(retrieveCourses, [page, pageSize]);



  const refreshList = () => {
    retrieveCourses();
    setCurrentCourse(null);
    setCurrentIndex(-1);
  };

  const setActiveCourse = (course, index) => {
    setCurrentCourse(course);
    setCurrentIndex(index);
  };


  // const findByName = () => {
  //   CourseDataService.findByName(searchName)
  //     .then(response => {
  //       setCourses(response.data);
  //       console.log(response.data);
  //     })
  //     .catch(e => {
  //       console.log(e);
  //     });
  // };

  const handlePageChange = (event, value) => {
    setPage(value);
  };
  const handlePageSizeChange = (event) => {
    setPageSize(event.target.value);
    setPage(1);
  };

  return (
    <div className="list row">
      <div className="col-md-8">
        <div className="input-group mb-3">
          <input
            type="text"
            className="form-control"
            placeholder="Search by name"
            value={searchName}
            onChange={onChangeSearchName}
          />
          <div className="input-group-append">
            <button
              className="btn btn-outline-secondary"
              type="button"
              onClick={retrieveCourses}
            >
              Search
            </button>
          </div>
        </div>
      </div>
      <div className="col-md-6">
        <h4>Courses List</h4>

        <div className="mt-3">
          {"Items per Page: "}
          <select onChange={handlePageSizeChange} value={pageSize}>
            {pageSizes.map((size) => (
              <option key={size} value={size}>
                {size}
              </option>
            ))}
          </select>

          <Pagination
            className="my-3"
            count={count}
            page={page}
            siblingCount={1}
            boundaryCount={1}
            variant="outlined"
            shape="rounded"
            onChange={handlePageChange}
          />
        </div>

        <ul className="list-group">
          {courses &&
            courses.map((course, index) => (
              <li
                className={
                  "list-group-item " + (index === currentIndex ? "active" : "")
                }
                onClick={() => setActiveCourse(course, index)}
                key={index}
              >
                {course.Name}
              </li>
            ))}
        </ul>


      </div>
      <div className="col-md-6">
        {currentCourse ? (
          <div>
            <h4>Course</h4>
            <div>
              <label>
                <strong>Name:</strong>
              </label>{" "}
              {currentCourse.Name}
            </div>
            <div>
              <label>
                <strong>Course Type:</strong>
              </label>{" "}
              {currentCourse.CertType.String}
            </div>
            <div>
              <label>
                <strong>Difficulty:</strong>
              </label>{" "}
              {currentCourse.Difficulty.String}
            </div>
            <div>
              <label>
                <strong>Organizer:</strong>
              </label>{" "}
              {currentCourse.Provider.String}
            </div>
            <div>
              <label>
                <strong>Students enrolled:</strong>
              </label>{" "}
              {currentCourse.Enrolled.Int32}
            </div>
            <div>
              <label>
                <strong>Course Rating :</strong>
              </label>{" "}
              {currentCourse.Rating.Float64}
            </div>
            <div>
              <label>
                <strong>Course ID :</strong>
              </label>{" "}
              {currentCourse.ID}
            </div>
            <div>
              <label>
                <strong>Description :</strong>
              </label>{" "}
              {currentCourse.Desc.Valid ? currentCourse.Desc.String : "Nil"}
            </div>

            

            <Link
              to={"/courses/" + currentCourse.ID}
              className="badge badge-warning"
            >
              Edit
            </Link>
          </div>
        ) : (
          <div>
            <br />
            <p>Please click on a Course...</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default CourseList;