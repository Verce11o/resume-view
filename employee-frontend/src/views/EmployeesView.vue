<template>
    <div class="employees">
      <div class="page-header">
        <h1>Employees</h1>
        <button @click="openAddModal" class="btn btn-primary">
          Add Employee
        </button>
      </div>
  
      <div class="card">
        <table class="table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Position</th>
              <th class="actions">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="employee in employees" :key="employee.id">
              <td>{{ employee.first_name }} {{ employee.last_name }}</td>
              <td>{{ getPositionName(employee.position_id) }}</td>
              <td class="actions">
                <button
                  @click="editEmployee(employee)"
                  class="btn-link"
                >
                  Edit
                </button>
                <button
                  @click="deleteEmployee(employee)"
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
          <h3>{{ isEditing ? 'Edit Employee' : 'Add Employee' }}</h3>
        </template>
        
        <form @submit.prevent="handleSubmit" class="form">
          <div class="form-group">
            <label class="label">First Name</label>
            <input
              v-model="form.first_name"
              type="text"
              class="input"
              required
            />
          </div>
  
          <div class="form-group">
            <label class="label">Last Name</label>
            <input
              v-model="form.last_name"
              type="text"
              class="input"
              required
            />
          </div>
  
          <div class="form-group">
            <label class="label">Position name</label>
            <input
              v-model="form.position_name"
              type="text"
              class="input"
              required
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
  
  const employees = ref([])
  const positions = ref([])
  const isModalOpen = ref(false)
  const isLoading = ref(false)
  const isEditing = ref(false)
  const currentEmployee = ref(null)
  
  const form = ref({
    first_name: '',
    last_name: '',
    position_id: ''
  })
  
  onMounted(() => {
    fetchEmployees()
    fetchPositions()
  })
  
  const fetchEmployees = async () => {
    try {
      const response = await api.get('/employee')
      employees.value = response.data.employees
    } catch (error) {
      console.error('Error fetching employees:', error)
    }
  }
  
  const fetchPositions = async () => {
    try {
      const response = await api.get('/position')
      positions.value = response.data.positions
    } catch (error) {
      console.error('Error fetching positions:', error)
    }
  }
  
  const getPositionName = (positionId) => {
    const position = positions.value.find(p => p.id === positionId)

    positions.value.forEach(element => {
        console.log(element)
    });
    return position ? position.name : 'Unknown'
  }
  
  const openAddModal = () => {
    isEditing.value = false
    currentEmployee.value = null
    form.value = {
      first_name: '',
      last_name: '',
      position_name: ''
    }
    isModalOpen.value = true
  }
  
  const editEmployee = (employee) => {
    isEditing.value = true
    currentEmployee.value = employee
    form.value = { ...employee }
    isModalOpen.value = true
  }
  
  const deleteEmployee = async (employee) => {
    if (!confirm('Are you sure you want to delete this employee?')) return
  
    try {
      await api.delete(`/employee/${employee.id}`)
      await fetchEmployees()
    } catch (error) {
      alert('Error deleting employee')
    }
  }
  
  const closeModal = () => {
    isModalOpen.value = false
    form.value = {
      first_name: '',
      last_name: '',
      position_name: ''
    }
  }
  
  const handleSubmit = async () => {
    try {
      isLoading.value = true
      if (isEditing.value) {
        await api.put(`/employee/${currentEmployee.value.id}`, form.value)
      } else {
        await api.post('/employee', form.value)
      }
      await fetchEmployees()
      closeModal()
    } catch (error) {
      alert('Error saving employee')
    } finally {
      isLoading.value = false
    }
  }
  </script>
  
  <style>
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }
  
  .page-header h1 {
    margin: 0;
    font-size: 1.5rem;
  }
  
  .table {
    width: 100%;
    border-collapse: collapse;
  }
  
  .table th,
  .table td {
    padding: 12px;
    text-align: left;
    border-bottom: 1px solid var(--color-border);
  }
  
  .table th {
    font-weight: 600;
    color: var(--color-text-secondary);
  }
  
  .table th.actions,
  .table td.actions {
    text-align: right;
  }
  
  .btn-link {
    background: none;
    border: none;
    color: var(--color-primary);
    cursor: pointer;
    padding: 4px 8px;
    font-size: 14px;
  }
  
  .btn-link:hover {
    text-decoration: underline;
  }
  
  .btn-link.delete {
    color: #dc2626;
  }
  
  .form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  
  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
  }
  </style>