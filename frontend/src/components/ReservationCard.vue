<template>
  <OnyxCard class="reservation-card">
    <div class="reservation-header">
      <div class="price-tag">
        <span class="price-amount">${{ price }}</span>
        <span class="price-period"> / night</span>
      </div>
    </div>

    <!-- Onyx Date Picker -->
    <OnyxUnstableDatePickerV2
      label="Select Dates"
      v-model="intervalDates"
      :min="minDate"
      selectionMode="range"
      class="date-picker-input"
    />

    <!-- Onyx Action Button -->
    <OnyxButton label="Reserve this property" mode="default" color="primary" class="reserve-btn" />

    <!-- Check-in and Check-out Display Grid -->
    <div v-if="nights > 0" class="date-summary-box" >
      <div class="summary-cell">
        <span class="cell-label">CHECK-IN</span>
        <span class="cell-value">{{ startDateFormatted }}</span>
      </div>
      <div class="summary-cell border-left">
        <span class="cell-label">CHECK-OUT</span>
        <span class="cell-value">{{ endDateFormatted }}</span>
      </div>
    </div>

    <!-- Dynamic Price Calculation (Base Nights Total Only) -->
    <div v-if="nights > 0" class="price-breakdown">
      <div class="breakdown-row total-row">
        <span>Total ({{ nights }} {{ nights === 1 ? 'night' : 'nights' }})</span>
        <span>${{ totalPrice }}</span>
      </div>
    </div>
  </OnyxCard>
</template>

<script setup lang="ts">
import { type DateRange, OnyxButton, OnyxCard, OnyxUnstableDatePickerV2 } from 'sit-onyx'
import { computed, ref } from 'vue'

const props = defineProps<{
  price: number
  disabledDays?: Date[]
}>()

const intervalDates = ref<DateRange>()
const minDate = new Date()

// Date Extraction Helpers
const getStartDate = (range?: DateRange): Date | null => {
  if (!range) return null
  if (Array.isArray(range)) return range[0] ? new Date(range[0]) : null
  if (typeof range === 'object' && 'start' in range && range.start) return new Date(range.start)
  return null
}

const getEndDate = (range?: DateRange): Date | null => {
  if (!range) return null
  if (Array.isArray(range)) return range[1] ? new Date(range[1]) : null
  if (typeof range === 'object' && 'end' in range && range.end) return new Date(range.end)
  return null
}

const formatDate = (date: Date | null): string => {
  if (!date || isNaN(date.getTime())) return 'Select date'
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const year = date.getFullYear()
  return `${month}/${day}/${year}`
}

const startDateFormatted = computed(() => formatDate(getStartDate(intervalDates.value)))
const endDateFormatted = computed(() => formatDate(getEndDate(intervalDates.value)))

// Calculated Values
const nights = computed(() => {
  const start = getStartDate(intervalDates.value)
  const end = getEndDate(intervalDates.value)
  if (!start || !end) return 0
  const diffMs = end.getTime() - start.getTime()
  const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24))
  return days > 0 ? days : 0
})

const totalPrice = computed(() => props.price * nights.value)
</script>

<style scoped>
.reservation-card {
  background-color: #ffffff;
  color: #111827;
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.08);
  border: 1px solid #e5e7eb;
}

.reservation-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.25rem;
}

.price-amount {
  font-size: 1.5rem;
  font-weight: 800;
}

.price-period {
  font-size: 0.875rem;
  color: #6b7280;
}

.date-summary-box {
  display: flex;
  border: 1px solid #d1d5db;
  border-radius: 10px;
  overflow: hidden;
  margin-bottom: 1rem;
}

.summary-cell {
  flex: 1;
  padding: 1.625rem 0.875rem;
  display: flex;
  flex-direction: column;
}

.border-left {
  border-left: 1px solid #d1d5db;
}

.cell-label {
  font-size: 0.625rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  color: #374151;
  text-transform: uppercase;
}

.cell-value {
  font-size: 0.875rem;
  font-weight: 500;
  color: #111827;
  margin-top: 0.25rem;
}

.date-picker-input {
  margin-bottom: 1rem;
}

:deep(.reserve-btn) {
  width: 100% !important;
  background-color: #2563eb !important;
  color: #ffffff !important;
  border-radius: 8px !important;
  padding: 0.75rem !important;
  font-weight: 700 !important;
  font-size: 0.9375rem !important;
}

.price-breakdown {
  margin-top: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.breakdown-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.9375rem;
  color: #4b5563;
}

.total-row {
  padding-top: 0.875rem;
  border-top: 1px solid #e5e7eb;
  font-weight: 700;
  font-size: 1rem;
  color: #111827;
}
</style>
