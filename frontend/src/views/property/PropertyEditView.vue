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
      category="CURRENT LISTING"
      title="Edit your Property Listing"
      description="Make updates to yout pricing, description, or listing details. Keep information fresh"
      @submit="handleSuccess"
      :initial-data="initialProperty"
    />
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { OnyxButton, useToast } from 'sit-onyx'
import { iconChevronLeft } from '@sit-onyx/icons'
import PropertyForm from '@/components/PropertyForm.vue'
import {
  GetPropertyResponse,
  UpdatePropertyRequest,
  UpdatePropertyResponse,
  type Property,
} from '@/generated/proto/property-api'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import { uploadImageToCloudinary } from '@/utils/cloudinary'
import { onMounted, ref } from 'vue'

const presetName = import.meta.env.VITE_CLOUDINARY_PROFILE_PRESET
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const toast = useToast()
const initialProperty = ref<Property | null>(null)
const isLoading = ref(true)

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

onMounted(async () => {
  const id = route.params.id

  try {
    const response = await api.get<GetPropertyResponse>(`/api/property?id=${id}`)
    initialProperty.value = response.data.property ?? null
    if (initialProperty.value === null) {
      goBack()
    }
  } catch (error) {
    console.error('Failed to load property details:', error)
  } finally {
    isLoading.value = false
  }
})

const handleSuccess = async (
  updatedProperty: Property,
  imagesToUpload: Map<string, File>,
  urlsToDelete: string[],
) => {
  try {
    // upload images
    const substitutionEntries = await Promise.all(Array.from(imagesToUpload, uploadImage))
    const urlSubstitution = new Map<string, string>(substitutionEntries)

    updatedProperty.imageUrls.map((url) => urlSubstitution.get(url) ?? url)

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

    console.log(request)
    const response = await api.patch<UpdatePropertyResponse>('/api/property', request)

    console.log(response)
    toast.show({
      headline: 'Success',
      description: 'Property updated successfully.',
      color: 'success',
    })

    router.push('/properties')
  } catch (error: any) {
    console.log(error)
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
