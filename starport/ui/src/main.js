import Vue from 'vue'
import App from './App.vue'
import router from './router'

import '@/styles/_normalize.css'
import '@/styles/_colors.css'
import '@/styles/_typography.css'

Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
