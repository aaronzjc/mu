<template>
<nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
        <a class="navbar-item" href="/">
            <img src="../assets/logo.webp" alt="Mu: 快乐摸鱼~">
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
                <a class="navbar-item" @click="toggleTheme">{{ theme === "light" ? "黑夜模式" : "白天模式" }}</a>
            </div>
            <div class="navbar-item navbar-opt">
                <a @click="toLogin">登录</a>
            </div>
            <div class="navbar-item navbar-opt">
                <a @click="toggleTheme">{{ theme === "light" ? "黑夜模式" : "白天模式" }}</a>
            </div>
        </template>
        <template v-else>
            <div class="mini-navbar-opt" v-show="open">
                <span class="navbar-item">欢迎，{{ $store.getters['account/getUsername'] }}</span>
                <a :key="idx" v-for="(r, idx) in rs" @click="go(r.path)" class="navbar-item">{{ r.title }}</a>
                <a class="navbar-item" @click="toggleTheme">{{ theme === "light" ? "黑夜模式" : "白天模式" }}</a>
                <a class="navbar-item" @click="logout">退出登录</a>
            </div>

            <div class="navbar-item has-dropdown is-hoverable navbar-opt">
                <a class="navbar-link">{{ $store.getters['account/getUsername'] }}</a>

                <div class="navbar-dropdown is-right">
                    <a :key="idx" v-for="(r, idx) in rs" @click="go(r.path)" class="navbar-item">{{ r.title }}</a>
                    <a class="navbar-item" @click="toggleTheme">{{ theme === "light" ? "黑夜模式" : "白天模式" }}</a>
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

const LIGHT = "light";
const DARK = "dark";
const THEME_KEY = "theme";

export default {
    name: "Navbar",
    mounted() {
        this.rs = routes[0].children;
        var t = localStorage.getItem(THEME_KEY);
        this.initTheme(t);
    },
    data() {
        return {
            open: false,
            rs: [],
            theme: "",
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
        },
        initTheme(type) {
            if (type != LIGHT && type != DARK) {
                type = LIGHT;
            }
            var ht = document.getElementsByTagName("html")[0];
            if (type === DARK) {
                ht.className = ht.className.trim() + " " + DARK;
            } else if (type === "light") {
                ht.className = ht.className.replace(DARK, "");
            }
            this.theme = type;
            localStorage.setItem(THEME_KEY, this.theme);
        },
        toggleTheme() {
            if (this.theme === LIGHT) {
                this.initTheme(DARK);
            } else if (this.theme === DARK) {
                this.initTheme(LIGHT);
            }
            this.initTheme(this.theme);
            this.open = false;
        }
    }
}
</script>