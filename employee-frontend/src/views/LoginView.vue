<template>
    <div class="login-container">
      <div class="login-card">
        <h2>Sign in to your account</h2>
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label for="employee-id" class="label">Employee ID</label>
            <input
              id="employee-id"
              v-model="employeeId"
              type="text"
              class="input"
              required
              placeholder="Enter your employee ID"
            />
          </div>
  
          <button
            type="submit"
            class="btn btn-primary submit-btn"
            :disabled="isLoading"
          >
            {{ isLoading ? 'Signing in...' : 'Sign in' }}
          </button>
        </form>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '../stores/auth'
  import api from '../services/api'
  
  const router = useRouter()
  const authStore = useAuthStore()
  
  const employeeId = ref('')
  const isLoading = ref(false)
  
  const handleSubmit = async () => {
    try {
      isLoading.value = true
      const response = await api.post('/auth/signin', {
        id: employeeId.value
      })
      authStore.setToken(response.data.message)
      router.push('/employees')
    } catch (error) {
      alert('Invalid employee ID')
    } finally {
      isLoading.value = false
    }
  }
  </script>
  
  <style>
  .login-container {
    min-height: calc(100vh - 64px);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 24px;
  }
  
  .login-card {
    background: white;
    padding: 32px;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 400px;
  }
  
  .login-card h2 {
    text-align: center;
    margin-bottom: 32px;
  }
  
  .form-group {
    margin-bottom: 24px;
  }
  
  .submit-btn {
    width: 100%;
  }
  </style>