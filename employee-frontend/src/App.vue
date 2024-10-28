<template>
  <div class="app">
    <nav v-if="authStore.token" class="nav">
      <div class="container nav-container">
        <div class="nav-links">
          <router-link
            v-for="item in navigation"
            :key="item.path"
            :to="item.path"
            class="nav-link"
          >
            {{ item.name }}
          </router-link>
        </div>
        <button @click="handleLogout" class="btn btn-secondary">
          Sign Out
        </button>
      </div>
    </nav>
    <main class="main">
      <div class="container">
        <router-view></router-view>
      </div>
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from './stores/auth'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const navigation = [
  { name: 'Employees', path: '/employees' },
  { name: 'Positions', path: '/positions' }
]

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>

<style>
.nav {
  background: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.nav-container {
  height: 64px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-links {
  display: flex;
  gap: 24px;
}

.nav-link {
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
}

.nav-link:hover {
  color: var(--color-text);
}

.nav-link.router-link-active {
  color: var(--color-primary);
}

.main {
  padding: 24px 0;
}
</style>