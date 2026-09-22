<template>
  <OnyxCard class="reservation-card">
    <div class="reservation-header">
      <div class="price-tag">
        <span class="price-amount">${{ price }}</span>
        <span class="price-period"> / night</span>
      </div>
    </div>

    <OnyxUnstableDatePickerV2
      label="Select Check-In and Check-Out Dates"
      v-model="intervalDates"
      :min="minDate"
      selectionMode="range"
      class="date-picker-input"
      :disabledDays="isDisabled"
      :error="hasDisabledDaysInRange ? 'Pick a valid range' : ''"
    />

    <div v-if="nights > 0" class="date-summary-box">
      <div class="summary-cell">
        <span class="cell-label">CHECK-IN</span>
        <span class="cell-value">{{ startDateFormatted }}</span>
      </div>
      <div class="summary-cell border-left">
        <span class="cell-label">CHECK-OUT</span>
        <span class="cell-value">{{ endDateFormatted }}</span>
      </div>
    </div>

    <OnyxButton
      label="Reserve this property"
      mode="default"
      color="primary"
      class="reserve-btn"
      :disabled="isSubmitDisabled"
      @click="emit('submit', emitReservation)"
    />

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
  existingReservations?: [string, string][]
}>()

const emit = defineEmits<{
  (e: 'submit', reservation: { start: string; end: string }): void
}>()

const intervalDates = ref<DateRange>()

const minDate = computed(() => {
  const d = new Date()
  d.setHours(0, 0, 0, 0)
  return d
})

const dayTimestamp = (date: Date): number => {
  const normalized = new Date(date)
  normalized.setHours(0, 0, 0, 0)
  return normalized.getTime()
}

const isDisabled = (targetDate: Date): boolean => {
  if (!props.existingReservations?.length) return false
  const target = dayTimestamp(targetDate)

  return props.existingReservations.some(([checkIn, checkOut]) => {
    const start = dayTimestamp(new Date(checkIn))
    const end = dayTimestamp(new Date(checkOut))
    return start <= target && target < end
  })
}

const parsedRange = computed<{ start: Date | null; end: Date | null }>(() => {
  const range = intervalDates.value
  if (!range) return { start: null, end: null }
  return { start: range.start, end: range.end || null }
})

const hasDisabledDaysInRange = computed(() => {
  const { start, end } = parsedRange.value
  if (!start || !end) return false

  const current = new Date(start)
  current.setHours(0, 0, 0, 0)
  const targetEnd = dayTimestamp(end)

  while (current.getTime() <= targetEnd) {
    if (isDisabled(current)) return true
    current.setDate(current.getDate() + 1)
  }

  return false
})

const nights = computed(() => {
  const { start, end } = parsedRange.value
  if (!start || !end) return 0
  const diffMs = end.getTime() - start.getTime()
  const days = Math.round(diffMs / (1000 * 60 * 60 * 24))
  return days > 0 ? days : 0
})

const totalPrice = computed(() => props.price * nights.value)
const isSubmitDisabled = computed(() => nights.value === 0 || hasDisabledDaysInRange.value)

const dateFormatter = new Intl.DateTimeFormat('en-US', {
  month: '2-digit',
  day: '2-digit',
  year: 'numeric',
})

const outputFormatter = new Intl.DateTimeFormat('en-CA', {
  month: '2-digit',
  day: '2-digit',
  year: 'numeric',
})

const emitReservation = computed(() => {
  return {
    start: outputFormatter.format(parsedRange.value.start || undefined),
    end: outputFormatter.format(parsedRange.value.end || undefined),
  }
})

const formatDate = (date: Date | null): string => {
  return date ? dateFormatter.format(date) : 'Select date'
}

const startDateFormatted = computed(() => formatDate(parsedRange.value.start))
const endDateFormatted = computed(() => formatDate(parsedRange.value.end))
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
  width: 100%;
  border-radius: 8px;
  padding: 0.75rem;
  font-weight: 700;
  font-size: 0.9375rem;
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
