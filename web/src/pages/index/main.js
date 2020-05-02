/* styles */
import './scss/main.scss'

/* register service worker */
import "./registerSW"

import Vue from 'vue'

/* router & store */
import router from "./router/router"
import client from "@/tools/http"
import * as ls from "@/tools/ls"

client.interceptors.response.use(resp => {
  let res = resp.data;
  if (res.code === 10003) {
    router.push({"name": "login"}).catch(() => {});
    return Promise.reject(resp);
  }

  return resp;
});

router.beforeEach((to, from, next) => {
  let token = to.query.token;

  if (token != "" && token != undefined && token != null) {
    ls.Set("token", token, -1)
    return router.replace({path: '/'})
  }

  return next()
});

import store from "./store"
import App from './App.vue'

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app')
