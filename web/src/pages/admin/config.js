import {
    mdiAccount,
    mdiDotsGrid,
    mdiFirefox,
    mdiLogin,
    mdiMonitor,
    mdiServer,
    mdiSquareEditOutline,
    mdiTable
} from '@mdi/js'

export const menus = [
    {
        route: 'home',
        title: '后台总览',
        icon: mdiMonitor,
        active: false
    },
    {
        route: 'site',
        title: '网站管理',
        icon: mdiFirefox
    },
    {
        route: 'node',
        title: '节点管理',
        icon: mdiServer
    },
    {
        route: 'user',
        title: '用户管理',
        icon: mdiAccount
    }
]

export const nodeType = {
    1: "国内",
    2: "海外"
};

export const crawType = {
    1: "JSON",
    2: "HTML"
};