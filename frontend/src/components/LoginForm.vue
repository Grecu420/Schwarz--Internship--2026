<template>
  <OnyxForm @submit.prevent="handleLogin" class="login-form" novalidate>
    
    <OnyxInput
      v-model="formData.email"
      type="email"
      label="Email address"
      placeholder="john@example.com"
      required-marker="required"
      reserve-message-space
      :error="getFieldError('email')"
      @blur="v$.email.$touch()"
    >
      <template #leadingIcons>
        <OnyxIcon :icon="iconMail" />
      </template>
    </OnyxInput>

    <OnyxInput
      v-model="formData.password"
      type="password"
      label="Password"
      placeholder="••••••••"
      required-marker="required"
      reserve-message-space
      :error="getFieldError('password')"
      @blur="v$.password.$touch()"
    >
      <template #leadingIcons>
        <OnyxIcon :icon="iconPasswordLock" />
      </template>
    </OnyxInput>

    <div class="form-actions">
      <router-link to="/forgot-password" class="forgot-link">Forgot password?</router-link>
    </div>

    <OnyxButton 
      type="submit" 
      mode="default"
      class="submit-btn" 
      :loading="isLoading" 
      :disabled="isLoading"
      :label="isLoading ? 'Logging in...' : 'Log in'"
    />
  </OnyxForm>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useVuelidate } from '@vuelidate/core'
import { required, email } from '@vuelidate/validators'
import { OnyxForm, OnyxInput, OnyxButton, OnyxIcon, useToast } from 'sit-onyx'
import { iconMail, iconPasswordLock } from '@sit-onyx/icons'

import axios from 'axios'
import api from '@/api'
import { LoginRequest, LoginResponse } from '@/generated/proto/auth-api'
import { GetUserResponse } from '@/generated/proto/user-api'
import { useAuthStore } from '@/stores/auth'

const toast = useToast()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const formData = reactive({
  email: '',
  password: ''
})

const rules = {
  email: { required, email },
  password: { required }
}

const v$ = useVuelidate(rules, formData)
const isLoading = ref(false)

const getFieldError = (field: keyof typeof formData) => {
  const fieldValidation = v$.value[field]
  if (!fieldValidation || !fieldValidation.$error) return undefined

  if (fieldValidation.required?.$invalid) return 'This field is required.'
  if (fieldValidation.email?.$invalid) return 'Please enter a valid email address.'

  return 'Invalid value.'
}

const handleLogin = async () => {
  v$.value.$touch()
  const isFormValid = await v$.value.$validate()
  
  if (!isFormValid) {
    toast.show({
      headline: 'Validation Error',
      description: 'Please correct the highlighted errors before submitting.',
      color: 'danger'
    })
    return
  }

  isLoading.value = true

  try {
    const loginRequest = LoginRequest.create({
      email: formData.email,
      password: formData.password
    })
    
    const loginResponse = await api.post<LoginResponse>('/api/login', loginRequest)
    const token = loginResponse.data.JWT

    const userResponse = await api.get<GetUserResponse>(`/api/user?email=${formData.email}`, {
      headers: { Authorization: `Bearer ${token}` }
    })

    const user = userResponse.data.user

    if (user) {
      authStore.setSession(token, user)

      const redirectPath = (route.query.redirect as string) || '/'
      router.push(redirectPath)
    }

  } catch (error: unknown) {
    let errorMessage = 'An unexpected error occurred.'

    if (axios.isAxiosError(error)) {
      const code = error.response?.data?.code
      if (code === 5 || code === 7) {
        errorMessage = 'Wrong email or password.'
      } else {
        errorMessage = 'Error connecting to the server.'
      }
    } else if (error instanceof Error) {
      errorMessage = error.message
    }

    toast.show({
      headline: 'Login Failed',
      description: errorMessage,
      color: 'danger'
    })
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-form {
  display: flex;
  flex-direction: column;
  gap: 0.25rem; 
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: -0.5rem;
  margin-bottom: 0.5rem;
}

.forgot-link {
  color: var(--onyx-color-base-primary-500, #1a41b5);
  font-size: 0.875rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s;
}

.forgot-link:hover {
  color: var(--onyx-color-base-primary-600, #143393);
  text-decoration: underline;
}

:deep(.submit-btn) {
  width: 100%;
  margin-top: 0.5rem;
  border-radius: 9999px; 
}
</style>