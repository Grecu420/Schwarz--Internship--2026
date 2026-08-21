<template>
  <div class="register-container">
    <form @submit.prevent="handleRegister" class="register-form">
      <h2>Create an Account</h2>

      
      <div class="form-group">
        <label for="firstName">First Name</label>
        <input id="firstName" v-model="formData.firstName" type="text" placeholder="John" />
        <span class="error-text" v-if="v$.firstName.$error">First name is required (min 2 chars).</span>
      </div>

     
      <div class="form-group">
        <label for="lastName">Last Name</label>
        <input id="lastName" v-model="formData.lastName" type="text" placeholder="Doe" />
        <span class="error-text" v-if="v$.lastName.$error">Last name is required (min 2 chars).</span>
      </div>

      
      <div class="form-group">
        <label for="userName">Username</label>
        <input id="userName" v-model="formData.userName" type="text" placeholder="johndoe123" />
        <span class="error-text" v-if="v$.userName.$error">Username must be between 4 and 20 chars.</span>
      </div>

      
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
        {{ isLoading ? 'Loading...' : 'Register' }}
      </button>

      
      <p v-if="errorMessage" class="global-error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="global-success">{{ successMessage }}</p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useVuelidate } from '@vuelidate/core'
import { required, email, minLength, maxLength } from '@vuelidate/validators'


import { CreateUserRequest, User } from '../generated/proto/user-api'

const formData = reactive({
  firstName: '',
  lastName: '',
  userName: '',
  email: '',
  password: ''
})

const rules = {
  firstName: { required, minLength: minLength(2) },
  lastName: { required, minLength: minLength(2) },
  userName: { required, minLength: minLength(4), maxLength: maxLength(20) },
  email: { required, email },
  password: { required, minLength: minLength(8) }
}

const v$ = useVuelidate(rules, formData)
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const handleRegister = async () => {

  const isFormValid = await v$.value.$validate()
  if (!isFormValid) return

  errorMessage.value = ''
  successMessage.value = ''
  isLoading.value = true

  try {
  
    const userPayload = User.create({
      firstName: formData.firstName,
      lastName: formData.lastName,
      userName: formData.userName,
      email: formData.email,
      password: formData.password
    })

    const requestObj = CreateUserRequest.create({ user: userPayload })
    const jsonBody = CreateUserRequest.toJSON(requestObj)

const response = await fetch('/api/user', { 
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(jsonBody)
})

    if (!response.ok) {
      throw new Error(`Server responded with status: ${response.status}`)
    }

   
    successMessage.value = 'User was successfully created!'
    Object.assign(formData, { firstName: '', lastName: '', userName: '', email: '', password: '' })
    v$.value.$reset()

  } catch (error: any) {
    errorMessage.value = error.message || 'Error connecting to the server.'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  font-family: Arial, sans-serif;
}

.register-form {
  background: #f9f9f9;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.register-form h2 {
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