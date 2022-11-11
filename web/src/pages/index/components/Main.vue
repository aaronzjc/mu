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

const API = {
    LoginInfo: "/auth/info/index"
}

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
            try {
                let resp = await Get(API.LoginInfo)
                if (resp.data.code == 10000) {
                    let info = resp.data.data;
                    await store.dispatch("account/initUser", {
                        id: info.id,
                        username: info.username,
                        nickname: info.nickname,
                        avatar: info.avatar
                    });
                }
            } catch(err) {
                // err happends
            }
        }
        onMounted(fetchUserInfo)
    },
    components: {
        Navbar
    }
}
</script>
