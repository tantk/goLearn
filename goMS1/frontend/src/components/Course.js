import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import CourseDataService from "../services/CourseService";

const Course = props => {
  const initialCourseState = {
    ID: "",
    Name: "",
    Desc: "",
    CertType:"",
    Enroll:0,
    Provider:"",
    Rating:0,
    Difficulty:"",
  };

  const [currentCourse, setCurrentCourse] = useState(initialCourseState);
  const [message, setMessage] = useState("");

  const getCourse = id => {
    CourseDataService.get(id)
      .then(response => {
        setCurrentCourse(
          {
            ID:response.data.ID,
            Name:response.data.Name,
            Enroll:response.data.Enrolled.Int32,
            CertType:response.data.CertType.String,
            Provider:response.data.Provider.String,
            Rating:response.data.Rating.Float64,
            Desc:response.data.Desc.String,
            Difficulty:response.data.Difficulty.String,
          });
        console.log(response.data);
      })
      .catch(e => {
        console.log(e);
      });
  };

  useEffect(() => {
    getCourse(props.match.params.id);
  }, [props.match.params.id]);

  const handleInputChange = event => {
    const { name, value } = event.target;
    setCurrentCourse({ ...currentCourse, [name]: value });
  };

  const updateCourse = () => {
    var data = {
      ID:currentCourse.ID,
      Name:currentCourse.Name,
      Enroll:parseInt(currentCourse.Enroll),
      CertType:currentCourse.CertType,
      Provider:currentCourse.Provider,
      Rating:parseFloat(currentCourse.Rating),
      Desc:currentCourse.Desc,
      Difficulty:currentCourse.Difficulty,
    };
    CourseDataService.update(currentCourse.ID, data)
      .then(response => {
        console.log(currentCourse);
        //console.log(response.data);
        setMessage("The course was updated successfully!");
      })
      .catch(e => {
        console.log(e);
      });
  };

  const deleteCourse = () => {
    CourseDataService.remove(currentCourse.ID)
      .then(response => {
        console.log(response.data);
        props.history.push("/courses");
      })
      .catch(e => {
        console.log(e);
      });
  };

  return (
    <div>
    {currentCourse ? (
      <div className="edit-form">
        <h4>Course</h4>
        <form>
          <div className="form-group">
            <label htmlFor="Name">Name</label>
            <input
              type="text"
              className="form-control"
              id="Name"
              name="Name"
              value={currentCourse.Name}
              onChange={handleInputChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="CertType">Course Type</label>
            <input
              type="text"
              className="form-control"
              id="CertType"
              name="CertType"
              value={currentCourse.CertType}
              onChange={handleInputChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="Difficulty">Difficulty</label>
            <input
              type="text"
              className="form-control"
              id="Difficulty"
              name="Difficulty"
              value={currentCourse.Difficulty}
              onChange={handleInputChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="Provider">Organizer</label>
            <input
              type="text"
              className="form-control"
              id="Provider"
              name="Provider"
              value={currentCourse.Provider}
              onChange={handleInputChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="Enroll">Students enrolled</label>
            <input
              type="number"
              className="form-control"
              id="Enroll"
              name="Enroll"
              value={currentCourse.Enroll}
              onChange={handleInputChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="Rating">Rating</label>
            <input
              type="number"
              className="form-control"
              id="Rating"
              name="Rating"
              value={currentCourse.Rating}
              onChange={handleInputChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="Desc">Description</label>
            <input
              type="text"
              className="form-control"
              id="Desc"
              name="Desc"
              value={currentCourse.Desc}
              onChange={handleInputChange}
            />
          </div>

        </form>

        <button className="badge badge-danger mr-2" onClick={deleteCourse}>
          Delete
        </button>

        <button
          type="submit"
          className="badge badge-success"
          onClick={updateCourse}
        >
          Update
        </button>
        <p>{message}</p>
      </div>
    ) : (
      <div>
        <br />
        <p>Please click on a Course...</p>
      </div>
    )}
  </div>
  );
};

export default Course;