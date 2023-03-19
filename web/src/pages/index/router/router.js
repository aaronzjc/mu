import { createRouter, createWebHashHistory } from "vue-router";
import Main from "../components/Main.vue";
import Content from "../components/Content.vue";

const Favor = () => import("../components/Favor.vue");
const Login = () => import("../components/Login.vue");

const routes = [
  {
    path: "/",
    name: "default",
    title: "首页",
    component: Main,
    redirect: "index",
    children: [
      {
        path: "/",
        name: "index",
        title: "首页",
        component: Content,
      },
      {
        path: "/favor",
        name: "favor",
        title: "我的收藏",
        component: Favor,
      },
    ],
  },
];

export { routes };

const publicRoutes = [
  {
    path: "/login",
    name: "login",
    component: Login,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes: routes.concat(publicRoutes),
});

router.beforeEach((to, from, next) => {
  let token = to.query.token;
  if (token != "" && token != undefined && token != null) {
    localStorage.setItem(import.meta.env.VITE_TOKEN_KEY, token);
    router.replace({ path: "/" });
  } else {
    next();
  }
});

export default router;
