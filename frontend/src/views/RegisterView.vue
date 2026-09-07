<template>
  <main class="split-layout">
    <div class="left-panel">
      <div class="form-wrapper">
        <div class="header">
          <div class="logo-container">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2L2 22H22L12 2Z" fill="white"/>
              </svg>
            </div>
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

.logo-icon {
  width: 32px;
  height: 32px;
  background-color: #1a41b5;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon svg {
  width: 20px;
  height: 20px;
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