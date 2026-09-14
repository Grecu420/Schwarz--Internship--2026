<template>
  <div class="property-creation-page-wrapper">
    <div class="header-container">
      <div class="back-button-wrapper">
        <OnyxButton
          label="Back"
          :icon="iconChevronLeft"
          mode="outline"
          class="back-btn"
          @click="goBack"
        />
      </div>
      <span class="category-tag">PROPERTY MANAGEMENT</span>
      <h1 class="page-title">Create New Property</h1>
      <p class="page-description">
        Fill in the details below to publish a new property listing.
      </p>
    </div>

    <PropertyForm @submit="handleSuccess" />
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

const uploadToCloud = false

const uploadImage = async (data: [string, File]): Promise<[string, string]> => {
  const [url, file] = data
  console.log('upload ', url)

  if (uploadToCloud) {
    const newUrl = await uploadImageToCloudinary(file, presetName)
    return [url, newUrl]
  } else {
    return [url, 'uploaded:' + url]
  }
}

const handleSuccess = async (
  createdProperty: Property,
  imagesToUpload: Map<string, File>,
  urlsToDelete: string[],
) => {
  try {
    // check authentication
    if (authStore.checkTokenExpiration()) {
      router.push('/login')
    }


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
.property-creation-page-wrapper {
  min-height: calc(100vh - 80px);
  background-color: #fafafa;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem 1.5rem 5rem 1.5rem;
}

.header-container {
  text-align: center;
  max-width: 520px;
  width: 100%;
  margin-bottom: 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.back-button-wrapper {
  align-self: flex-start;
  margin-bottom: 1rem;
}

.category-tag {
  font-size: 0.75rem;
  font-weight: 800;
  color: #2563eb;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.page-title {
  font-size: 2rem;
  font-weight: 800;
  color: #111827;
  margin: 0.25rem 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.page-description {
  font-size: 0.875rem;
  color: #6b7280;
  line-height: 1.5;
  margin: 0;
}
</style>