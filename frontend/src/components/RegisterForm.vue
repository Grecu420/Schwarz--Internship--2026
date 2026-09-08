<template>
  <OnyxForm @submit.prevent="handleRegister" class="register-form" novalidate>
    
    <div class="form-row">
      <OnyxInput
        v-model="formData.firstName"
        label="First Name"
        placeholder="John"
        required-marker="required"
        reserve-message-space
        :error="getFieldError('firstName')"
        @blur="v$.firstName.$touch()"
      >
        <template #leadingIcons>
          <OnyxIcon :icon="iconUser" />
        </template>
      </OnyxInput>

      <OnyxInput
        v-model="formData.lastName"
        label="Last Name"
        placeholder="Doe"
        required-marker="required"
        reserve-message-space
        :error="getFieldError('lastName')"
        @blur="v$.lastName.$touch()"
      >
        <template #leadingIcons>
          <OnyxIcon :icon="iconUser" />
        </template>
      </OnyxInput>
    </div>

    <OnyxInput
      v-model="formData.userName"
      label="Username"
      placeholder="johndoe123"
      required-marker="required"
      reserve-message-space
      :error="getFieldError('userName')"
      @blur="v$.userName.$touch()"
    >
      <template #leadingIcons>
        <OnyxIcon :icon="iconUser" />
      </template>
    </OnyxInput>

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
      placeholder="Minimum 8 characters"
      required-marker="required"
      reserve-message-space
      :error="getFieldError('password')"
      @blur="v$.password.$touch()"
    >
      <template #leadingIcons>
        <OnyxIcon :icon="iconPasswordLock" />
      </template>
    </OnyxInput>

    <OnyxButton 
      type="submit" 
      mode="default"
      class="submit-btn" 
      :loading="isLoading" 
      :disabled="isLoading"
      :label="isLoading ? 'Creating account...' : 'Create account'"
    />
  </OnyxForm>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useVuelidate } from '@vuelidate/core'
import { required, email, minLength, maxLength } from '@vuelidate/validators'
import { OnyxForm, OnyxInput, OnyxButton, OnyxIcon, useToast } from 'sit-onyx'
import { iconUser, iconMail, iconPasswordLock } from '@sit-onyx/icons'

import { CreateUserRequest, User } from '../generated/proto/user-api'

const toast = useToast()
const router = useRouter()

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

const getFieldError = (field: keyof typeof formData) => {
  const fieldValidation = v$.value[field]
  if (!fieldValidation || !fieldValidation.$error) return undefined

  if (fieldValidation.required?.$invalid) return 'This field is required.'
  if (fieldValidation.email?.$invalid) return 'Please enter a valid email address.'
  if (fieldValidation.minLength?.$invalid) return `Minimum ${fieldValidation.minLength.$params.min} characters required.`
  if (fieldValidation.maxLength?.$invalid) return `Maximum ${fieldValidation.maxLength.$params.max} characters allowed.`

  return 'Invalid value.'
}

const handleRegister = async () => {
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
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(jsonBody)
    })

    if (!response.ok) {
      let errorMessage = 'An error occurred during registration.'
      
      try {
        const errorData = await response.json()
        if (errorData && (errorData.message || errorData.error)) {
          errorMessage = errorData.message || errorData.error
        }
      } catch (e) {
        if (response.status === 409) {
          errorMessage = 'Username or email is already taken.'
        }
      }

      throw new Error(errorMessage)
    }

    toast.show({
      headline: 'Account created!',
      description: 'Welcome to StayNest. Please log in.',
      color: 'success'
    })
    
    router.push('/login')

  } catch (error: any) {
    toast.show({
      headline: 'Registration Failed',
      description: error.message || 'Error connecting to the server.',
      color: 'danger'
    })
    isLoading.value = false
  }
}
</script>

<style scoped>
.register-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: flex;
  gap: 1rem;
}

.form-row > * {
  flex: 1;
}

:deep(.submit-btn) {
  width: 100%;
  margin-top: 0.5rem;
  border-radius: 9999px;
}
</style>