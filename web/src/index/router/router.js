import Vue from 'vue'
import Router from 'vue-router'
import Login from "../components/Login";
import Main from "../components/Main";
import Index from "../components/Index";
import Favor from "../components/Favor";

Vue.use(Router);

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

const publicRouters = [
    {
        path: '/login',
        name: 'login',
        component: Login
    }
];

const router = new Router({
    routes: routes.concat(publicRouters)
});

export {routes}

export default router