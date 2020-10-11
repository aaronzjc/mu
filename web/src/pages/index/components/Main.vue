<template>
<main class="container">
    <Navbar></Navbar>
    <router-view></router-view>
</main>
</template>

<script>
import Navbar from "./Navbar"
import { Get } from "@/tools/http"
import { Get as lsGet } from "@/tools/ls"
import { onMounted } from 'vue'
import { useStore } from "vuex";

export default {
    name: "Main",
    setup() {
        const store = useStore()
        async function fetchUserInfo() {
            // 如果本地没有用户cookie，不发送请求了。
            let token = lsGet("token")
            if (!token) {
                return false
            }
            let resp = await Get("/api/info")
            if (resp.data.code == 10000) {
                let info = resp.data.data;
                await store.dispatch("account/initUser", {
                    id: info.id,
                    username: info.username,
                    nickname: info.nickname,
                    avatar: info.avatar
                });
            } else {
                console.log(info)
            }
        }
        onMounted(fetchUserInfo)
    },
    components: {
        Navbar
    }
}
</script>
