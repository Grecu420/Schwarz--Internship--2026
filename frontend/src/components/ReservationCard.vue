<template>
  <div class="reservation-card">
    <div class="reservation-header">
      <div class="price-tag">
        <span class="price-amount">${{ price }}</span>
        <span class="price-period"> / night</span>
      </div>
    </div>
    <p>Reservation component to be implemented.</p>

    <!-- Check-in and Check-out Display Grid -->
    <div class="date-summary-box" v-if="intervalDates">
      <div class="summary-cell">
        <span class="cell-label">CHECK-IN</span>
        <span class="cell-value">{{ formatDate(intervalDates.start) }}</span>
      </div>
      <div class="summary-cell border-left">
        <span class="cell-label">CHECK-OUT</span>
        <span class="cell-value">{{ formatDate(intervalDates.end || null) }}</span>
      </div>
    </div>

    <OnyxUnstableDatePickerV2
      label="Date Interval"
      v-model="intervalDates"
      :min="minDate"
      selectionMode="range"
    />

    <OnyxButton label="Reserve this property" mode="default" class="reserve-btn" />
  </div>
</template>

<script setup lang="ts">
import { type DateRange, OnyxButton, OnyxUnstableDatePickerV2 } from 'sit-onyx'
import { ref } from 'vue'

defineProps<{
  price: number,
  disabledDays?: Date[]
}>()


const intervalDates = ref<DateRange>()

const minDate = new Date ()


const formatDate = (date: Date | null): string => {
  if (!date || isNaN(date.getTime())) return 'Select date'
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const year = date.getFullYear()
  return `${month}/${day}/${year}`
}
</script>

<style scoped>
.reservation-card {
  background-color: #ffffff;
  color: #111827;
  border-radius: 20px;
  padding: 1.75rem;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.05);
}

.reservation-header {
  display: flex;
  align-items: baseline;
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

:deep(.reserve-btn) {
  width: 100% !important;
  background-color: #2563eb !important;
  color: #ffffff !important;
  border-radius: 10px !important;
  padding: 0.75rem !important;
  font-weight: 700 !important;
  font-size: 0.9375rem !important;
}
</style>
