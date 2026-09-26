<template>
  <div class="pagination-controls">
    <OnyxIconButton
      mode="outline"
      label="Previous Page"
      :icon="iconChevronLeft"
      :disabled="isPrevDisabled || isLoading"
      class="pagination-btn"
      @click="emit('previous-page')"
    />

    <span class="page-indicator">Page {{ currentPage + 1 }}</span>

    <OnyxIconButton
      mode="outline"
      label="Next Page"
      :icon="iconChevronRight"
      :disabled="isNextDisabled || isLoading"
      class="pagination-btn"
      @click="emit('next-page')"
    />
  </div>
</template>

<script setup lang="ts">
import { OnyxIconButton } from 'sit-onyx'
import { iconChevronLeft, iconChevronRight } from '@sit-onyx/icons'

defineProps<{
  currentPage: number
  isPrevDisabled: boolean
  isNextDisabled: boolean
  isLoading?: boolean
}>()

const emit = defineEmits<{
  (e: 'previous-page'): void
  (e: 'next-page'): void
}>()
</script>

<style scoped>
.pagination-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1.5rem;
}

.page-indicator {
  color: #374151;
  font-size: 0.875rem;
  font-weight: 600;
}

:deep(.pagination-btn) {
  border-radius: 50% !important;
  width: 40px !important;
  height: 40px !important;
  padding: 0 !important;
  display: flex !important;
  justify-content: center !important;
  align-items: center !important;
  border: 1px solid #e5e7eb !important;
  background-color: #ffffff !important;
  color: #374151 !important;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.03) !important;
}

:deep(.pagination-btn:not(:disabled):hover) {
  background-color: #f9fafb !important;
  border-color: #d1d5db !important;
  color: #111827 !important;
}

:deep(.pagination-btn:disabled) {
  opacity: 0.4;
  cursor: not-allowed;
  background-color: #f3f4f6 !important;
  border-color: #e5e7eb !important;
}
</style>