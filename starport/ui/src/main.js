import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import UUID from 'vue-uuid'

import '@/styles/main.css'

Vue.config.productionTip = false

Vue.use(UUID)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
