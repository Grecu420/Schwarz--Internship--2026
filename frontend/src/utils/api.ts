import axios, { type AxiosInstance, type InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/stores/auth'; 

const api: AxiosInstance = axios.create();

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const authStore = useAuthStore(); 
  const token = authStore.token; 

  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;