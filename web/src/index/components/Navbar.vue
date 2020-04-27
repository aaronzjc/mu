<template>
<nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
        <a class="navbar-item" href="/">
            <img alt="Mu: 快乐摸鱼~" src="data:image/webp;base64,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" />
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