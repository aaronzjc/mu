import { createApp } from 'vue'
import App from './App'

const app = createApp(App)

import * as ls from "@/tools/ls"
import router from "./router/router"

router.beforeEach((to, from, next) => {
    let token = to.query.token;

    if (token != "" && token != undefined && token != null) {
        ls.Set("token", token, -1)
        return router.replace({path: '/'})
    }

    return next()
});

app.use(router)

import client from "@/tools/http"

client.interceptors.response.use(resp => {
    let res = resp.data;
    if (res.code === 10003) {
        router.push({"name": "login"}).catch(() => {});
        return Promise.reject(resp);
    }

    return resp;
});

app.mount('#app')