<template>
<div class="container">
    <div class="login has-text-centered columns">
        <div class="column">
            <h4 class="title is-4">登录</h4>
            <a v-for="(item, idx) in auth" :key="idx" :href="item.url"  class="button is-medium is-white" :title="item.name">
            <span class="icon is-medium">
              <i class="fab fa-github"></i>
            </span>
            </a>
        </div>
    </div>
</div>
</template>

<script>
import "../assets/font-awesome"

import {Get} from "../tools/http";

export default {
    name: "Login",
    created() {
        this.fetchConfig();
    },
    data() {
        return {
            auth: []
        }
    },
    methods: {
        fetchConfig() {
            Get("/auth_config").then(resp => {
                if (resp.data.code === 10000) {
                    this.auth = resp.data.data;
                } else {
                    alert(resp.data.msg);
                }
            })
        }
    }
}
</script>

<style scoped>
.login {
    margin: 4rem auto;
    width: 30%;
    height: 40%;

}
</style>