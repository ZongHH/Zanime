import { createApp } from 'vue'
import App from './App.vue'
import store from './static/js/storeVuex.js'
import router from './static/js/router.js'

const app = createApp(App)

app.use(router)
app.use(store)

app.mount('#app')
