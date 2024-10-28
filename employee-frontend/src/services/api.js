import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:3009',
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          localStorage.removeItem('token')
          window.location.href = '/'
          break
        case 403:
          console.error('Access forbidden')
          break
        case 404:
          console.error('Resource not found')
          break
        case 500:
          console.error('Server error')
          break
        default:
          console.error('An error occurred')
      }
    } else if (error.request) {
      console.error('Network error')
    }
    return Promise.reject(error)
  }
)

export const endpoints = {
  auth: {
    signIn: (id) => api.post('/auth/signin', { id }),
  },
  employees: {
    getAll: () => api.get('/employee'),
    getById: (id) => api.get(`/employee/${id}`),
    create: (data) => api.post('/employee', data),
    update: (id, data) => api.put(`/employee/${id}`, data),
    delete: (id) => api.delete(`/employee/${id}`)
  },
  positions: {
    getAll: () => api.get('/position'),
    getById: (id) => api.get(`/position/${id}`),
    create: (data) => api.post('/position', data),
    update: (id, data) => api.put(`/position/${id}`, data),
    delete: (id) => api.delete(`/position/${id}`)
  }
}

export default api