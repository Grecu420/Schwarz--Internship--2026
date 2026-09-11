<template>
  <header class="navbar">
    <div class="navbar-content">
      <router-link to="/" class="logo-link">
        <div class="logo-container">
          <svg width="32" height="28" viewBox="0 0 32 28" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="32" height="28" rx="6" fill="#1E40AF"/>
            <g clip-path="url(#clip0_3_13)">
              <path d="M17.75 19.2495V14.5829C17.75 14.4282 17.6885 14.2798 17.5791 14.1704C17.4697 14.061 17.3214 13.9995 17.1667 13.9995H14.8333C14.6786 13.9995 14.5303 14.061 14.4209 14.1704C14.3115 14.2798 14.25 14.4282 14.25 14.5829V19.2495M10.75 12.8331C10.75 12.6634 10.7869 12.4958 10.8584 12.3418C10.9298 12.1879 11.034 12.0514 11.1636 11.9418L15.2469 8.44181C15.4575 8.26384 15.7243 8.1662 16 8.1662C16.2757 8.1662 16.5425 8.26384 16.7531 8.44181L20.8364 11.9418C20.966 12.0514 21.0702 12.1879 21.1416 12.3418C21.2131 12.4958 21.25 12.6634 21.25 12.8331V18.0831C21.25 18.3926 21.1271 18.6893 20.9083 18.9081C20.6895 19.1269 20.3928 19.2498 20.0833 19.2498H11.9167C11.6072 19.2498 11.3105 19.1269 11.0917 18.9081C10.8729 18.6893 10.75 18.3926 10.75 18.0831V12.8331Z" stroke="white" stroke-width="2" stroke-linecap="round"/>
            </g>
            <defs>
              <clipPath id="clip0_3_13">
                <rect width="14" height="14" fill="white" transform="translate(9 7)"/>
              </clipPath>
            </defs>
          </svg>
          <span class="logo-text">StayNest</span>
        </div>
      </router-link>

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

.logo-link {
  text-decoration: none;
  flex: 1;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-text {
  font-size: 1rem; 
  font-weight: 800; 
  color: #111;
  line-height: 1;
  letter-spacing: -0.03em;
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

</style>