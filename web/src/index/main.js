/* styles */
import './scss/main.scss'

/* register service worker */
import "./registerSW"

import Vue from 'vue'

/* router & store */
import router from "./router/router"
import store from "./store"

import App from './App.vue'

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app')
