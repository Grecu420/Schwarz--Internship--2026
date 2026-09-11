<template>
  <OnyxForm @submit.prevent="handleSubmit" class="property-form" novalidate>
    <!-- Main Cover Image -->
    <div class="upload-block">
      <!-- Existing Cover Image Preview (Edit Mode) -->
      <div v-if="formData.mainImageUrl" class="existing-preview-card">
        <span class="preview-label">Main Image</span>
        <div class="image-wrapper">
          <OnyxImage
            :height="100"
            :width="100"
            :src="formData.mainImageUrl"
            alt="Existing Cover Image"
            shape="rounded"
          />
        </div>
        <OnyxIconButton label="delete" :icon="iconTrash" @click="deleteMainImage" />
      </div>
      <OnyxFileUpload
        v-else
        v-model="mainImageFile"
        :accept="imageAccept"
        label="Main image"
        maxSize="4MiB"
        required
        show-error
        size="small"
        style="width: 30rem; max-width: 100%"
      />
    </div>

    <MapComponent />

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
      <legend class="section-title">Location Coordinates</legend>
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

      <OnyxButton
        type="button"
        mode="outline"
        class="location-btn"
        label="Detect Current Location"
        :loading="isGettingLocation"
        @click="getCurrentLocation"
      />
    </fieldset>

    <!-- Property Images -->
    <fieldset class="form-section">
      <legend class="section-title">Property Images</legend>

      <!-- Additional Gallery Images -->
      <!-- <div class="gallery-section">
        <span class="sub-label">Gallery Images (Optional)</span>
        <div class="image-input-row">
          <OnyxInput
            v-model="galleryUrlInput"
            label="Gallery Image URL"
            placeholder="https://example.com/photo.jpg"
            reserve-message-space
          >
            <template #leadingIcons>
              <OnyxIcon :icon="iconPicture" />
            </template>
          </OnyxInput>

          <OnyxButton
            type="button"
            label="Add to Gallery"
            :disabled="!galleryUrlInput"
            @click="addGalleryImage"
          />
        </div>

        <div v-if="formData.galleryImageUrls.length > 0" class="image-preview-grid">
          <div
            v-for="(url, index) in formData.galleryImageUrls"
            :key="url + index"
            class="image-card"
          >
            <img :src="url" alt="Gallery image preview" />
            <button
              type="button"
              class="remove-btn"
              @click="removeGalleryImage(index)"
              aria-label="Remove image"
            >
              &times;
            </button>
          </div>
        </div>
      </div> -->
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
import { ref, reactive, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
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
import { iconHome, iconMap, iconTag, iconPicture, iconX, iconTrash } from '@sit-onyx/icons'

import { Property } from '../generated/proto/property-api'
import MapComponent from './MapComponent.vue'

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
  (e: 'submit', property: Property): void
}>()

const toast = useToast()
const router = useRouter()

const galleryUrlInput = ref('')
const isLoading = ref(false)
const isGettingLocation = ref(false)

const imageAccept = ['.png', '.jpg'] as FileType[]

const mainImageFile = ref<File | null>(null)
const galleryImageFiles = ref<File[]>([])

const deletedImageUrls = ref<string[]>([])
const imagesToUpload = ref<File[]>([])

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

watch(mainImageFile, (newImage, oldImage) => {
  console.log(newImage?.name)
  if (newImage) {
    formData.mainImageUrl = URL.createObjectURL(newImage)
  }

  console.log(formData.mainImageUrl)
})

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

const isBlobURL = (url: string) => {
  return url.startsWith('blob:')
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

const getLocationFieldError = (coord: 'lat' | 'long') => {
  const coordValidation = v$.value.location[coord]
  if (!coordValidation || !coordValidation.$error) return undefined

  if (coordValidation.required?.$invalid) return 'Required.'
  if (coordValidation.numeric?.$invalid) return 'Must be a valid coordinate.'

  return 'Invalid.'
}

const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    toast.show({
      headline: 'Geolocation Unavailable',
      description: 'Your browser does not support automatic location detection.',
      color: 'danger',
    })
    return
  }

  isGettingLocation.value = true
  navigator.geolocation.getCurrentPosition(
    (position) => {
      formData.location.lat = position.coords.latitude
      formData.location.long = position.coords.longitude
      isGettingLocation.value = false
      toast.show({
        headline: 'Location updated',
        description: 'Coordinates applied successfully.',
        color: 'success',
      })
    },
    (error) => {
      isGettingLocation.value = false
      toast.show({
        headline: 'Location Error',
        description: error.message || 'Failed to detect current location.',
        color: 'danger',
      })
    },
  )
}

const isValidUrl = (url: string) => {
  const urlPattern = /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/
  return urlPattern.test(url)
}
const getMainImageError = () => {
  if (v$.value.mainImage.$error) return 'A main cover image is required.'
  return undefined
}

const deleteMainImage = () => {
  formData.mainImageUrl = ''
  mainImageFile.value = null
}

const addGalleryImage = () => {
  if (!galleryUrlInput.value) return

  if (!isValidUrl(galleryUrlInput.value)) {
    toast.show({
      headline: 'Invalid URL',
      description: 'Please enter a valid image URL.',
      color: 'danger',
    })
    return
  }

  formData.galleryImageUrls.push(galleryUrlInput.value.trim())
  galleryUrlInput.value = ''
}

const removeGalleryImage = (index: number) => {
  formData.galleryImageUrls.splice(index, 1)
}

const handleImageError = () => {
  toast.show({
    headline: 'Image Load Warning',
    description: 'Unable to preview the main image URL.',
    color: 'warning',
  })
}

const handleSubmit = async () => {
  v$.value.$touch()
  const isFormValid = await v$.value.$validate()

  console.log(mainImageFile.value)

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
    const propertyObj = Property.create({
      id: props.initialData?.id || undefined,
      name: formData.name,
      description: formData.description,
      address: formData.address,
      price: formData.price || 0,
      location: {
        lat: formData.location.lat!,
        long: formData.location.long!,
      },
      imageUrls: [formData.mainImageUrl, ...formData.galleryImageUrls],
    })
    toast.show({
      headline: props.isEdit ? 'Property Updated!' : 'Property Created!',
      description: `Your property has been successfully ${props.isEdit ? 'updated' : 'created'}.`,
      color: 'success',
    })

    emit('submit', propertyObj)
    router.push('/properties')
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

.image-input-row {
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
}

.image-input-row > :first-child {
  flex: 1;
}

.image-preview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 0.75rem;
  margin-top: 0.5rem;
}

.image-card {
  position: relative;
  aspect-ratio: 1;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid var(--onyx-color-neutral-300, #e2e8f0);
}

.image-card img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.remove-btn {
  position: absolute;
  top: 2px;
  right: 2px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  border: none;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 14px;
}

.validation-error {
  color: var(--onyx-color-danger-500, #ef4444);
  font-size: 0.8rem;
  margin-top: 0.25rem;
}

:deep(.submit-btn) {
  width: 100%;
  margin-top: 0.5rem;
  border-radius: 9999px;
}
</style>
