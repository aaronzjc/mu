<template>
    <main class="container">
        <Navbar></Navbar>
        <router-view></router-view>
    </main>
</template>

<script>
import Navbar from "./Navbar"
import {Get} from "../tools/http";

export default {
    name: "Main",
    created() {
        this.fetchUserInfo();
    },
    methods: {
        fetchUserInfo() {
            // 如果本地没有用户cookie，不发送请求了。
            var str = document.cookie;
            if (str === "") return false;
            var skip = true;
            var cookieArr = str.split("; ");
            for (var i = 0; i < cookieArr.length; i++) {
                var arr = cookieArr[i].split("=");
                if (arr[0] === "_token" && arr[1] !== ""){
                    skip = false;
                    break;
                }
            }
            if (skip) return false;
            Get("/info").then(resp => {
                if (resp.data.code == 10000) {
                    var info = resp.data.data;
                    this.$store.dispatch("account/initUser", {
                        id: info.id,
                        username: info.username,
                        avatar: info.avatar
                    });
                }
            });
        }
    },
    components: {
        Navbar
    }
}
</script>