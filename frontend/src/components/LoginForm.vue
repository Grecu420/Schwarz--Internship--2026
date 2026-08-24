<template>
  <div class="login-container">
    <form @submit.prevent="handleLogin" class="login-form">

      <div class="form-group">
        <label for="email">Email Address</label>
        <input id="email" v-model="formData.email" type="email" placeholder="john@example.com" />
        <span class="error-text" v-if="v$.email.$error">Please enter a valid email address.</span>
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input id="password" v-model="formData.password" type="password" placeholder="••••••••" />
        <span class="error-text" v-if="v$.password.$error">Password must be at least 8 characters.</span>
      </div>

      <button type="submit" class="submit-btn" :disabled="isLoading">
        {{ isLoading ? 'Loading...' : 'Login' }}
      </button>

      <p v-if="errorMessage" class="global-error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="global-success">{{ successMessage }}</p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useVuelidate } from '@vuelidate/core'
import { required, email, minLength } from '@vuelidate/validators'

import { LoginRequest, LoginResponse } from '@/generated/proto/auth-api'
import { useRouter } from 'vue-router';
import api from '@/api';


const router = useRouter();

const formData = reactive({
  email: '',
  password: ''
})

const rules = {
  email: { required, email },
  password: { required, minLength: minLength(8) }
}

const v$ = useVuelidate(rules, formData)
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const handleLogin = async () => {

  // Validate before request
  const isFormValid = await v$.value.$validate()
  if (!isFormValid) return

  errorMessage.value = ''
  successMessage.value = ''
  isLoading.value = true

  try {

    const loginRequest = LoginRequest.create({
      email: formData.email,
      password: formData.password
    })
    // Send request

    const response = await api.post<LoginResponse>('/api/login', loginRequest);
    const token = response.data.JWT;

    localStorage.setItem('jwt_token', token);

    successMessage.value = 'Successful login'
    // Exit login page
    router.push("/")

  } catch (error: any) {
    errorMessage.value = error.message || 'Error connecting to the server.'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  font-family: Arial, sans-serif;
}

.login-form {
  background: #f9f9f9;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.login-form h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: bold;
  color: #555;
}

input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-sizing: border-box;
}

input:focus {
  border-color: #4CAF50;
  outline: none;
}

.error-text {
  color: #d9534f;
  font-size: 0.85rem;
  margin-top: 0.25rem;
  display: block;
}

.submit-btn {
  width: 100%;
  padding: 0.75rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  margin-top: 1rem;
  transition: background-color 0.3s;
}

.submit-btn:hover:not(:disabled) {
  background-color: #45a049;
}

.submit-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.global-error {
  color: #d9534f;
  text-align: center;
  margin-top: 1rem;
  font-weight: bold;
}

.global-success {
  color: #5cb85c;
  text-align: center;
  margin-top: 1rem;
  font-weight: bold;
}
</style>