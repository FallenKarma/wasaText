import apiClient from '../client'

const usersApi = {
  fetchUsers() {
    return apiClient.get('/users')
  },

  fetchCurrentUser() {
    return apiClient.get('/users/me')
  },

  updateUsername(name) {
    return apiClient.put('/users/me/username', { name })
  },

  uploadPhoto(photoFile) {
    const formData = new FormData()
    formData.append('photo', photoFile)

    return apiClient.put('/users/me/photo', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
  },
}

export default usersApi
