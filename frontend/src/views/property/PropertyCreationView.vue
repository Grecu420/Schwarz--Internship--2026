<template>
  <div class="page-container">
    <header class="page-header">
      <OnyxButton
        label="Back"
        :icon="iconChevronLeft"
        mode="outline"
        class="back-btn"
        @click="goBack"
      />

      <div class="header-titles">
        <h1 class="page-title">Create New Property</h1>
        <p class="page-subtitle">Fill in the details below to publish a new property listing.</p>
      </div>
    </header>

    <main class="content-card">
      <PropertyForm @submit="handleSuccess" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { OnyxButton, useToast } from 'sit-onyx'
import { iconChevronLeft } from '@sit-onyx/icons'
import PropertyForm from '@/components/PropertyForm.vue'
import {
  CreatePropertyRequest,
  CreatePropertyResponse,
  type Property,
} from '@/generated/proto/property-api'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import { uploadImageToCloudinary } from '@/utils/cloudinary'
const presetName = import.meta.env.VITE_CLOUDINARY_PROFILE_PRESET
const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const goBack = () => {
  router.push('/properties')
}

const uploadImage = async (data: [string, File]): Promise<[string, string]> => {
  const [url, file] = data
  console.log('upload ', url)
  // const newUrl = await uploadImageToCloudinary(file, presetName)
  return [url, 'modded:' + url]
}

const handleSuccess = async (
  createdProperty: Property,
  imagesToUpload: Map<string, File>,
  urlsToDelete: string[],
) => {
  try {
    // upload images

    const substitutionEntries = await Promise.all(Array.from(imagesToUpload, uploadImage))
    const urlSubstitution = new Map<string, string>(substitutionEntries)

    const propertyRequest = CreatePropertyRequest.create({
      userId: createdProperty.userId || authStore.user?.id,
      name: createdProperty.name,
      description: createdProperty.description,
      address: createdProperty.address,
      price: createdProperty.price,
      location: createdProperty.location,
      imageUrls: createdProperty.imageUrls.map((url) => urlSubstitution.get(url) ?? url),
    })

    console.log(propertyRequest)
    const response = await api.post<CreatePropertyResponse>('/api/property', propertyRequest)

    console.log(response)
    toast.show({
      headline: 'Success',
      description: 'Property created successfully.',
      color: 'success',
    })

    router.push('/properties')
  } catch (error: any) {
    console.log(error)
    toast.show({
      headline: 'Creation Failed',
      description: error?.message || 'Failed to submit property request.',
      color: 'danger',
    })
  }
}
</script>

<style scoped>
.page-container {
  max-width: 720px;
  margin: 2rem auto;
  padding: 0 1rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.page-header {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  align-items: flex-start;
}

.header-titles {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--onyx-color-neutral-900, #0f172a);
  margin: 0;
}

.page-subtitle {
  font-size: 0.95rem;
  color: var(--onyx-color-neutral-600, #475569);
  margin: 0;
}

.content-card {
  background: var(--onyx-color-surface-default, #ffffff);
  border: 1px solid var(--onyx-color-neutral-200, #e2e8f0);
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

@media (max-width: 640px) {
  .content-card {
    padding: 1.25rem;
  }
}
</style>
