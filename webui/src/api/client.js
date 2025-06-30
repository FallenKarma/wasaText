import axios from 'axios'
import { setupInterceptors } from './interceptors'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
  timeout: 10000,
})

setupInterceptors(apiClient)

export default apiClient
