<template>
<aside class="menu is-one-fifth">
    <p class="menu-label">系统</p>
    <ul class="menu-list">
        <li v-for="(r, idx) in state.menus" :key="idx">
            <router-link :to="r.path" :class="{ 'is-active' : route.name == r.name }">{{ r.title }}</router-link>
        </li>
    </ul>
</aside>
</template>

<script>
import { onMounted, reactive } from 'vue'
import {routes} from "../router/router"
import { useRoute } from 'vue-router'

export default {
    name: "Menu",
    setup() {
        const route = useRoute()
        const state = reactive({
            menus: []
        })

        onMounted(() => {
            let menus = [];
            for (let i = 0; i < routes.length; i++) {
                let children = routes[i]['children'];
                for (let j = 0; j < children.length; j++) {
                    menus.push(children[j])
                }
            }
            state.menus = menus;
        })

        return {
            route,
            state
        }
    }
}
</script>