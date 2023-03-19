import { createApp } from "vue";
import App from "./App.vue";
import "./styles/main.scss";

const app = createApp(App);

/** 路由 */
import router from "./router";
app.use(router);


import client from "@/lib/http";
client.interceptors.response.use((resp) => {
  let res = resp.data;
  if (res.code === 10003) {
    router.push({ name: "login" }).catch(() => {});
    return Promise.reject(resp);
  }

  return resp;
});

/** 全局状态管理 */
import { createPinia } from "pinia";
const pinia = createPinia();
app.use(pinia);

/** 自定义组件 */
import { Toast } from "@adm/components/toast";
app.use(Toast);

app.mount("#app");
