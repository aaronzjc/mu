<template>
    <aside class="menu is-one-fifth">
        <p class="menu-label">系统</p>
        <ul class="menu-list">
            <li v-for="(r, idx) in menus" :key="idx">
                <router-link :to="r.path" :class="{ 'is-active' : $route.name == r.name }">{{ r.title }}</router-link>
            </li>
        </ul>
    </aside>
</template>

<script>
import {routes} from "../router/router"

export default {
    name: "Menu",
    mounted() {
        var menus = [];
        for (var i = 0; i < routes.length; i++) {
            var children = routes[i]['children'];
            for (var j = 0; j < children.length; j++) {
                menus.push(children[j])
            }
        }

        this.menus = menus;
    },
    data() {
        return {
            menus: []
        }
    }
}
</script>

<style scoped>

</style>