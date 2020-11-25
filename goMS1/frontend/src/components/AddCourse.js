import React, { useState } from "react";
import CourseDataService from "../services/CourseService";

const AddCourse = () => {
  const initialCourseState = {
    id: null,
    name: "",
    description: "",
    enroll: "",
    certType:"",
    provider:"",
    rating:"",
    difficulty:"",
  };
  const [course, setCourse] = useState(initialCourseState);
  const [submitted, setSubmitted] = useState(false);

  const handleInputChange = event => {
    const { name, value } = event.target;
    setCourse({ ...course, [name]: value });
  };

  const saveCourse = () => {
    var data = {

      Name:course.name,
      Enroll:parseInt(course.enroll),
      CertType:course.certType,
      Provider:course.provider,
      Rating:parseFloat(course.rating),
      Desc:course.description,
      Difficulty:course.difficulty,
    };
    console.log("before post")
    CourseDataService.create(data)
      .then(response => {
        setSubmitted(true);
        console.log(response.data);
      })
      .catch(e => {
        console.log(e);
      });
    console.log("after post")
  };

  const newCourse = () => {
    setCourse(initialCourseState);
    setSubmitted(false);
  };

  return (
    <div className="submit-form">
    {submitted ? (
      <div>
        <h4>You submitted successfully!</h4>
        <button className="btn btn-success" onClick={newCourse}>
          Add
        </button>
      </div>
    ) : (
      <div>
        <div className="form-group">
          <label htmlFor="name">Name</label>
          <input
            type="text"
            className="form-control"
            id="name"
            required
            value={course.name}
            onChange={handleInputChange}
            name="name"
          />
        </div>
        <div className="form-group">
          <label htmlFor="certType">Course Type</label>
          <input
            type="text"
            className="form-control"
            id="certType"
            required
            value={course.certType}
            onChange={handleInputChange}
            name="certType"
          />
        </div>
        <div className="form-group">
          <label htmlFor="difficulty">Difficulty</label>
          <input
            type="text"
            className="form-control"
            id="difficulty"
            required
            value={course.difficulty}
            onChange={handleInputChange}
            name="difficulty"
          />
        </div>
        <div className="form-group">
          <label htmlFor="provider">Organizer</label>
          <input
            type="text"
            className="form-control"
            id="provider"
            required
            value={course.provider}
            onChange={handleInputChange}
            name="provider"
          />
        </div>
        <div className="form-group">
          <label htmlFor="enroll">Students enrolled</label>
          <input
            type="number"
            className="form-control"
            id="enroll"
            required
            value={course.enroll}
            onChange={handleInputChange}
            name="enroll"
          />
        </div>
        <div className="form-group">
          <label htmlFor="number">Course rating</label>
          <input
            type="number"
            step="0.01"
            className="form-control"
            id="rating"
            required
            value={course.rating}
            onChange={handleInputChange}
            name="rating"
          />
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <input
            type="text"
            className="form-control"
            id="description"
            required
            value={course.description}
            onChange={handleInputChange}
            name="description"
          />
        </div>

        <button onClick={saveCourse} className="btn btn-success">
          Submit
        </button>
      </div>
    )}
  </div>
  );
};

export default AddCourse;