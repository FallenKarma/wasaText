// src/api/endpoints/auth.js
import apiClient from '../client'

export const authApi = {
  login(userData) {
    return apiClient.post('/session', userData)
  },
}
