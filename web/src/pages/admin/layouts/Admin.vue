<template>
    <div class="container-box">
        <div
            :class="[
                'content-box has-aside-left has-aside-mobile-transition has-navbar-fixed-top has-aside-expanded',
                { 'has-aside-mobile-expanded': styleStore.isAsideMobileOpen }
            ]"
        >
            <Menu></Menu>

            <RouterView v-slot="{ Component, route }">
                <Transition name="page">
                    <component :is="Component" :key="route.path" />
                </Transition>
            </RouterView>
        </div>
    </div>
</template>

<script setup>
import Menu from '@adm/components/Menu.vue'
import { useStyleStore } from '@adm/stores/style.js'
import {Get} from '@/lib/http'
import { useMainStore } from '../../index/store/main';
import { onBeforeMount } from 'vue';

const styleStore = useStyleStore()
const mainStore = useMainStore()

const API = {
    loginInfo: "/auth/info/admin"
}

async function fetchLogin() {
    try {
        let resp = await Get(API.loginInfo)
        if (resp.data.code === 10000) {
            mainStore.setUser({
                username: resp.data.data.username,
                avatar: resp.data.data.avatar
            })
        }
    } catch(err) {
        console.log(err)
    }
}
onBeforeMount(fetchLogin)
</script>
