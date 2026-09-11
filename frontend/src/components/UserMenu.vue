<template>
  <div class="profile-menu-container">
    <button 
      class="profile-menu-btn" 
      @click.stop="toggleMenu"
      :aria-expanded="isMenuOpen"
      aria-label="User menu"
    >
      <OnyxIcon :icon="iconMenu" class="hamburger-icon" />
      
      <OnyxAvatar :src="profileImageUrl" :fullName="fullName" size="32px" />
    </button>

    <Transition name="dropdown">
      <div v-if="isMenuOpen" class="dropdown-menu">
        <div v-if="authStore.user" class="user-info">
          <span class="user-name">{{ fullName }}</span>
          <span class="user-email">{{ authStore.user.email }}</span>
        </div>

        <hr v-if="authStore.user" class="divider" />

        <router-link to="/profile" class="dropdown-item" @click="isMenuOpen = false">
          Edit Profile
        </router-link>

        <button class="dropdown-item logout-btn" @click="handleLogout">
          Log out
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { OnyxAvatar, OnyxIcon } from 'sit-onyx'
import { iconMenu } from '@sit-onyx/icons'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const isMenuOpen = ref(false)

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const fullName = computed(() => {
  if (!authStore.user) return 'Guest User'
  const first = authStore.user.firstName || ''
  const last = authStore.user.lastName || ''
  return `${first} ${last}`.trim() || authStore.user.userName || 'User'
})

const profileImageUrl = computed(() => {
  if (!authStore.user) return ''
  return authStore.user.profileImageUrl
}) 

const closeMenu = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.profile-menu-container')) {
    isMenuOpen.value = false
  }
}

const handleLogout = () => {
  isMenuOpen.value = false
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  window.addEventListener('click', closeMenu)
})

onUnmounted(() => {
  window.removeEventListener('click', closeMenu)
})
</script>

<style scoped>
.profile-menu-container {
  position: relative;
}

.profile-menu-btn {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.35rem 0.35rem 0.35rem 0.875rem;
  background-color: #ffffff;
  border: 1px solid #dddddd;
  border-radius: 9999px;
  cursor: pointer;
  transition: box-shadow 0.2s ease, border-color 0.2s ease;
  color: #555555;
}

.profile-menu-btn:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  border-color: #c5c5c5;
}

:deep(.hamburger-icon) {
  width: 18px;
  height: 18px;
  color: #555555;
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  width: 220px;
  background-color: #ffffff;
  border: 1px solid #ebebeb;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
  padding: 0.5rem 0;
  z-index: 100;
  display: flex;
  flex-direction: column;
}

.user-info {
  padding: 0.5rem 1.25rem;
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 700;
  font-size: 0.875rem;
  color: #111111;
}

.user-email {
  font-size: 0.75rem;
  color: #717171;
  word-break: break-all;
}

.divider {
  border: none;
  border-top: 1px solid #ebebeb;
  margin: 0.35rem 0;
}

.dropdown-item {
  padding: 0.65rem 1.25rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #222222;
  text-decoration: none;
  background: none;
  border: none;
  text-align: left;
  cursor: pointer;
  transition: background-color 0.15s ease;
  width: 100%;
}

.dropdown-item:hover {
  background-color: #f7f7f7;
}

.dropdown-item.logout-btn {
  color: #d9381e;
  font-weight: 600;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>