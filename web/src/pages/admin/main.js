import Vue from 'vue'

import router from "./router/router"
import * as ls from "@/tools/ls"
import client from "@/tools/http"

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

import App from './App.vue'

Vue.config.productionTip = false;

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
