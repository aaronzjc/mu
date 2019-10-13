<template>
<nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
        <a class="navbar-item" href="/">
            <img src="../assets/logo.png" alt="Mu: 快乐摸鱼~">
        </a>
        <a role="button" :class='[ "navbar-burger", { "is-active": open } ]' aria-label="menu" aria-expanded="false" @click="open = !open">
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
        </a>
    </div>

    <div class="navbar-end">
        <template v-if="!$store.getters['account/isLogin']">
            <div class="mini-navbar-opt" v-show="open">
                <a class="navbar-item" @click="toLogin">登录</a>
            </div>
            <div class="navbar-item navbar-opt">
                <a @click="toLogin">登录</a>
            </div>
        </template>
        <template v-else>
            <div class="mini-navbar-opt" v-show="open">
                <span class="navbar-item">欢迎，{{ $store.getters['account/getUsername'] }}</span>
                <a :key="idx" v-for="(r, idx) in rs" @click="go(r.path)" class="navbar-item">{{ r.title }}</a>
                <a class="navbar-item" @click="logout">退出登录</a>
            </div>

            <div class="navbar-item has-dropdown is-hoverable navbar-opt">
                <a class="navbar-link">{{ $store.getters['account/getUsername'] }}</a>

                <div class="navbar-dropdown is-right">
                    <a :key="idx" v-for="(r, idx) in rs" @click="go(r.path)" class="navbar-item">{{ r.title }}</a>
                    <span class="navbar-divider"></span>
                    <a class="navbar-item" @click="logout">退出登录</a>
                </div>
            </div>
        </template>
    </div>
</nav>
</template>

<script>
import {routes} from "../router/router";
import {Get} from "../tools/http";

export default {
    name: "Navbar",
    mounted() {
        this.rs = routes[0].children;
    },
    data() {
        return {
            open: false,
            rs: [],
        }
    },
    methods: {
        go(path) {
            this.$router.push(path).catch(() => {});
            this.open = false;
        },
        toLogin() {
            this.$router.push({name: "login"}).catch(() => {});
        },
        logout() {
            Get("/logout").then(resp => {
                if (resp.data.code == 10000) {
                    window.location.href = "/"
                }
            })
        }
    }
}
</script>

<style lang="scss" scoped>
@import 'bulma/sass/utilities/mixins';

.navbar {
    margin-bottom: 1rem;
}

.navbar-burger:hover {
    background: none;
}

.mini-navbar-opt {
    background: #f1f1f1;
}

@include until($desktop) {
    .navbar-item.navbar-opt {
        display: none;
    }
    .mini-navbar-opt {
        display: block;
    }
}
</style>