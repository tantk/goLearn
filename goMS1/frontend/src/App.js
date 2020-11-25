import React, { Component } from 'react'
import "bootstrap/dist/css/bootstrap.min.css"
import { Link,Switch, Route} from "react-router-dom";

import AddCourse from "./components/AddCourse";
import Course from "./components/Course";
import CourseList from "./components/CourseList";

class App extends Component {
  render() {
    return (
      <div>
        <nav className="navbar navbar-expand navbar-dark bg-dark">
          <a href="/courses" className="navbar-brand">
            Go Microservice 1
          </a>
          <div className="navbar-nav mr-auto">
            <li className="nav-item">
              <Link to={"/courses"} className="nav-link">
                Courses
              </Link>
            </li>
            <li className="nav-item">
              <Link to={"/add"} className="nav-link">
                Add
              </Link>
            </li>
          </div>
        </nav>

        <div className="container mt-3">
          <Switch>
            <Route exact path={["/", "/courses"]} component={CourseList} />
            <Route exact path="/add" component={AddCourse} />
            <Route path="/courses/:id" component={Course} />
          </Switch>
        </div>
      </div>
    );
  }
}

export default App
