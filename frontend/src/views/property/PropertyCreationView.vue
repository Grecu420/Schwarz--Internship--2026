<template>
  <div class="property-creation-page-wrapper">
    <div class="back-button-container">
      <OnyxButton
        label="Back"
        :icon="iconChevronLeft"
        mode="outline"
        class="back-btn"
        @click="goBack"
      />
    </div>

    <PropertyForm
      category="JOIN AS A HOST"
      title="List a New Property"
      description="Share your unique space with travelers from around the world. We make it simple."
      @submit="handleSuccess"
    />
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
    const response = await api.post<CreatePropertyResponse>('/api/property', propertyRequest)
    router.push('/properties')
    toast.show({
      headline: 'Success',
      description: 'Property created successfully.',
      color: 'success',
    })
  } catch (error: any) {
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
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2.5rem 2rem 5rem 2rem;
}

.back-button-container {
  width: 100%;
  max-width: 1140px;
  margin-bottom: 1.5rem;
}

:deep(.back-btn) {
  border-radius: 9999px !important;
}
</style>
