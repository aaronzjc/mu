<template>
<nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
        <a class="navbar-item" href="/">
            <img alt="Mu: 快乐摸鱼~" style="width: 28px;height: 28px;" src="../assets/logo.webp" />
        </a>
        <a role="button" :class='[ "navbar-burger", { "is-active": state.open } ]' aria-label="menu" aria-expanded="false" @click="state.open = !state.open">
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
        </a>
    </div>

    <div class="navbar-end">
        <template v-if="!store.getters['account/isLogin']">
            <div class="mini-navbar-opt" v-show="state.open">
                <a class="navbar-item" @click="toLogin">登录</a>
                <a class="navbar-item" @click="toggleTheme">{{ state.theme === "light" ? "黑夜模式" : "白天模式" }}</a>
            </div>
            <div class="navbar-item navbar-opt">
                <a @click="toLogin">登录</a>
            </div>
            <div class="navbar-item navbar-opt">
                <a @click="toggleTheme">{{ state.theme === "light" ? "黑夜模式" : "白天模式" }}</a>
            </div>
        </template>
        <template v-else>
            <div class="mini-navbar-opt" v-show="state.open">
                <span class="navbar-item">欢迎，{{ store.getters['account/getNickname'] }}</span>
                <a :key="idx" v-for="(r, idx) in state.rs" @click="go(r.path)" class="navbar-item">{{ r.title }}</a>
                <a class="navbar-item" @click="toggleTheme">{{ state.theme === "light" ? "黑夜模式" : "白天模式" }}</a>
                <a class="navbar-item" @click="logout">退出登录</a>
            </div>

            <div class="navbar-item has-dropdown is-hoverable navbar-opt">
                <a class="navbar-link">{{ store.getters['account/getNickname'] }}</a>

                <div class="navbar-dropdown is-right">
                    <a :key="idx" v-for="(r, idx) in state.rs" @click="go(r.path)" class="navbar-item">{{ r.title }}</a>
                    <a class="navbar-item" @click="toggleTheme">{{ state.theme === "light" ? "黑夜模式" : "白天模式" }}</a>
                    <span class="navbar-divider"></span>
                    <a class="navbar-item" @click="logout">退出登录</a>
                </div>
            </div>
        </template>
    </div>
</nav>
</template>

<script>
import { routes } from "../router/router";
import { Get, Set, Del } from "@/tools/ls";
import { onMounted, reactive } from 'vue';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

const LIGHT = "light";
const DARK = "dark";
const THEME_KEY = "theme";

export default {
    name: "Navbar",
    setup() {
        const router = useRouter()
        const store = useStore()
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
            Set(THEME_KEY, state.theme, -1);
        }

        let go = (path) => {
            router.push(path).catch(() => {})
            state.open = false
        }
        let toLogin = () => {
            router.push({
                name: "login"
            }).catch(() => {});
        }
        let logout = () => {
            Del("token")
            window.location.href = "/"
        }
        let toggleTheme = () => {
            if (state.theme === LIGHT) {
                initTheme(DARK)
            } else if (state.theme === DARK) {
                initTheme(LIGHT)
            }
            initTheme(state.theme)
            state.open = false
        }

        onMounted(() => {
            state.rs = routes[0].children
            initTheme(Get(THEME_KEY))
        })
        return {
            store,
            state,
            initTheme,
            go,
            toLogin,
            logout,
            toggleTheme
        }
    }
}
</script>
