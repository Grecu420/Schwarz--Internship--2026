<template>
  <header class="navbar">
    <div class="navbar-content">
      
      <div class="logo-section">
        <Logo/>
      </div>

      <nav class="nav-section">
        <router-link to="/">Home</router-link>
        <router-link to="/reservations">Reservations</router-link>
        <router-link to="/properties">Properties</router-link>
        <router-link to="/conversations">Conversations</router-link>
      </nav>

      <div class="actions-section">
        <UserMenu />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import UserMenu from '@/components/UserMenu.vue' 
import Logo from './Logo.vue';
</script>

<style scoped>
.navbar {
  background-color: #ffffff;
  border-bottom: 1px solid #ebebeb;
  display: flex;
  justify-content: center;
}

.navbar-content {
  width: 100%;
  max-width: 1440px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 3rem; 
}

.nav-section {
  display: flex;
  gap: 2.5rem;
  flex: 1;
  justify-content: center;
}

.nav-section a {
  text-decoration: none;
  color: #717171;
  font-weight: 500;
  font-size: 0.875rem;
  transition: color 0.2s ease;
}

.nav-section a:hover,
.nav-section a.router-link-active {
  color: #222222;
  font-weight: 600;
}

.actions-section {
  display: flex;
  justify-content: flex-end;
  flex: 1;
}

.logo-section {
  display: flex;
  flex: 1;
  justify-content: flex-start;
}

</style>