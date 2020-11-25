import http from "../http-common";

const getAll = (params) => {
  return http.get("/courses",{params});
};

const get = id => {
  return http.get(`/courses/${id}`);
};

const create = data => {
  return http.post("/courses/new", data);
};

const update = (id, data) => {
  return http.put(`/courses/${id}`, data);
};

const remove = id => {
  return http.delete(`/courses/${id}`);
};


export default {
  getAll,
  get,
  create,
  update,
  remove,
};