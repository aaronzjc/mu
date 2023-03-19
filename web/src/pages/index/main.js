import { createApp } from "vue";
import App from "./App.vue";

const app = createApp(App);

/** 样式 */
import "./scss/main.scss";

/** 路由 */
import router from "./router/router";
app.use(router);

import client from "@/lib/http";
client.interceptors.response.use((resp) => {
  let res = resp.data;
  if (res.code === 10003) {
    localStorage.removeItem(import.meta.env.VITE_TOKEN_KEY);
    router.push({ name: "login" }).catch(() => {});
    return Promise.reject(resp);
  }
  return resp;
});

/** 状态 */
import { createPinia } from "pinia";
const pinia = createPinia();
app.use(pinia);

app.mount("#app");

/* register sw */
import "./registerSW";
