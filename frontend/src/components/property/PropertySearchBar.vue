<template>
  <div class="search-bar-card">
    <div class="search-row">
      <OnyxUnstableSearch
        v-model:show-filters="showFilters"
        v-model="filters.search"
        placeholder="Search property name..."
        class="search-input"
      >
        <div class="filters-vertical-container">
          <!-- 1. Location & Radius Selector -->
          <div class="location-radius-group">
            <OnyxSelect
              v-model="filters.location"
              class="filter-field location-select"
              label="Location"
              :options="cityOptions"
              placeholder="Select city"
              list-label="Cities"
              withSearch
            />

            <OnyxSelect
              v-model="filters.radius"
              class="filter-field radius-select"
              label="Radius"
              :options="radiusOptions"
              placeholder="+ 0 km"
              list-label="Radius Options"
              :disabled="!filters.location"
            />
          </div>

          <!-- 2. Date Range Picker for Check-in & Check-out -->
          <OnyxUnstableDatePickerV2
            label="Check-In and Check-Out Dates"
            selectionMode="range"
            :min="minDate"
            v-model="filters.dates"
            class="filter-field"
          />

          <!-- 3. Price Range Slider -->
          <div class="filter-field">
            <OnyxSlider
              label="Price Range"
              mode="range"
              v-model="priceRange"
              control="input"
              :min="props.price[0] ?? 0"
              :max="props.price[1] ?? 1000"
              class="price-slider"
            />
          </div>

          <!-- 4. Sort Selector (At the end) -->
          <OnyxSelect
            v-model="filters.sort"
            class="filter-field"
            label="Sort By"
            :options="sortOptions"
            placeholder="Select sort order"
            list-label="Sort options"
          />

          <!-- 5. Clear Button -->
          <OnyxButton
            label="Clear filters"
            class="clear-button"
            @click="handleClear"
          />
        </div>
      </OnyxUnstableSearch>

      <!-- Submit Button -->
      <OnyxButton label="Search" color="primary" class="submit-btn" @click="handleSubmit" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  OnyxButton,
  OnyxSlider,
  OnyxUnstableSearch,
  OnyxSelect,
  OnyxUnstableDatePickerV2,
  type DateRange,
} from 'sit-onyx'
import type { City } from '@/utils/cities'

export type FilterState = {
  search?: string
  location?: string
  radius?: number
  price?: [number, number]
  dates?: DateRange
  sort?: string
}

const filters = ref<FilterState>({ radius: 5000 })

const props = defineProps<{
  cities: City[]
  price: [number, number]
  sortOptions: { value: string; label: string }[]
}>()

const priceRange = ref(props.price)

const emit = defineEmits<{
  (e: 'submit', filters: FilterState): void
}>()

const minDate = computed(() => {
  const d = new Date()
  d.setHours(0, 0, 0, 0)
  return d
})
const showFilters = ref(false)

const radiusOptions = [
  { value: 5000, label: '+ 5 km' },
  { value: 10000, label: '+ 10 km' },
  { value: 20000, label: '+ 20 km' },
  { value: 50000, label: '+ 50 km' },
  { value: 100000, label: '+ 100 km' },
]

// Map cities prop to OnyxSelect option format ({ value, label })
const cityOptions = computed(() => {
  if (!props.cities) return []
  return props.cities.map((city) => ({ value: city.name, label: city.name }))
})

const handleClear = () => {
  filters.value = { radius: 5000 }
  priceRange.value = props.price
  handleSubmit()
}

const handleSubmit = () => {
  filters.value.price = priceRange.value
  emit('submit', filters.value)
}
</script>

<style scoped>
.search-bar-card {
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 1.25rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.search-row {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  width: 100%;
}

.search-input {
  flex: 1;
}

:deep(.submit-btn) {
  border-radius: 10px !important;
  padding: 0.4rem 1.2rem !important;
  font-weight: 600 !important;
}

.filters-vertical-container {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  width: 100%;
  padding-top: 1rem;
}

.location-radius-group {
  display: flex;
  gap: 0.75rem;
  width: 100%;
}

.location-select {
  flex: 3;
  max-width: none;
}

.radius-select {
  flex: 2;
  max-width: none;
}

.filter-field,
.clear-button {
  width: 100%;
}

:deep(.clear-button){
  border-radius: 10px;
}

:deep(.price-slider .onyx-stepper),
:deep(.price-slider .onyx-input-wrap),
:deep(.price-slider input) {
  width: 4.5rem !important;
  min-width: 4.5rem !important;
}

:deep(.onyx-form-element-v2__content) {
  border-radius: 50px !important;
  border: 1px solid lightgray !important; 
  overflow: hidden !important; 
  display: flex !important;
  align-items: center !important;
  gap: 0 !important; 
}

:deep(.onyx-form-element-v2__input-container) {
  border: none !important;
  box-shadow: none !important;
  background: transparent !important;
}

:deep(.onyx-form-element-v2__input) {
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
  background: transparent !important;
}

:deep(.onyx-form-element-v2__slot--trailing) {
  border: none !important; 
  background: transparent !important;
}

:deep(.onyx-form-element-v2__slot--trailing .onyx-form-element-action__button) {
  border-radius: 0 !important;
  border: none !important;
  background: transparent !important;
  box-shadow: none !important;
}




</style>
