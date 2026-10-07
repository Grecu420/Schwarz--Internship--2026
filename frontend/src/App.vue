<script setup lang="ts">
import { RouterView } from 'vue-router'
import { OnyxToast } from 'sit-onyx'

import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
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
html,
body {
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
