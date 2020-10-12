/* styles */
import "./scss/main.scss"

/* register sw */
import "./registerSW"

import { createApp } from 'vue'
import App from './App'

const app = createApp(App)

/* vue-router staff */
import router from "./router/router";
import client from "@/tools/http"
import * as ls from "@/tools/ls"

client.interceptors.response.use(resp => {
    let res = resp.data;
    if (res.code === 10003) {
      ls.Del("token")
      router.push({"name": "login"}).catch(() => {});
      return Promise.reject(resp);
    }

    return resp;
});
  
router.beforeEach((to, from, next) => {
  let token = to.query.token;

  if (token != "" && token != undefined && token != null) {
    ls.Set("token", token, -1)
    router.replace({path: '/'})
    return
  }

  next()
});

app.use(router)

/* vuex staff */
import { store } from "./store"

app.use(store)

/* mount #app */
app.mount('#app')