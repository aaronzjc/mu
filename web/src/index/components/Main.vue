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