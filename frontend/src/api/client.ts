import { userManager } from "@/config";
import axios from "axios";

export const publicClient = axios.create({
  baseURL: "/api",
  timeout: 10000,
});

export const privateClient = axios.create({
  baseURL: "/api",
  timeout: 10000,
});

privateClient.interceptors.request.use(async (config) => {
  const user = await userManager.getUser();
  if (user?.access_token) {
    config.headers.Authorization = `Bearer ${user.access_token}`;
  }
  return config;
});
