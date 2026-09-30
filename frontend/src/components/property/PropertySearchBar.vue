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
          <!-- 1. Location Selector -->
          <OnyxSelect
            v-model="filters.location"
            class="filter-field"
            label="Location"
            :options="cityOptions"
            placeholder="Select city"
            list-label="Cities"
            withSearch
          />

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
            mode="outline"
            class="clear-button"
            @click="handleClear"
          />
        </div>
      </OnyxUnstableSearch>

      <!-- Submit Button -->
      <OnyxButton label="Search" color="primary" class="submit-button" @click="handleSubmit" />
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
  price?: [number, number]
  dates?: DateRange
  sort?: string
}

const filters = ref<FilterState>({})

const props = defineProps<{
  cities: City[]
  price: [number, number]
  sortOptions?: string[]
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

const sortOptions = [
  { value: 'name', label: 'Name' },
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
]

// Map cities prop to OnyxSelect option format ({ value, label })
const cityOptions = computed(() => {
  if (!props.cities) return []
  return props.cities.map((city) => ({ value: city.name, label: city.name }))
})

const handleClear = () => {
  filters.value = {}
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
  width: 100%;
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
  background-color: #1e40af !important;
  border-radius: 8px !important;
  padding: 0.6rem 1.5rem !important;
  font-weight: 600 !important;
  margin-left: auto;
}

.filters-vertical-container {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  width: 100%;
  padding-top: 1rem;
}

.filter-field,
.clear-button {
  width: 100%;
  max-width: 400px;
}

:deep(.price-slider .onyx-stepper),
:deep(.price-slider .onyx-input-wrap),
:deep(.price-slider input) {
  width: 4.5rem !important;
  min-width: 4.5rem !important;
}
</style>