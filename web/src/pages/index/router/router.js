import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'
import Main from "../components/Main";
import Index from "../components/Index";

const Favor = () => import(/* webpackChunkName: "idx-comps" */ "../components/Favor");
const Login = () => import(/* webpackChunkName: "idx-comps" */ "../components/Login");

const routes = [
    {
        path: '/',
        name: 'default',
        title: "首页",
        component: Main,
        redirect: "index",
        children: [
            {
                path: "/",
                name: "index",
                title: "首页",
                component: Index
            },
            {
                path: "/favor",
                name: "favor",
                title: "我的收藏",
                component: Favor
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