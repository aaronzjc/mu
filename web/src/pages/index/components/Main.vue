<template>
    <main class="container">
        <Navbar></Navbar>
        <router-view></router-view>
    </main>
</template>

<script setup>
import Navbar from "./Navbar.vue"
import { Get } from "@/lib/http"
import { onMounted } from 'vue'
import { useMainStore } from "@idx/store/main";

const API = {
    LoginInfo: "/auth/info/index"
}

const store = useMainStore()

async function fetchUserInfo() {
    // 如果本地没有用户cookie，不发送请求了。
    let token = localStorage.getItem(import.meta.env.VITE_TOKEN_KEY)
    if (!token) {
        return false
    }
    try {
        let resp = await Get(API.LoginInfo)
        if (resp.data.code == 10000) {
            let info = resp.data.data;
            store.setUser({
                username: info.username,
                avatar: info.avatar
            })
        }
    } catch (err) {
        // err happends
    }
}
onMounted(fetchUserInfo)
</script>
