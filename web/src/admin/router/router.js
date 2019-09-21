import Vue from 'vue'
import Router from 'vue-router'
import Site from "../components/Site"
import Node from "../components/Node"

Vue.use(Router);

const routes = [
    {
        path: '/',
        name: 'default',
        component: Site,
        redirect: "/site",
        hide: true,
    },
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
];

const router = new Router({
    routes
});

export {routes}

export default router