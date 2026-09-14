<template>
  <OnyxForm @submit.prevent="handleSubmit" class="property-form" novalidate>
    <!-- Main Cover Image -->
    <div class="upload-block">
      <!-- Existing Cover Image Preview (Edit Mode) -->
      <div v-if="formData.mainImageUrl" class="existing-preview-card">
        <span class="preview-label">Main Image</span>
        <div class="image-wrapper">
          <OnyxImage
            :height="200"
            :width="300"
            :src="formData.mainImageUrl"
            alt="Existing Cover Image"
            shape="rounded"
          />
        </div>
        <OnyxIconButton label="delete" :icon="iconTrash" type="button" @click="deleteMainImage" />
      </div>
      <OnyxFileUpload
        v-else
        v-model="mainImageFile"
        :accept="imageAccept"
        label="Main image"
        maxSize="4MiB"
        required
        show-error
        size="medium"
        style="width: 30rem; max-width: 40%"
      />
    </div>

    <OnyxInput
      v-model="formData.name"
      label="Property Name"
      placeholder="Property Name"
      required-marker="required"
      reserve-message-space
      :error="getFieldError('name')"
      @blur="v$.name.$touch()"
    >
      <template #leadingIcons>
        <OnyxIcon :icon="iconHome" />
      </template>
    </OnyxInput>

    <OnyxTextarea
      v-model="formData.description"
      label="Description"
      placeholder="Detailed description of the property..."
      required-marker="required"
      reserve-message-space
      :error="getFieldError('description')"
      @blur="v$.description.$touch()"
    />

    <OnyxStepper
      v-model.number="formData.price"
      hideButtons
      label="Price (per night)"
      required-marker="required"
      reserve-message-space
      :error="getFieldError('price')"
      @blur="v$.price.$touch()"
    >
      <template #leadingIcons>
        <OnyxIcon :icon="iconTag" />
      </template>
    </OnyxStepper>

    <OnyxInput
      v-model="formData.address"
      label="Address"
      placeholder="Address"
      required-marker="required"
      reserve-message-space
      :error="getFieldError('address')"
      @blur="v$.address.$touch()"
    >
      <template #leadingIcons>
        <OnyxIcon :icon="iconMap" />
      </template>
    </OnyxInput>

    <!-- Location Selection -->
    <fieldset class="form-section">
      <legend class="section-title">Location</legend>
      <div class="form-row">
        <OnyxStepper
          label="Latitude"
          v-model="formData.location.lat"
          hideButtons
          required-marker="required"
          reserve-message-space
          :error="getLocationFieldError('lat')"
          @blur="v$.location.lat.$touch()"
        />

        <OnyxStepper
          label="Longitude"
          v-model="formData.location.long"
          hideButtons
          required-marker="required"
          reserve-message-space
          :error="getLocationFieldError('long')"
          @blur="v$.location.long.$touch()"
        />
      </div>

      <!-- Map Component Integration -->
      <MapComponent
        :key="mapKey"
        :initial-center="mapInitialCenter"
        :initial-zoom="15"
        @change="handleMapChange"
      />

      <OnyxButton
        type="button"
        mode="outline"
        class="location-btn"
        label="Detect Current Location"
        :loading="isGettingLocation"
        @click="() => fetchDeviceLocation(false)"
      />
    </fieldset>

    <!-- Property Images Gallery -->
    <fieldset class="form-section">
      <legend class="section-title">Gallery Images</legend>

      <div v-if="formData.galleryImageUrls.length > 0" class="gallery-grid">
        <div
          v-for="(url, index) in formData.galleryImageUrls"
          :key="index"
          class="existing-preview-card"
        >
          <div class="image-wrapper">
            <OnyxImage :height="100" :width="100" :src="url" alt="Gallery Image" shape="rounded" />
          </div>
          <OnyxIconButton
            label="Delete image"
            :icon="iconTrash"
            type="button"
            @click="removeGalleryImage(index)"
          />
        </div>
      </div>

      <OnyxFileUpload
        v-model="uploadedFileInput"
        :accept="imageAccept"
        label="Add gallery images"
        maxSize="4MiB"
        multiple
        size="medium"
        style="width: 30rem; max-width: 100%"
      />
    </fieldset>

    <OnyxButton
      type="submit"
      mode="default"
      class="submit-btn"
      :loading="isLoading"
      :disabled="isLoading"
      :label="submitButtonLabel"
    />
  </OnyxForm>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useVuelidate } from '@vuelidate/core'
import { required, minLength, minValue, numeric } from '@vuelidate/validators'
import {
  OnyxForm,
  OnyxInput,
  OnyxButton,
  OnyxIcon,
  useToast,
  OnyxStepper,
  OnyxFileUpload,
  OnyxIconButton,
  type FileType,
  OnyxImage,
  OnyxTextarea,
} from 'sit-onyx'
import { iconHome, iconMap, iconTag, iconTrash } from '@sit-onyx/icons'
import { Property } from '../generated/proto/property-api'
import MapComponent, { type LocationPayload } from './MapComponent.vue'
import { useAuthStore } from '@/stores/auth.ts'

const defaultLocation: [number, number] = [44.495, 26.08]

const props = withDefaults(
  defineProps<{
    initialData?: Property | null
    isEdit?: boolean
  }>(),
  {
    initialData: null,
    isEdit: false,
  },
)

const emit = defineEmits<{
  (e: 'submit', property: Property, imagesToUpload: Map<string, File>, urlsToDelete: string[]): void
}>()

const toast = useToast()
const authStore = useAuthStore()

const isLoading = ref(false)

// Form data
// ======================

const formData = reactive({
  name: '',
  description: '',
  address: '',
  price: null as number | null,
  location: {
    lat: null as number | null,
    long: null as number | null,
  },
  mainImageUrl: '',
  galleryImageUrls: [] as string[],
})

// File Storage
// ======================
const imageAccept = ['.png', '.jpg'] as FileType[]
const mainImageFile = ref<File | null>(null)
const uploadedFileInput = ref<File[]>([])

const imagesToUpload = ref<Map<string, File>>(new Map())
const urlsToDelete = ref<string[]>([])

// Pre-fill form data when initialData prop is provided or updated
watch(
  () => props.initialData,
  (data) => {
    if (data) {
      formData.name = data.name || ''
      formData.description = data.description || ''
      formData.address = data.address || ''
      formData.price = data.price ?? null
      formData.location.lat = data.location?.lat ?? null
      formData.location.long = data.location?.long ?? null

      mapKey.value++

      if (data.imageUrls && data.imageUrls.length > 0) {
        formData.mainImageUrl = data.imageUrls[0] ?? ''
        formData.galleryImageUrls = data.imageUrls.slice(1)
      } else {
        formData.mainImageUrl = ''
        formData.galleryImageUrls = []
      }
    }
  },
  { immediate: true, deep: true },
)

const submitButtonLabel = computed(() => {
  if (isLoading.value) {
    return props.isEdit ? 'Updating property...' : 'Creating property...'
  }
  return props.isEdit ? 'Update Property' : 'Create Property'
})

const rules = {
  name: { required, minLength: minLength(3) },
  description: { required, minLength: minLength(10) },
  address: { required },
  price: { required, minValue: minValue(1) },
  location: {
    lat: { required, numeric },
    long: { required, numeric },
  },
  mainImageUrl: { required },
}
const v$ = useVuelidate(rules, formData)

const getFieldError = (field: keyof typeof formData) => {
  const fieldValidation = v$.value[field]
  if (!fieldValidation || !fieldValidation.$error) return undefined

  if (fieldValidation.required?.$invalid) return 'This field is required.'
  if (fieldValidation.minLength?.$invalid)
    return `Minimum ${fieldValidation.minLength.$params.min} characters required.`
  if (fieldValidation.minValue?.$invalid)
    return `Value must be at least ${fieldValidation.minValue.$params.min}.`

  return 'Invalid value.'
}

// Map data
// ======================

const isGettingLocation = ref(false)
const mapKey = ref(0)

const mapInitialCenter = computed<[number, number]>(() => {
  if (formData.location.lat !== null && formData.location.long !== null) {
    return [formData.location.lat, formData.location.long]
  }
  return defaultLocation
})

function fetchDeviceLocation(silent = false) {
  if (!navigator.geolocation) {
    if (!silent) {
      toast.show({
        headline: 'Geolocation Unavailable',
        description: 'Your browser does not support automatic location detection.',
        color: 'danger',
      })
    }
    return
  }

  isGettingLocation.value = true
  navigator.geolocation.getCurrentPosition(
    (position) => {
      formData.location.lat = position.coords.latitude
      formData.location.long = position.coords.longitude
      mapKey.value++
      isGettingLocation.value = false

      if (!silent) {
        toast.show({
          headline: 'Location updated',
          description: 'Coordinates applied successfully.',
          color: 'success',
        })
      }
    },
    (error) => {
      isGettingLocation.value = false
      if (!silent) {
        toast.show({
          headline: 'Location Error',
          description: error.message || 'Failed to detect current location.',
          color: 'danger',
        })
      }
    },
  )
}

onMounted(() => {
  // Only attempt auto-detection if initial data wasn't provided
  if (!props.initialData?.location?.lat || !props.initialData?.location?.long) {
    fetchDeviceLocation(true)
  }
})

function handleMapChange(payload: LocationPayload) {
  const [lat, long] = payload.location
  formData.location.lat = lat
  formData.location.long = long

  if (payload.address) {
    formData.address = payload.address
    v$.value.address.$touch()
  }

  v$.value.location.lat.$touch()
  v$.value.location.long.$touch()
}

const getLocationFieldError = (coord: 'lat' | 'long') => {
  const coordValidation = v$.value.location[coord]
  if (!coordValidation || !coordValidation.$error) return undefined

  if (coordValidation.required?.$invalid) return 'Required.'
  if (coordValidation.numeric?.$invalid) return 'Must be a valid coordinate.'

  return 'Invalid.'
}

// Image handling
// ======================

const trackImageRemoval = (url?: string) => {
  if (!url) return
  if (url.startsWith('blob:')) {
    imagesToUpload.value.delete(url)
    URL.revokeObjectURL(url) // Free up browser memory
  } else {
    urlsToDelete.value.push(url)
  }
}

watch(mainImageFile, (newImage) => {
  if (!newImage) return

  trackImageRemoval(formData.mainImageUrl)
  formData.mainImageUrl = URL.createObjectURL(newImage)
  imagesToUpload.value.set(formData.mainImageUrl, newImage)
})

const deleteMainImage = () => {
  trackImageRemoval(formData.mainImageUrl)
  formData.mainImageUrl = ''
  mainImageFile.value = null
}

watch(uploadedFileInput, (newFiles) => {
  newFiles.forEach((file) => {
    const url = URL.createObjectURL(file)
    formData.galleryImageUrls.push(url)
    imagesToUpload.value.set(url, file)
  })

  uploadedFileInput.value = []
})

const removeGalleryImage = (index: number) => {
  const [removedUrl] = formData.galleryImageUrls.splice(index, 1)
  trackImageRemoval(removedUrl)
}

const handleSubmit = async () => {
  if (!authStore.user) {
    toast.show({ headline: 'Auth Error', description: 'User session is missing', color: 'danger' })
    return
  }

  v$.value.$touch()
  const isFormValid = await v$.value.$validate()

  if (!isFormValid) {
    toast.show({
      headline: 'Validation Error',
      description: 'Please correct all fields before submitting.',
      color: 'danger',
    })
    return
  }

  isLoading.value = true

  try {
    const allImages = [formData.mainImageUrl, ...formData.galleryImageUrls].filter(Boolean)

    const propertyObj = Property.create({
      id: props.initialData?.id,
      userId: authStore.user.id,
      name: formData.name,
      description: formData.description,
      address: formData.address,
      price: formData.price || 0,
      location: {
        lat: formData.location.lat!,
        long: formData.location.long!,
      },
      imageUrls: allImages,
    })

    emit('submit', propertyObj, imagesToUpload.value, urlsToDelete.value)
  } catch (error: any) {
    toast.show({
      headline: props.isEdit ? 'Update Failed' : 'Creation Failed',
      description: error.message || 'Error connecting to the server.',
      color: 'danger',
    })
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.property-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-row {
  display: flex;
  gap: 1rem;
}

.form-row > * {
  flex: 1;
}

.form-section {
  border: 1px solid var(--onyx-color-neutral-300, #e2e8f0);
  border-radius: 8px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.section-title {
  font-weight: 600;
  font-size: 0.9rem;
  padding: 0 0.5rem;
}

.location-btn {
  align-self: flex-start;
}

.gallery-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 0.5rem;
}

.existing-preview-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

:deep(.submit-btn) {
  width: 100%;
  margin-top: 0.5rem;
  border-radius: 9999px;
}
</style>
