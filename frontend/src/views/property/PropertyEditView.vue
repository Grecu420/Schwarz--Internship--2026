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

    <p v-if="isLoading"> Loading...</p>

    <PropertyForm
      v-if="!isLoading && initialProperty"
      category="CURRENT LISTING"
      title="Edit your Property Listing"
      description="Make updates to your pricing, description, or listing details. Keep information fresh"
      @submit="handleSuccess"
      @delete="handleDelete"
      :initial-data="initialProperty"
      :is-edit="true"
    />

  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { OnyxButton, useToast } from 'sit-onyx'
import { iconChevronLeft } from '@sit-onyx/icons'
import PropertyForm from '@/components/PropertyForm.vue'
import {
  DeletePropertyResponse,
  GetPropertyResponse,
  UpdatePropertyRequest,
  UpdatePropertyResponse,
  type Property,
} from '@/generated/proto/property-api'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import { uploadImageToCloudinary } from '@/utils/cloudinary'
import { onMounted, ref } from 'vue'

const presetName = import.meta.env.VITE_CLOUDINARY_PROPERTY_PRESET
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const toast = useToast()
const initialProperty = ref<Property | null>(null)
const isLoading = ref(true)

const goBack = () => {
  router.push('/properties')
}

const uploadToCloud = true

const uploadImage = async (data: [string, File]): Promise<[string, string]> => {
  const [url, file] = data

  if (uploadToCloud) {
    const newUrl = await uploadImageToCloudinary(file, presetName)
    return [url, newUrl]
  } else {
    return [url, 'uploaded:' + url]
  }
}

onMounted(async () => {
  const id = route.params.id
  isLoading.value = true
  try {
    const response = await api.get<GetPropertyResponse>(`/api/property?id=${id}`)
    initialProperty.value = response.data.property ?? null
    if (initialProperty.value === null) {
      goBack()
    }
  } catch (error: any) {
    toast.show({
      headline: 'Failed to load property details',
      description: error?.message || 'Failed to send property request.',
      color: 'danger',
    })
    goBack()
  } finally {
    isLoading.value = false
  }
})

const handleDelete = async (propertyID: number, urlsToDelete: string[]) => {
  isLoading.value = true
  try {
    const response = await api.delete<DeletePropertyResponse>(`/api/property?id=${propertyID}`)
    toast.show({
      headline: 'Success',
      description: 'Property deleted successfully.',
      color: 'success',
    })
    goBack()
  } catch (error: any) {
    toast.show({
      headline: 'Delete Failed',
      description: error?.message || 'Failed to delete property listing.',
      color: 'danger',
    })
  } finally {
    isLoading.value = false
  }
}

const handleSuccess = async (
  updatedProperty: Property,
  imagesToUpload: Map<string, File>,
  urlsToDelete: string[],
) => {
  try {
    // upload images
    const substitutionEntries = await Promise.all(Array.from(imagesToUpload, uploadImage))
    const urlSubstitution = new Map<string, string>(substitutionEntries)

    updatedProperty.imageUrls = updatedProperty.imageUrls.map(
      (url) => urlSubstitution.get(url) ?? url,
    )

    const fieldMask: string[] = ['name', 'description', 'address', 'price', 'location', 'imageUrls']

    const initProperty = initialProperty.value

    if (initProperty === null) {
      goBack()
      return
    }
    const request = UpdatePropertyRequest.create({
      property: updatedProperty,
      fieldMask: fieldMask,
    })

    const json = UpdatePropertyRequest.toJSON(request)
    const response = await api.patch<UpdatePropertyResponse>('/api/property', json)
    router.push('/properties')
    toast.show({
      headline: 'Success',
      description: 'Property updated successfully.',
      color: 'success',
    })
  } catch (error: any) {
    toast.show({
      headline: 'Update Failed',
      description: error?.message || 'Failed to submit update request.',
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
