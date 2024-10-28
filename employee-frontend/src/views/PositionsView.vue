<template>
    <div class="positions">
      <div class="page-header">
        <h1>Positions</h1>
        <button @click="openAddModal" class="btn btn-primary">
          Add Position
        </button>
      </div>
  
      <div class="card">
        <table class="table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Salary</th>
              <th class="actions">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="position in positions" :key="position.id">
              <td>{{ position.name }}</td>
              <td>${{ position.salary.toLocaleString() }}</td>
              <td class="actions">
                <button
                  @click="editPosition(position)"
                  class="btn-link"
                >
                  Edit
                </button>
                <button
                  @click="deletePosition(position)"
                  class="btn-link delete"
                >
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
  
      <Modal v-model="isModalOpen">
        <template #title>
          <h3>{{ isEditing ? 'Edit Position' : 'Add Position' }}</h3>
        </template>
        
        <form @submit.prevent="handleSubmit" class="form">
          <div class="form-group">
            <label class="label">Name</label>
            <input
              v-model="form.name"
              type="text"
              class="input"
              required
            />
          </div>
  
          <div class="form-group">
            <label class="label">Salary</label>
            <input
              v-model.number="form.salary"
              type="number"
              class="input"
              required
              min="0"
            />
          </div>
  
          <div class="modal-actions">
            <button
              type="button"
              class="btn btn-secondary"
              @click="closeModal"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="btn btn-primary"
              :disabled="isLoading"
            >
              {{ isLoading ? 'Saving...' : (isEditing ? 'Update' : 'Add') }}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import api from '../services/api'
  import Modal from '../components/ui/Modal.vue'
  
  const positions = ref([])
  const isModalOpen = ref(false)
  const isLoading = ref(false)
  const isEditing = ref(false)
  const currentPosition = ref(null)
  
  const form = ref({
    name: '',
    salary: 0
  })
  
  onMounted(() => {
    fetchPositions()
  })
  
  const fetchPositions = async () => {
    try {
      const response = await api.get('/position')
      positions.value = response.data.positions
    } catch (error) {
      console.error('Error fetching positions:', error)
    }
  }
  
  const openAddModal = () => {
    isEditing.value = false
    currentPosition.value = null
    form.value = {
      name: '',
      salary: 0
    }
    isModalOpen.value = true
  }
  
  const editPosition = (position) => {
    isEditing.value = true
    currentPosition.value = position
    form.value = { ...position }
    isModalOpen.value = true
  }
  
  const deletePosition = async (position) => {
    if (!confirm('Are you sure you want to delete this position?')) return
  
    try {
      await api.delete(`/position/${position.id}`)
      await fetchPositions()
    } catch (error) {
      alert('Error deleting position')
    }
  }
  
  const closeModal = () => {
    isModalOpen.value = false
    form.value = {
      name: '',
      salary: 0
    }
  }
  
  const handleSubmit = async () => {
    try {
      isLoading.value = true
      if (isEditing.value) {
        await api.put(`/position/${currentPosition.value.id}`, form.value)
      } else {
        await api.post('/position', form.value)
      }
      await fetchPositions()
      closeModal()
    } catch (error) {
      alert('Error saving position')
    } finally {
      isLoading.value = false
    }
  }
  </script>