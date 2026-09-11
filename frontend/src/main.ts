import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createOnyx } from 'sit-onyx'

import 'leaflet/dist/leaflet.css'



import App from './App.vue'
import router from './router'
import "sit-onyx/style.css"

createApp(App)
  .use(createPinia())
  .use(router)
  .use(createOnyx({})) 
  .mount('#app')