<script setup lang="ts">
import { RouterView } from 'vue-router'
import { OnyxToast } from 'sit-onyx';

import Header from '@/components/Header.vue';
import Footer from '@/components/Footer.vue';
</script>

<template>
  <div class="app-layout">
    <Header v-if="!$route.meta.hideHeaderFooter" />

    <main class="main-content">
      <RouterView />
    </main>

    <Footer v-if="!$route.meta.hideHeaderFooter" />

    <OnyxToast />
  </div>
</template>

<style>
:root,
:root,
html[data-onyx-theme] {
  --onyx-color-base-primary-500: #1a41b5 !important;
  --onyx-color-base-primary-600: #143393 !important;
  --onyx-color-base-primary-700: #0f2671 !important;
  --onyx-color-base-primary-100: #e8ecf8 !important;

  --onyx-color-brand-500: #1a41b5 !important;
  --onyx-color-brand-600: #143393 !important;
  --onyx-color-brand-700: #0f2671 !important;

  --onyx-color-component-focus-ring: #1a41b5 !important;
  --onyx-color-form-border-focus: #1a41b5 !important;

  --onyx-font-family-paragraph: var(--onyx-font-family-base) !important;
}

.onyx-button:hover {
  background-color: var(--onyx-color-base-primary-600) !important;
  border-color: var(--onyx-color-base-primary-600) !important;
}

.onyx-button:active {
  background-color: var(--onyx-color-base-primary-700) !important;
  border-color: var(--onyx-color-base-primary-700) !important;
}

.onyx-toast, 
.onyx-toast-container,
[data-onyx-toast-container] {
  top: auto !important;         
  left: auto !important;        
  bottom: 2rem !important;     
  right: 2rem !important;      
  transform: none !important;   
}

html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: var(--onyx-font-family-base, system-ui, sans-serif); 
}

.app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1; 
}
</style>