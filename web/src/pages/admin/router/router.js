import { createRouter, createWebHistory } from 'vue-router'

import Site from "../components/Site"
import Node from "../components/Node"
import Dashboard from "../components/Dashboard";
import Login from "../components/Login";
import User from "../components/User";

const routes = [
    {
        path: '/',
        name: 'default',
        title: "系统",
        component: Dashboard,
        redirect: "site",
        children : [
            {
                path: '/site',
                name: 'site',
                component: Site,
                title: "站点管理"
            },
            {
                path: '/node',
                name: 'node',
                component: Node,
                title: "节点管理"
            },
            {
                path: '/user',
                name: 'user',
                component: User,
                title: "用户列表"
            }
        ]
    },
];

const publicRoutes = [
    {
        path: '/login',
        name: 'login',
        component: Login
    }
];

const router = createRouter({
    history: createWebHashHistory(),
    routes: routes.concat(publicRoutes)
})

export {routes}

export default router