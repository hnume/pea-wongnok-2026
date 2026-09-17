import axios from "axios";
import { getSession } from "next-auth/react";

const instance = axios.create({
  baseURL: "http://localhost:8080/api/v1",
  timeout: 5000,
  headers: { "Content-Type": "application/json" },
});

instance.interceptors.request.use(async config => {
  const session = await getSession()

  config.headers.Authorization = `Bearer ${session?.accessToken}`
  return config
})

export { instance as axios };
