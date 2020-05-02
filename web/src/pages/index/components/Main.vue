<template>
    <main class="container">
        <Navbar></Navbar>
        <router-view></router-view>
    </main>
</template>

<script>
import Navbar from "./Navbar"
import {Get} from "@/tools/http";
import {Get as lsGet} from "@/tools/ls"

export default {
    name: "Main",
    created() {
        this.fetchUserInfo();
    },
    methods: {
        fetchUserInfo() {
            // 如果本地没有用户cookie，不发送请求了。
            let token = lsGet("token")
            if (!token) {
                return false
            }
            Get("/api/info").then(resp => {
                if (resp.data.code == 10000) {
                    var info = resp.data.data;
                    this.$store.dispatch("account/initUser", {
                        id: info.id,
                        username: info.username,
                        avatar: info.avatar
                    });
                } else {
                    console.log(info)
                }
            });
        }
    },
    components: {
        Navbar
    }
}
</script>