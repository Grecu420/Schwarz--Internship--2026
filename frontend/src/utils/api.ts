import axios, { type AxiosInstance, type InternalAxiosRequestConfig } from 'axios';
import { useAuthStore } from '@/stores/auth'; 
import router from '@/router'; // Importăm router-ul Vue

const api: AxiosInstance = axios.create();

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const authStore = useAuthStore(); 
  const token = authStore.token; 

  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status;
    console.log('API Error Status:', status);

    const authStore = useAuthStore();

    if (status === 404) {
      console.log('Redirecting to 404...');
      router.push('/404');
      return Promise.reject(error);
    }

    if (status === 401) {
      console.log('Redirecting to Login...');
      authStore.logout();
      router.push('/login');
      return Promise.reject(error);
    }

    return Promise.reject(error);
  }
);

export default api;