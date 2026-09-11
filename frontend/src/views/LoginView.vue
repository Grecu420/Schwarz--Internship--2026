<template>
  <main class="split-layout">
    <div class="left-panel">
      <div class="form-wrapper">
        <div class="header">
  
          <Logo/>
          <h1>Welcome back</h1>
          <p>Discover unique stays, tailored experiences, and friendly hosts.</p>
        </div>

        <LoginForm />

        <p class="register-prompt">
          Don't have an account? <router-link to="/register" class="register-link">Register</router-link>
        </p>
      </div>
    </div>

    <div class="right-panel"></div>
  </main>
</template>

<script setup lang="ts">
import Logo from '@/components/Logo.vue';
import LoginForm from '../components/LoginForm.vue'
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

.header h1 {
  font-size: 2rem;
  margin: 0 0 0.5rem 0;
  color: #111;
}

.header p {
  color: #666;
  margin: 0;
}

.register-prompt {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 0.9rem;
}

.register-prompt a {
  color: #1a41b5;
  font-weight: bold;
  text-decoration: none;
}

.right-panel {
  flex: 1;
  background: url('https://images.unsplash.com/photo-1510798831971-661eb04b3739?q=80&w=1920&auto=format&fit=crop') center/cover no-repeat;
}

@media (max-width: 900px) {
  .right-panel {
    display: none;
  }
}

.register-link {
  color: #2c5282; 
  text-decoration: none;
  font-weight: 600;
  display: inline-block;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.register-link:hover {
  color: #1a41b5; 
  transform: translateY(-1px) scale(1.03); 
  filter: brightness(1.1);
}

.register-link:active {
  color: #0f2671;
  transform: translateY(0) scale(0.98);
}
</style>