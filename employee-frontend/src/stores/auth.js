import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token'))
  const user = ref(null)

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)

    if (newToken) {
      api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
    } else {
      delete api.defaults.headers.common['Authorization']
    }
  }

  const login = async (employeeId) => {
    try {
      const response = await api.post('/auth/signin', {
        id: employeeId
      })
      setToken(response.data.id)
      return response
    } catch (error) {
      throw error
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    delete api.defaults.headers.common['Authorization']
  }

  const isAuthenticated = () => {
    return !!token.value
  }

  return {
    token,
    user,
    setToken,
    login,
    logout,
    isAuthenticated
  }
})