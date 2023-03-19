<template>
    <nav class="navbar" role="navigation" aria-label="main navigation">
        <div class="navbar-brand">
            <a class="navbar-item" href="/">
                <img alt="Mu: 快乐摸鱼~" style="width: 28px;height: 28px;" src="../assets/logo.webp" />
            </a>
            <a role="button" :class='["navbar-burger", { "is-active": state.open }]' aria-label="menu" aria-expanded="false"
                @click="state.open = !state.open">
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
            </a>
        </div>

        <div class="navbar-end">
            <template v-if="!store.isLogin">
                <div class="mini-navbar-opt" v-show="state.open">
                    <span class="navbar-item" @click="toLogin">登录</span>
                    <span class="navbar-item" @click="toggleTheme">{{ state.theme === "light" ? "黑夜模式" : "白天模式" }}</span>
                </div>
                <div class="navbar-item navbar-opt">
                    <span @click="toLogin">登录</span>
                </div>
                <div class="navbar-item navbar-opt">
                    <span @click="toggleTheme">{{ state.theme === "light" ? "黑夜模式" : "白天模式" }}</span>
                </div>
            </template>
            <template v-else>
                <div class="mini-navbar-opt" v-show="state.open">
                    <span class="navbar-item">欢迎，{{ store.userInfo.username }}</span>
                    <span :key="idx" v-for="(r, idx) in state.rs" @click="go(r.path)" class="navbar-item">{{ r.title
                    }}</span>
                    <span class="navbar-item" @click="toggleTheme">{{ state.theme === "light" ? "黑夜模式" : "白天模式" }}</span>
                    <span class="navbar-item" @click="logout">退出登录</span>
                </div>

                <div class="navbar-item has-dropdown is-hoverable navbar-opt">
                    <span class="navbar-link">{{ store.userInfo.username }}</span>

                    <div class="navbar-dropdown is-right">
                        <span :key="idx" v-for="(r, idx) in state.rs" @click="go(r.path)" class="navbar-item">{{ r.title
                        }}</span>
                        <span class="navbar-item" @click="toggleTheme">{{ state.theme === "light" ? "黑夜模式" : "白天模式"
                        }}</span>
                        <span class="navbar-divider"></span>
                        <span class="navbar-item" @click="logout">退出登录</span>
                    </div>
                </div>
            </template>
        </div>
    </nav>
</template>

<script setup>
import { routes } from "../router/router";
import { onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useMainStore } from "../store/main";

const LIGHT = "light";
const DARK = "dark";
const THEME_KEY = "theme";

const router = useRouter()
const store = useMainStore()

const state = reactive({
    open: false,
    rs: [], // 菜单项
    theme: ""
})

const initTheme = (type) => {
    if (type !== LIGHT && type !== DARK) {
        type = LIGHT;
    }
    var ht = document.getElementsByTagName("html")[0];
    if (type === DARK) {
        ht.className = ht.className.trim() + " " + DARK;
    } else if (type === "light") {
        ht.className = ht.className.replace(DARK, "");
    }
    state.theme = type;
    localStorage.setItem(THEME_KEY, state.theme);
}

let go = (path) => {
    router.push(path).catch(() => { })
    state.open = false
}
let toLogin = () => {
    router.push({
        name: "login"
    }).catch(() => { });
}
let logout = () => {
    localStorage.removeItem(import.meta.env.VITE_TOKEN_KEY)
    window.location.href = "/"
}
let toggleTheme = () => {
    if (state.theme === LIGHT) {
        initTheme(DARK)
    } else if (state.theme === DARK) {
        initTheme(LIGHT)
    }
    state.open = false
}

onMounted(() => {
    state.rs = routes[0].children
    initTheme(localStorage.getItem(THEME_KEY))
})
</script>

<style>
.navbar-item span {
    color: #3273dc;
    cursor: pointer;
}
</style>