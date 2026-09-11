<template>
  <main class="split-layout">
    <div class="left-panel">
      <div class="form-wrapper">
        <div class="header">
          
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

          <h1>Create an account</h1>
          <p>Join our community of hosts and travelers around the globe.</p>
        </div>

        <RegisterForm />

        <p class="login-prompt">
          Already have an account? <router-link to="/login" class="login-link">Log in</router-link>
        </p>
      </div>
    </div>

    <div class="right-panel"></div>
  </main>
</template>

<script setup lang="ts">
import RegisterForm from '../components/RegisterForm.vue'
</script>

<style scoped>
:global(body), :global(html) {
  margin: 0;
  padding: 0;
  height: 100%;
}

.split-layout {
  display: flex;
  min-height: 100vh;
  width: 100%;
  font-family: var(--onyx-font-family-base, system-ui, sans-serif);
}

.left-panel {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
}

.form-wrapper {
  width: 100%;
  max-width: 420px;
}

.header {
  margin-bottom: 2rem;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 1.5rem;
}

.logo-text {
  font-size: 1.25rem;
  font-weight: 800;
  color: #111;
}

.header h1 {
  font-size: 2rem;
  margin: 0 0 0.5rem 0;
  color: #111;
}

.header p {
  color: #666;
  margin: 0;
}

.login-prompt {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 0.9rem;
}

.login-prompt a {
  color: #1a41b5;
  font-weight: bold;
  text-decoration: none;
}

.right-panel {
  flex: 1;
  background: url('https://images.unsplash.com/photo-1499793983690-e29da59ef1c2?q=80&w=2070') center/cover no-repeat;
}

@media (max-width: 900px) {
  .right-panel {
    display: none;
  }
}

.login-link {
  color: #2c5282; 
  text-decoration: none;
  font-weight: 600;
  display: inline-block;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.login-link:hover {
  color: #1a41b5; 
  transform: translateY(-1px) scale(1.03); 
  filter: brightness(1.1);
}

.login-link:active {
  color: #0f2671;
  transform: translateY(0) scale(0.98);
}
</style>