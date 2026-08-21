<script setup lang="ts">
import RegisterForm from '../components/RegisterForm.vue'
</script>

<template>
  <main class="page-container">
    <div class="header">
      <h1>Create an Account</h1>
      <p>Join us today! Please fill in your details below.</p>
    </div>
    
    <RegisterForm />
  </main>
</template>

<style scoped>
.page-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #f3f4f6;
  padding: 20px;
}

.header {
  text-align: center;
  margin-bottom: 24px;
}

.header h1 {
  font-size: 28px;
  color: #111827;
  margin-bottom: 8px;
}

.header p {
  color: #4b5563;
  font-size: 16px;
  margin: 0;
}
</style>