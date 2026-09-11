<template>
  <div class="profile-card">
    <div class="avatar-section">
      <input 
        type="file" 
        ref="fileInput" 
        accept="image/*" 
        style="display: none" 
        @change="onFileSelected" 
      />

      <div class="avatar-wrapper">
        <OnyxAvatar 
          :src="userForm.profileImageUrl" 
          :fullName="fullName" 
          size="96px" 
        />

        <div 
          v-if="userForm.profileImageUrl" 
          class="avatar-overlay" 
          @click="handleRemovePhoto"
          title="Remove Photo"
        >
          <OnyxIcon :icon="iconTrash" class="trash-icon" />
        </div>
      </div>

      <OnyxButton 
        type="button" 
        mode="outline" 
        label="Change Photo" 
        :icon="iconCamera"
        class="change-photo-btn"
        @click="triggerFileInput"
      />
    </div>

    <form @submit.prevent="handleSave" class="profile-form">
      <OnyxInput
        v-model="userForm.firstName"
        label="First Name"
        placeholder="Marcus"
      >
        <template #leadingIcons>
          <OnyxIcon :icon="iconUser" />
        </template>
      </OnyxInput>

      <OnyxInput
        v-model="userForm.lastName"
        label="Last Name"
        placeholder="Vance"
      >
        <template #leadingIcons>
          <OnyxIcon :icon="iconUser" />
        </template>
      </OnyxInput>

      <OnyxInput
        v-model="userForm.userName"
        label="Username"
        placeholder="marcusvance"
      >
        <template #leadingIcons>
          <OnyxIcon :icon="iconUser" />
        </template>
      </OnyxInput>

      <div class="form-actions">
        <OnyxButton 
          type="submit" 
          mode="default" 
          label="Save Changes" 
          :loading="isSaving"
          class="save-btn"
        />
      </div>
    </form>

    <hr class="card-divider" />

    <div class="danger-zone">
      <div class="danger-header">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#dc2626" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/>
          <line x1="12" y1="9" x2="12" y2="13"/>
          <line x1="12" y1="17" x2="12.01" y2="17"/>
        </svg>
        <span class="danger-title">DANGER ZONE</span>
      </div>

      <div class="danger-box">
        <p>
          Deleting your account is permanent. This will remove all your booking history, profile data, and saved favorites. This action cannot be undone.
        </p>
      </div>

      <button type="button" class="delete-account-btn" @click="handleDeleteAccount">
        Delete Account
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { OnyxInput, OnyxButton, OnyxIcon, OnyxAvatar, useToast } from 'sit-onyx'
import { iconUser, iconTrash, iconCamera } from '@sit-onyx/icons'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'
import { UpdateUserRequest, UpdateUserResponse, DeleteUserRequest, DeleteUserResponse } from '@/generated/proto/user-api'
import { uploadImageToCloudinary } from '@/utils/cloudinary'

const toast = useToast()
const router = useRouter()
const authStore = useAuthStore()
const isSaving = ref(false)

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const isPhotoRemoved = ref(false)

const userForm = reactive({
  firstName: '',
  lastName: '',
  userName: '',
  profileImageUrl: ''
})

const fullName = computed(() => {
  return `${userForm.firstName} ${userForm.lastName}`.trim() || userForm.userName || 'User'
})

onMounted(() => {
  if (authStore.user) {
    userForm.firstName = authStore.user.firstName || ''
    userForm.lastName = authStore.user.lastName || ''
    userForm.userName = authStore.user.userName || ''
    userForm.profileImageUrl = authStore.user.profileImageUrl || '' 
  }
})

const triggerFileInput = () => {
  fileInput.value?.click()
}

const onFileSelected = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (file) {
    selectedFile.value = file
    isPhotoRemoved.value = false 
    userForm.profileImageUrl = URL.createObjectURL(file) 
  }
}

const handleRemovePhoto = () => {
  if (selectedFile.value) {
    URL.revokeObjectURL(userForm.profileImageUrl)
    selectedFile.value = null
  }
  userForm.profileImageUrl = ''
  isPhotoRemoved.value = true
}

const handleSave = async () => {
  if (!authStore.user) {
    toast.show({ headline: 'Auth Error', description: 'User session is missing', color: 'danger' })
    return
  }

  isSaving.value = true
  try {
    let finalImageProfileUrl = userForm.profileImageUrl
    const updateMask = ['firstName', 'lastName', 'userName']

    if (selectedFile.value) {
      const presetName = import.meta.env.VITE_CLOUDINARY_PROFILE_PRESET 

      finalImageProfileUrl = await uploadImageToCloudinary(selectedFile.value, presetName)
      
      updateMask.push('profileImageUrl') 
    } 
    else if (isPhotoRemoved.value) {
      finalImageProfileUrl = ''
      updateMask.push('profileImageUrl')
    }

    const request = UpdateUserRequest.create({
      user: {
        id: authStore.user.id, 
        firstName: userForm.firstName,
        lastName: userForm.lastName,
        userName: userForm.userName,
        profileImageUrl: finalImageProfileUrl 
      },
      fieldMask: updateMask 
    })

    const jsonBody = UpdateUserRequest.toJSON(request)

    const response = await api.patch<UpdateUserResponse>('/api/user', jsonBody)

    if (response.data.user) {
      authStore.user = response.data.user
      localStorage.setItem('user', JSON.stringify(response.data.user))
    }

    if (selectedFile.value) {
      URL.revokeObjectURL(userForm.profileImageUrl)
      selectedFile.value = null
    }
    isPhotoRemoved.value = false

    toast.show({
      headline: 'Profile Updated',
      description: 'Your details have been saved successfully.',
      color: 'success'
    })
  } catch (error) {
    console.error('Error catched:', error)
    toast.show({
      headline: 'Error',
      description: 'Failed to update profile or upload image.',
      color: 'danger'
    })
  } finally {
    isSaving.value = false
  }
}

const handleDeleteAccount = async () => {
  if (!authStore.user) return

  if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
    try {
      const request = DeleteUserRequest.create({ 
        id: authStore.user.id 
      })

      await api.delete<DeleteUserResponse>(`/api/user?id=${authStore.user.id}`, { 
        data: request 
      })

      authStore.logout()
      router.push('/register')
      
    } catch (error) {
      toast.show({
        headline: 'Error',
        description: 'Failed to delete account.',
        color: 'danger'
      })
    }
  }
}
</script>

<style scoped>
.profile-card {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 20px;
  padding: 2.5rem;
  width: 100%;
  max-width: 640px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.875rem;
  margin-bottom: 2rem;
}

.avatar-wrapper {
  position: relative;
  width: 96px;
  height: 96px;
  border-radius: 50%;
  overflow: hidden; 
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background-color: rgba(220, 38, 38, 0.45); 
  backdrop-filter: blur(2px); 
  display: flex;
  justify-content: center;
  align-items: center;
  color: #ffffff;
  opacity: 0;
  cursor: pointer;
  transition: opacity 0.2s ease-in-out;
  z-index: 10;
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
}

.avatar-overlay:hover .trash-icon {
  transform: scale(1.15);
}

:deep(.change-photo-btn) {
  border-radius: 9999px !important;
  background-color: #ffffff !important;
  border-color: #e5e7eb !important;
  color: #374151 !important;
  font-weight: 500 !important;
  padding: 0.25rem 0.875rem !important;
}

:deep(.change-photo-btn:hover) {
  background-color: #f9fafb !important;
  border-color: #d1d5db !important;
}

.profile-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 0.5rem;
}

:deep(.save-btn) {
  background-color: #1e40af !important;
  border-radius: 8px !important;
  padding: 0.6rem 1.5rem !important;
  font-weight: 600 !important;
}

.card-divider {
  border: none;
  border-top: 1px solid #f3f4f6;
  margin: 2rem 0;
}

.danger-zone {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.danger-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.danger-title {
  font-size: 0.75rem;
  font-weight: 800;
  color: #dc2626;
  letter-spacing: 0.05em;
}

.danger-box {
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  border-radius: 10px;
  padding: 1rem 1.25rem;
}

.danger-box p {
  font-size: 0.8125rem;
  color: #991b1b;
  line-height: 1.5;
  margin: 0;
}

.delete-account-btn {
  align-self: flex-start;
  padding: 0.6rem 1.25rem;
  background-color: #ef4444;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.delete-account-btn:hover {
  background-color: #dc2626;
}
</style>