<template>
    <div class="container">
        <Navbar></Navbar>
        <router-view></router-view>
    </div>
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
            Get("/api/info").then(resp => {
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

<style lang="scss" scoped>
@import 'bulma/sass/utilities/mixins';

.container {
    max-width: 960px;
}

@include until($desktop) {
    .content-box {
        padding: 0 0.75rem;
    }
}
</style>