import axios from "axios";

const instance = axios.create({
  baseURL: "http://localhost:8080/api/v1",
  timeout: 5000,
  headers: { "Content-Type": "application/json" },
});

export { instance as axios };
