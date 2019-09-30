import Vue from 'vue'
import Router from 'vue-router'
import Site from "../components/Site"
import Node from "../components/Node"
import Dashboard from "../components/Dashboard";
import Login from "../components/Login";

Vue.use(Router);

const routes = [
    {
        path: '/',
        name: 'default',
        title: "系统",
        component: Dashboard,
        redirect: "/site",
        children : [
            {
                path: '/site',
                name: 'site',
                component: Site,
                title: "网站管理"
            },
            {
                path: '/node',
                name: 'node',
                component: Node,
                title: "节点管理"
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