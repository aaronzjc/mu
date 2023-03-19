import { createRouter, createWebHashHistory } from "vue-router";

import Admin from "@adm/layouts/Admin.vue";
import Home from "@adm/views/Home.vue";
import Site from "@adm/views/Site.vue";
import SiteEdit from "@adm/views/SiteEdit.vue";
import Node from "@adm/views/Node.vue";
import NodeEdit from "@adm/views/NodeEdit.vue";
import User from "@adm/views/User.vue";
import Login from "@adm/views/Login.vue";

const routes = [
  {
    path: "/",
    component: Admin,
    redirect: "home",
    children: [
      {
        path: "/home",
        name: "home",
        component: Home,
      },
      {
        path: "/site",
        name: "site",
        component: Site,
      },
      {
        meta: {
          hl: "site",
        },
        path: "/site/edit",
        name: "siteEdit",
        component: SiteEdit,
      },
      {
        path: "/node",
        name: "node",
        component: Node,
      },
      {
        meta: {
          hl: "node",
        },
        path: "/node/edit",
        name: "nodeEdit",
        component: NodeEdit,
      },
      {
        path: "/user",
        name: "user",
        component: User,
      },
    ],
  },
  {
    path: "/login",
    name: "login",
    component: Login,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior(to, from, savedPostion) {
    return savedPostion || { top: 0 };
  },
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
