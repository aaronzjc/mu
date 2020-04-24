/* styles */
import Bulma from "bulma"
import './scss/main.scss'

import Vue from 'vue'

/* router & store */
import router from "./router/router"
import store from "./store"

import App from './App.vue'

Vue.config.productionTip = false

Vue.use(Bulma);

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app')

if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register("./sw.js")
      .then(function(registration) {
        console.log('Registration successful, scope is:', registration.scope);
      })
      .catch(function(error) {
        console.log('Service worker registration failed, error:', error);
      });
}