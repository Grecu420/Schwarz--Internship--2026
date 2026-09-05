import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createOnyx } from "sit-onyx";

import "sit-onyx/style.css";
// import "sit-onyx/global.css";
import "@fontsource-variable/source-code-pro";
import "@fontsource-variable/source-sans-3";


import App from './App.vue'
import router from './router'

const app = createApp(App)

const onyx = createOnyx({
  // if you are using the Vue Router, make sure to pass it here be enable the router integration for onyx
  router: router,
});


app.use(createPinia())
app.use(onyx)

app.mount('#app')
